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
		Limit(q.PageSize).
		Offset(q.Page).
		Order(q.Sort).
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (r *ProductRepository) GetProductsPublic(ctx context.Context, q dto.BaseFilters) ([]dto.ProductPublicResponse, int, error) {
	var products []Product

	count, err := r.db.NewSelect().
		Model(&products).
		Relation("Category").
		Limit(q.PageSize).
		Offset(q.Page).
		Order(q.Sort).
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	items := make([]dto.ProductPublicResponse, count)
	for i, p := range products {
		items[i] = dto.ProductPublicResponse{
			Name:         p.Name,
			Slug:         p.Slug,
			Description:  p.Description,
			Barcode:      p.Barcode,
			DisplayPrice: float64(p.Price) / 100.00,
			Category:     &dto.CategoryPublicResponse{},
		}

		if p.Category != nil {
			items[i].Category = &dto.CategoryPublicResponse{
				Name: p.Category.Name,
				Slug: p.Category.Slug,
				UUID: p.Category.UUID,
			}
		}
	}

	return items, count, nil
}

func (r *ProductRepository) GetProductsBackoffice(ctx context.Context, q dto.BaseFilters) ([]dto.ProductBackofficeResponse, int, error) {
	var products []Product

	count, err := r.db.NewSelect().
		Model(&products).
		Relation("Category").
		Limit(q.PageSize).
		Offset(q.Page).
		Order(q.Sort).
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	items := make([]dto.ProductBackofficeResponse, count)
	for i, p := range products {
		items[i] = dto.ProductBackofficeResponse{
			Name:         p.Name,
			Slug:         p.Slug,
			Description:  p.Description,
			SKU:          p.SKU,
			Barcode:      p.Barcode,
			DisplayPrice: float64(p.Price) / 100.00,
			Category:     &dto.CategoryPublicResponse{},
		}

		if p.Category != nil {
			items[i].Category = &dto.CategoryPublicResponse{
				Name: p.Category.Name,
				Slug: p.Category.Slug,
				UUID: p.Category.UUID,
			}
		}
	}

	return items, count, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *Product) error {
	res, err := r.db.NewInsert().Model(product).Exec(ctx)

	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return err
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
