package data

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/XaiPhyr/rdev-go-api/internal/dto"
	"github.com/uptrace/bun"
)

type StockMovementRepository struct {
	db *bun.DB
}

func NewStockMovementRepository(db *bun.DB) *StockMovementRepository {
	return &StockMovementRepository{db: db}
}

func (r *StockMovementRepository) GetStockMovementByUUID(ctx context.Context, uuid string) (*StockMovement, error) {
	sm := new(StockMovement)

	err := r.db.NewSelect().
		Model(sm).
		Where("uuid = ?", uuid).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return sm, nil
}

func (r *StockMovementRepository) GetStockMovements(ctx context.Context, q dto.BaseFilters) ([]StockMovement, int, error) {
	var sm []StockMovement

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

func (r *StockMovementRepository) CreateStockMovement(ctx context.Context, sm *StockMovement) error {
	_, err := r.db.NewInsert().Model(sm).Exec(ctx)

	return err
}

func (r *StockMovementRepository) UpdateStockMovement(ctx context.Context, sm *StockMovement) error {
	_, err := r.db.NewUpdate().
		Model(sm).
		Column("product_id", "change_amount", "reason", "reference_id").
		Set("updated_at = ?", time.Now()).
		WherePK().
		Exec(ctx)

	return err
}

func (r *StockMovementRepository) DeleteStockMovement(ctx context.Context, uuid string) error {
	_, err := r.db.NewDelete().
		Model((*StockMovement)(nil)).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}

func (r *StockMovementRepository) UpdateStockMovementStatus(ctx context.Context, uuid string) error {
	_, err := r.db.NewUpdate().
		Model((*StockMovement)(nil)).
		Set("status = CASE WHEN status = 'A' THEN 'I' ELSE 'A' END").
		Set("updated_at = ?", time.Now()).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}

func (r *StockMovementRepository) ProcessBulkUpload(ctx context.Context, rows [][]string) error {
	// Bulk upload for products using excelize
	// BATCH INSERT instead of single line
	// Stage 1: []Product insert on conflict sku update RETURNING id
	// double check prices using function to avoid panic if typo with string
	// Stage 2: []Inventory insert on conflict product_id set quantity = inventory.quantity + EXCLUDED.quantity to add new quantity to the current quantity
	// Stage 3: []StockMovement always insert no update for audit trail and with tag FROM_IMPORTS

	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		var categories []Category
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

		var products []Product

		err := tx.NewSelect().Model(&categories).Scan(ctx)
		if err != nil {
			return err
		}

		for _, cat := range categories {
			categoryMap[cat.Name] = cat.ID
		}

		convertToInt := func(s string) int64 {
			s = strings.TrimSpace(s)
			s = strings.ReplaceAll(s, "$", "")
			s = strings.ReplaceAll(s, ",", "")

			price, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
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
				price:       convertToInt(r[3]) * 100,
				quantity:    convertToInt(r[4]),
				barcode:     r[5],
				category_id: categoryID,
			}
			skuMap[sku] = data

			product := Product{
				SKU:         sku,
				Name:        data.name,
				Description: data.name,
				Slug:        data.slug,
				Price:       data.price,
				Barcode:     data.barcode,
				CategoryID:  categoryID,
			}

			products = append(products, product)
		}

		_, err = tx.NewInsert().
			Model(&products).
			On("CONFLICT (sku) DO UPDATE").
			Set("name = EXCLUDED.name").
			Set("price = EXCLUDED.price").
			Returning("id", "sku").
			Exec(ctx)
		if err != nil {
			return err
		}

		var inventories []Inventory
		var stock_movements []StockMovement

		for _, p := range products {
			if _, ok := skuMap[p.SKU]; ok {
				qty := skuMap[p.SKU].quantity

				inventory := Inventory{
					ProductID:         p.ID,
					Quantity:          qty,
					LowStockThreshold: 5,
				}

				inventories = append(inventories, inventory)

				if qty > 0 {
					stock_movement := StockMovement{
						ProductID:    p.ID,
						ChangeAmount: qty,
						Reason:       "FROM_IMPORT",
					}

					stock_movements = append(stock_movements, stock_movement)
				}
			}
		}

		_, err = tx.NewInsert().
			Model(&inventories).
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
