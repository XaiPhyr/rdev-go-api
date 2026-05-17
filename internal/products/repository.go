package products

import (
	"context"
	"time"

	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewProductRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetProductByUUID(ctx context.Context, uuid string) (*models.Product, error) {
	product := new(models.Product)

	err := r.db.NewSelect().
		Model(product).
		Where("uuid = ?", uuid).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *Repository) GetProducts(ctx context.Context, q dto.BaseFilters) ([]models.Product, int, error) {
	var products []models.Product

	count, err := r.db.NewSelect().
		Model(&products).
		Relation("Category").
		Relation("Inventory").
		Relation("StockMovement").
		Limit(q.PageSize).
		Offset(q.Page).
		Order(q.Sort).
		WhereGroup(" OR ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			search := "%" + q.Search + "%"

			return sq.WhereOr("p.name ILIKE ?", search).
				WhereOr("p.slug ILIKE ?", search)
		}).
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (r *Repository) GetProductsPublic(ctx context.Context, q dto.BaseFilters) ([]models.Product, int, error) {
	var products []models.Product

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

func (r *Repository) GetProductsBackoffice(ctx context.Context, q dto.BaseFilters) ([]models.Product, int, error) {
	var products []models.Product

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

func (r *Repository) CreateProduct(ctx context.Context, product *models.Product, initQty int64) error {
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(product).Exec(ctx); err != nil {
			return err
		}

		inventory := &models.Inventory{
			ProductID:         product.ID,
			Quantity:          initQty,
			LowStockThreshold: 5,
		}

		if _, err := tx.NewInsert().Model(inventory).Exec(ctx); err != nil {
			return err
		}

		stock_movement := &models.StockMovement{
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

func (r *Repository) UpdateProduct(ctx context.Context, product *models.Product) error {
	_, err := r.db.NewUpdate().
		Model(product).
		Column("category_id", "name", "slug", "description", "sku", "barcode", "price", "cost_price").
		Set("updated_at = ?", time.Now()).
		WherePK().
		Exec(ctx)

	return err
}

func (r *Repository) DeleteProduct(ctx context.Context, uuid string) error {
	_, err := r.db.NewDelete().
		Model((*models.Product)(nil)).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}

func (r *Repository) UpdateProductStatus(ctx context.Context, uuid string) error {
	_, err := r.db.NewUpdate().
		Model((*models.Product)(nil)).
		Set("status = CASE WHEN status = 'A' THEN 'I' ELSE 'A' END").
		Set("updated_at = ?", time.Now()).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}
