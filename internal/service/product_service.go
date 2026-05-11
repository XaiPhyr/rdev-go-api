package service

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api/internal/data"
	"github.com/XaiPhyr/rdev-go-api/internal/dto"

	"github.com/redis/go-redis/v9"
)

type ProductRepository interface {
	GetProductByUUID(ctx context.Context, uuid string) (*data.Product, error)
	GetProducts(ctx context.Context, filters dto.BaseFilters) ([]data.Product, int, error)
	GetProductsPublic(ctx context.Context, q dto.BaseFilters) ([]dto.ProductPublicResponse, int, error)
	GetProductsBackoffice(ctx context.Context, q dto.BaseFilters) ([]dto.ProductBackofficeResponse, int, error)
	CreateProduct(ctx context.Context, category *data.Product, initQty int64) error
	UpdateProduct(ctx context.Context, category *data.Product) error
	DeleteProduct(ctx context.Context, uuid string) error
	UpdateProductStatus(ctx context.Context, uuid string) error
}

type ProductService interface {
	GetProductByUUID(ctx context.Context, uuid string) (*data.Product, error)
	GetProducts(ctx context.Context, q dto.Query) ([]data.Product, int, error)
	GetProductsPublic(ctx context.Context, q dto.Query) ([]dto.ProductPublicResponse, int, error)
	GetProductsBackoffice(ctx context.Context, q dto.Query) ([]dto.ProductBackofficeResponse, int, error)
	CreateProduct(ctx context.Context, req dto.ProductRequest) error
	UpdateProduct(ctx context.Context, uuid string, req dto.ProductRequest) error
	DeleteProduct(ctx context.Context, uuid string) error
	UpdateProductStatus(ctx context.Context, uuid string) error
}

type productService struct {
	r     ProductRepository
	es    *EmailService
	redis *redis.Client
}

func NewProductService(r ProductRepository, es *EmailService, redis *redis.Client) *productService {
	return &productService{r: r, es: es, redis: redis}
}

func (s *productService) GetProductByUUID(ctx context.Context, uuid string) (*data.Product, error) {
	return s.r.GetProductByUUID(ctx, uuid)
}

func (s *productService) GetProducts(ctx context.Context, q dto.Query) ([]data.Product, int, error) {
	filters := q.SanitizeQuery([]string{"name", "barcode"})

	return s.r.GetProducts(ctx, filters)
}

func (s *productService) GetProductsPublic(ctx context.Context, q dto.Query) ([]dto.ProductPublicResponse, int, error) {
	filters := q.SanitizeQuery([]string{"name", "barcode"})

	return s.r.GetProductsPublic(ctx, filters)
}

func (s *productService) GetProductsBackoffice(ctx context.Context, q dto.Query) ([]dto.ProductBackofficeResponse, int, error) {
	filters := q.SanitizeQuery([]string{"name", "barcode"})

	return s.r.GetProductsBackoffice(ctx, filters)
}

func (s *productService) CreateProduct(ctx context.Context, req dto.ProductRequest) error {
	product := &data.Product{}
	var qty int64 = 0

	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Slug != nil {
		product.Slug = *req.Slug
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.SKU != nil {
		product.SKU = *req.SKU
	}
	if req.Barcode != nil {
		product.Barcode = *req.Barcode
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.CostPrice != nil {
		product.CostPrice = *req.CostPrice
	}
	if req.Quantity != nil {
		qty = *req.Quantity
	}

	return s.r.CreateProduct(ctx, product, qty)
}

func (s *productService) UpdateProduct(ctx context.Context, uuid string, req dto.ProductRequest) error {
	product, err := s.r.GetProductByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Slug != nil {
		product.Slug = *req.Slug
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.SKU != nil {
		product.SKU = *req.SKU
	}
	if req.Barcode != nil {
		product.Barcode = *req.Barcode
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.CostPrice != nil {
		product.CostPrice = *req.CostPrice
	}

	return s.r.UpdateProduct(ctx, product)
}

func (s *productService) DeleteProduct(ctx context.Context, uuid string) error {
	return s.r.DeleteProduct(ctx, uuid)
}

func (s *productService) UpdateProductStatus(ctx context.Context, uuid string) error {
	return s.r.UpdateProductStatus(ctx, uuid)
}
