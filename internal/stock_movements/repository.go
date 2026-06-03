package stock_movements

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/helpers"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

type excelData struct {
	name        string
	slug        string
	price       int64
	quantity    int64
	barcode     string
	category_id int64
}

func NewStockMovementRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetStockMovementByUUID(ctx context.Context, uuid string) (*models.StockMovement, error) {
	sm := new(models.StockMovement)

	err := r.db.NewSelect().
		Model(sm).
		Where("uuid = ?", uuid).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return sm, nil
}

func (r *Repository) GetStockMovements(ctx context.Context, q dto.BaseFilters) ([]models.StockMovement, int, error) {
	var sm []models.StockMovement

	count, err := r.db.NewSelect().
		Model(&sm).
		Limit(q.PageSize).
		Offset(q.Page).
		Order(q.Sort).
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	return sm, count, nil
}

func (r *Repository) CreateStockMovement(ctx context.Context, sm *models.StockMovement) error {
	_, err := r.db.NewInsert().Model(sm).Exec(ctx)

	return err
}

func (r *Repository) UpdateStockMovement(ctx context.Context, sm *models.StockMovement) error {
	_, err := r.db.NewUpdate().
		Model(sm).
		Column("product_id", "change_amount", "reason", "reference_id").
		Set("updated_at = ?", time.Now()).
		WherePK().
		Exec(ctx)

	return err
}

func (r *Repository) DeleteStockMovement(ctx context.Context, uuid string) error {
	_, err := r.db.NewDelete().
		Model((*models.StockMovement)(nil)).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}

func (r *Repository) UpdateStockMovementStatus(ctx context.Context, uuid string) error {
	_, err := r.db.NewUpdate().
		Model((*models.StockMovement)(nil)).
		Set("status = CASE WHEN status = 'A' THEN 'I' ELSE 'A' END").
		Set("updated_at = ?", time.Now()).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}

func (r *Repository) ProcessBulkUpload(ctx context.Context, rows [][]string, pics [][]byte) ([]BulkUploadErrResponse, error) {
	// Bulk upload for products using excelize
	// BATCH INSERT instead of single line
	// Stage 1: []Product insert on conflict sku update RETURNING id
	// double check prices using function to avoid panic if typo with string
	// Stage 2: []Inventory insert on conflict product_id set quantity = inventory.quantity + EXCLUDED.quantity to add new quantity to the current quantity
	// Stage 3: []StockMovement always insert no update for audit trail and with tag FROM_IMPORTS

	// to consider when Processing Bulk Uploads
	// Goroutine: Best for "Right Now" background processing.
	// Cron Job: Best for "Late Night" batch processing.
	// Worker Pool: The middle ground—processes immediately but limits how many run at once so your server doesn't explode.

	var validProducts []models.Product
	var invalidProducts []BulkUploadErrResponse
	var categories []models.Category

	skuMap := make(map[string]excelData)
	skuMapExists := make(map[string]int64)
	nameMapExists := make(map[string]string)
	categoryMap := make(map[string]int64)

	if len(rows) <= 1 {
		return nil, errors.New("failed to load rows")
	}

	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		err := tx.NewSelect().Model(&categories).Scan(ctx)
		if err != nil {
			return err
		}

		for _, cat := range categories {
			categoryMap[cat.Name] = cat.ID
		}

		existingProducts, err := checkProductExisting(ctx, tx, rows)
		if err != nil {
			return fmt.Errorf("cannot fetch products %v", err)
		}

		if len(existingProducts) > 0 {
			for _, p := range existingProducts {
				skuMapExists[p.SKU] = p.ID
				nameMapExists[p.Name] = p.SKU
			}
		}

		validProducts, invalidProducts, err = validateProducts(rows, skuMap, categoryMap, skuMapExists, nameMapExists)

		_, err = tx.NewInsert().
			Model(&validProducts).
			On("CONFLICT (sku) DO UPDATE").
			Set("name = EXCLUDED.name").
			Set("slug = EXCLUDED.slug").
			Set("description = EXCLUDED.description").
			Set("price = EXCLUDED.price").
			Set("barcode = EXCLUDED.barcode").
			Set("category_id = EXCLUDED.category_id").
			Returning("id", "sku").
			Exec(ctx)
		if err != nil {
			return err
		}

		if err := insertInventory(ctx, tx, validProducts, skuMap); err != nil {
			return err
		}

		if err := insertStockMovement(ctx, tx, validProducts, skuMap); err != nil {
			return err
		}

		return err
	})

	return invalidProducts, err
}

