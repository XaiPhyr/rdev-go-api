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
	GetProductsPublic(ctx context.Context, q dto.BaseFilters) ([]data.Product, int, error)
	GetProductsBackoffice(ctx context.Context, q dto.BaseFilters) ([]data.Product, int, error)
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
	CreateProduct(ctx context.Context, req dto.ProductRequest, audit dto.AuditLogRequest) error
	UpdateProduct(ctx context.Context, uuid string, req dto.ProductRequest, audit dto.AuditLogRequest) error
	DeleteProduct(ctx context.Context, uuid string, audit dto.AuditLogRequest) error
	UpdateProductStatus(ctx context.Context, uuid string, audit dto.AuditLogRequest) error
}

type productService struct {
	r        ProductRepository
	es       *EmailService
	redis    *redis.Client
	auditLog AuditLogService
}

func NewProductService(r ProductRepository, es *EmailService, redis *redis.Client, auditLog AuditLogService) *productService {
	return &productService{r: r, es: es, redis: redis, auditLog: auditLog}
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

	products, count, err := s.r.GetProductsPublic(ctx, filters)
	items := make([]dto.ProductPublicResponse, len(products))
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

	return items, count, err
}

func (s *productService) GetProductsBackoffice(ctx context.Context, q dto.Query) ([]dto.ProductBackofficeResponse, int, error) {
	filters := q.SanitizeQuery([]string{"name", "barcode"})

	products, count, err := s.r.GetProductsBackoffice(ctx, filters)

	items := make([]dto.ProductBackofficeResponse, len(products))
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

		if p.Inventory != nil {
			items[i].Quantity = p.Inventory.Quantity
		}
	}

	return items, count, err
}

func (s *productService) CreateProduct(ctx context.Context, req dto.ProductRequest, audit dto.AuditLogRequest) error {
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

	err := s.r.CreateProduct(ctx, product, qty)
	s.auditLog.CreateAuditLog(parseAuditLog(audit, product.UUID, "PRODUCT", nil, *product, err))

	return err
}

func (s *productService) UpdateProduct(ctx context.Context, uuid string, req dto.ProductRequest, audit dto.AuditLogRequest) error {
	product, err := s.r.GetProductByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	oldProduct := *product

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

	err = s.r.UpdateProduct(ctx, product)
	s.auditLog.CreateAuditLog(parseAuditLog(audit, uuid, "PRODUCT", oldProduct, *product, err))

	return err
}

func (s *productService) DeleteProduct(ctx context.Context, uuid string, audit dto.AuditLogRequest) error {
	product, err := s.r.GetProductByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.DeleteProduct(ctx, uuid)
	s.auditLog.CreateAuditLog(parseAuditLog(audit, uuid, "PRODUCT", nil, *product, err))

	return err
}

func (s *productService) UpdateProductStatus(ctx context.Context, uuid string, audit dto.AuditLogRequest) error {
	product, err := s.r.GetProductByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.UpdateProductStatus(ctx, uuid)
	s.auditLog.CreateAuditLog(parseAuditLog(audit, uuid, "PRODUCT", nil, *product, err))

	return err
}
