package service

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api/internal/data"
	"github.com/XaiPhyr/rdev-go-api/internal/dto"

	"github.com/redis/go-redis/v9"
)

type ProductService struct {
	r     *data.ProductRepository
	es    *EmailService
	redis *redis.Client
}

func NewProduct(r *data.ProductRepository, es *EmailService, redis *redis.Client) *ProductService {
	return &ProductService{r: r, es: es, redis: redis}
}

func (s *ProductService) GetProductByUUID(ctx context.Context, uuid string) (*data.Product, error) {
	return s.r.GetProductByUUID(ctx, uuid)
}

func (s *ProductService) GetProducts(ctx context.Context, q dto.Query) ([]data.Product, int, error) {
	filters := q.SanitizeQuery([]string{"name", "barcode"})

	return s.r.GetProducts(ctx, filters)
}

func (s *ProductService) GetProductsPublic(ctx context.Context, q dto.Query) ([]dto.ProductPublicResponse, int, error) {
	filters := q.SanitizeQuery([]string{"name", "barcode"})

	return s.r.GetProductsPublic(ctx, filters)
}

func (s *ProductService) GetProductsBackoffice(ctx context.Context, q dto.Query) ([]dto.ProductBackofficeResponse, int, error) {
	filters := q.SanitizeQuery([]string{"name", "barcode"})

	return s.r.GetProductsBackoffice(ctx, filters)
}

func (s *ProductService) UpdateProduct(ctx context.Context, uuid string, req dto.ProductRequestUpdate) error {
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

func (s *ProductService) DeleteProduct(ctx context.Context, uuid string) error {
	return s.r.DeleteProduct(ctx, uuid)
}

func (s *ProductService) UpdateProductStatus(ctx context.Context, uuid string) error {
	return s.r.UpdateProductStatus(ctx, uuid)
}