func validateProducts(rows [][]string, skuMap map[string]excelData, categoryMap, skuMapExists map[string]int64, nameMapExists map[string]string) ([]models.Product, []BulkUploadErrResponse, error) {
	var validProducts []models.Product
	var invalidProducts []BulkUploadErrResponse

	for i, r := range rows[1:] {
		catRow := strings.TrimSpace(r[6])
		categoryID, exists := categoryMap[catRow]
		if !exists {
			return nil, nil, fmt.Errorf("category %s not found on row %d", catRow, i+2)
		}

		sku := r[0]
		name := r[1]
		image := ""

		if len(r) == 8 {
			image = r[7]
		}

		data := excelData{
			name:        name,
			slug:        r[2],
			price:       int64(parseToNumeric(r[3]) * 100),
			quantity:    int64(parseToNumeric(r[4])),
			barcode:     r[5],
			category_id: categoryID,
		}
		skuMap[sku] = data

		_, skuExists := skuMapExists[sku]
		_, nameExists := nameMapExists[name]

		if !skuExists && nameExists {
			invalidProducts = append(invalidProducts, BulkUploadErrResponse{Row: i + 2, Name: name, SKU: sku, ExistingSKU: nameMapExists[name]})
			continue
		}

		item := models.Product{
			SKU:         sku,
			Name:        helpers.CleanSpecialChars(data.name),
			Description: helpers.CleanSpecialChars(data.name),
			Slug:        helpers.CleanSpecialChars(data.slug),
			Price:       data.price,
			Barcode:     data.barcode,
			CategoryID:  categoryID,
			Image:       image,
		}

		validProducts = append(validProducts, item)
	}

	return validProducts, invalidProducts, nil
}

func checkProductExisting(ctx context.Context, tx bun.Tx, rows [][]string) ([]models.Product, error) {
	var existingProducts []models.Product

	if len(rows) <= 1 {
		return existingProducts, nil
	}

	if len(rows[0]) < 9 {
		return nil, fmt.Errorf("invalid excel header structure")
	}

	excelSKUs := make([]string, 0, len(rows)-1)
	excelProductNames := make([]string, 0, len(rows)-1)

	for _, r := range rows[1:] {
		if len(r) < 2 {
			continue
		}

		sku := helpers.CleanSpecialChars(r[0])
		name := helpers.CleanSpecialChars(r[1])

		if sku != "" {
			excelSKUs = append(excelSKUs, sku)
		}
		if name != "" {
			excelProductNames = append(excelProductNames, name)
		}
	}

	if len(excelSKUs) == 0 && len(excelProductNames) == 0 {
		return existingProducts, nil
	}

	query := tx.NewSelect().
		Model(&existingProducts).
		WhereGroup(" OR ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.WhereOr("sku IN (?)", bun.List(excelSKUs)).
				WhereOr("name IN (?)", bun.List(excelProductNames))
		})

	return existingProducts, query.Scan(ctx)
}

func insertInventory(ctx context.Context, tx bun.Tx, validProducts []models.Product, skuMap map[string]excelData) error {
	var inventoryItems []models.Inventory

	for _, p := range validProducts {
		if _, ok := skuMap[p.SKU]; ok {
			qty := skuMap[p.SKU].quantity

			inventoryItem := models.Inventory{
				ProductID:         p.ID,
				Quantity:          qty,
				LowStockThreshold: 5,
			}

			inventoryItems = append(inventoryItems, inventoryItem)
		}
	}

	_, err := tx.NewInsert().
		Model(&inventoryItems).
		On("CONFLICT (product_id) DO UPDATE").
		Set("quantity = i.quantity + EXCLUDED.quantity").
		Returning("id").
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func insertStockMovement(ctx context.Context, tx bun.Tx, validProducts []models.Product, skuMap map[string]excelData) error {
	var stock_movements []models.StockMovement

	for _, p := range validProducts {
		if _, ok := skuMap[p.SKU]; ok {
			qty := skuMap[p.SKU].quantity

			if qty > 0 {
				stock_movement := models.StockMovement{
					ProductID:    p.ID,
					ChangeAmount: qty,
					Reason:       "FROM_IMPORT",
				}

				stock_movements = append(stock_movements, stock_movement)
			}
		}
	}

	if _, err := tx.NewInsert().Model(&stock_movements).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func parseToNumeric(s string) float64 {
	if s == "" {
		return 0
	}

	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "-", "")

	price, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Printf("Error converting price: %v\n", err)
		return 0
	}

	return price
}
