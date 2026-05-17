package stock_movements

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
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

func (r *Repository) ProcessBulkUpload(ctx context.Context, rows [][]string) error {
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

	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		var categories []models.Category
		type excelData struct {
			sku         string
			name        string
			slug        string
			price       int64
			quantity    int64
			barcode     string
			category_id int64
		}
		skuMap := make(map[string]excelData)
		categoryMap := make(map[string]int64)

		var items []models.Product

		err := tx.NewSelect().Model(&categories).Scan(ctx)
		if err != nil {
			return err
		}

		for _, cat := range categories {
			categoryMap[cat.Name] = cat.ID
		}

		parseToNumeric := func(s string) float64 {
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

		for i, r := range rows {
			if i == 0 && len(r) < 8 {
				continue
			}

			catRow := strings.TrimSpace(r[6])
			categoryID, exists := categoryMap[catRow]
			if !exists {
				return fmt.Errorf("category %s not found on row %d", catRow, i+1)
			}

			sku := r[0]
			data := excelData{
				name:        r[1],
				slug:        r[2],
				price:       int64(parseToNumeric(r[3]) * 100),
				quantity:    int64(parseToNumeric(r[4])),
				barcode:     r[5],
				category_id: categoryID,
			}
			skuMap[sku] = data

			item := models.Product{
				SKU:         sku,
				Name:        data.name,
				Description: data.name,
				Slug:        data.slug,
				Price:       data.price,
				Barcode:     data.barcode,
				CategoryID:  categoryID,
			}

			items = append(items, item)
		}

		_, err = tx.NewInsert().
			Model(&items).
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

		var inventoryItems []models.Inventory
		var stock_movements []models.StockMovement

		for _, p := range items {
			if _, ok := skuMap[p.SKU]; ok {
				qty := skuMap[p.SKU].quantity

				inventoryItem := models.Inventory{
					ProductID:         p.ID,
					Quantity:          qty,
					LowStockThreshold: 5,
				}

				inventoryItems = append(inventoryItems, inventoryItem)

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

		_, err = tx.NewInsert().
			Model(&inventoryItems).
			On("CONFLICT (product_id) DO UPDATE").
			Set("quantity = i.quantity + EXCLUDED.quantity").
			Returning("id").
			Exec(ctx)
		if err != nil {
			return err
		}

		if _, err := tx.NewInsert().Model(&stock_movements).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}
