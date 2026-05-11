package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/XaiPhyr/rdev-go-api/internal/dto"
	"github.com/uptrace/bun"
)

type ProductRepository struct {
	db *bun.DB
}

func NewProductRepository(db *bun.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProductByUUID(ctx context.Context, uuid string) (*Product, error) {
	product := new(Product)

	err := r.db.NewSelect().
		Model(product).
		Where("uuid = ?", uuid).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) GetProducts(ctx context.Context, q dto.BaseFilters) ([]Product, int, error) {
	var products []Product

	count, err := r.db.NewSelect().
		Model(&products).
		Relation("Category").
		Relation("Inventory").
		Relation("StockMovement").
		Limit(q.PageSize).
		Offset(q.Page).
		Order(q.Sort).
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (r *ProductRepository) GetProductsPublic(ctx context.Context, q dto.BaseFilters) ([]Product, int, error) {
	var products []Product

	count, err := r.db.NewSelect().
		Model(&products).
		Relation("Category").
		Relation("Inventory").
		Limit(q.PageSize).
		Offset(q.Page).
		Order(q.Sort).
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (r *ProductRepository) GetProductsBackoffice(ctx context.Context, q dto.BaseFilters) ([]Product, int, error) {
	var products []Product

	count, err := r.db.NewSelect().
		Model(&products).
		Relation("Category").
		Relation("Inventory").
		Limit(q.PageSize).
		Offset(q.Page).
		Order(q.Sort).
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *Product, initQty int64) error {
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(product).Exec(ctx); err != nil {
			return err
		}

		inventory := &Inventory{
			ProductID:         product.ID,
			Quantity:          initQty,
			LowStockThreshold: 5,
		}

		if _, err := tx.NewInsert().Model(inventory).Exec(ctx); err != nil {
			return err
		}

		stock_movement := &StockMovement{
			ProductID:    product.ID,
			ChangeAmount: initQty,
			Reason:       "INITIAL_STOCKS",
		}

		if _, err := tx.NewInsert().Model(stock_movement).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product *Product) error {
	res, err := r.db.NewUpdate().
		Model(product).
		Column("category_id", "name", "slug", "description", "sku", "barcode", "price", "cost_price").
		Set("updated_at = ?", time.Now()).
		WherePK().
		Exec(ctx)

	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return err
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, uuid string) error {
	res, err := r.db.NewDelete().
		Model((*Product)(nil)).
		Where("uuid = ?", uuid).
		Exec(ctx)

	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return err
}

func (r *ProductRepository) UpdateProductStatus(ctx context.Context, uuid string) error {
	res, err := r.db.NewUpdate().
		Model((*Product)(nil)).
		Set("status = CASE WHEN status = 'A' THEN 'I' ELSE 'A' END").
		Set("updated_at = ?", time.Now()).
		Where("uuid = ?", uuid).
		Exec(ctx)

	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return err
}
