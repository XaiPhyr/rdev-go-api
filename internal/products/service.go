package products

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api/internal/audit_logs"
	"github.com/XaiPhyr/rdev-go-api/internal/categories"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/email"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"

	"github.com/redis/go-redis/v9"
)

type ProductRepository interface {
	GetProductByUUID(ctx context.Context, uuid string) (*Product, error)
	GetProducts(ctx context.Context, filters dto.BaseFilters) ([]Product, int, error)
	GetProductsPublic(ctx context.Context, q dto.BaseFilters) ([]Product, int, error)
	GetProductsBackoffice(ctx context.Context, q dto.BaseFilters) ([]Product, int, error)
	CreateProduct(ctx context.Context, category *Product, initQty int64) error
	UpdateProduct(ctx context.Context, category *Product) error
	DeleteProduct(ctx context.Context, uuid string) error
	UpdateProductStatus(ctx context.Context, uuid string) error
}

type ProductService interface {
	GetProductByUUID(ctx context.Context, uuid string) (*Product, error)
	GetProducts(ctx context.Context, q dto.Query) ([]Product, int, error)
	GetProductsPublic(ctx context.Context, q dto.Query) ([]ProductPublicResponse, int, error)
	GetProductsBackoffice(ctx context.Context, q dto.Query) ([]ProductBackofficeResponse, int, error)
	CreateProduct(ctx context.Context, req ProductRequest, audit models.AuditLogRequest) error
	UpdateProduct(ctx context.Context, uuid string, req ProductRequest, audit models.AuditLogRequest) error
	DeleteProduct(ctx context.Context, uuid string, audit models.AuditLogRequest) error
	UpdateProductStatus(ctx context.Context, uuid string, audit models.AuditLogRequest) error
}

type service struct {
	r        ProductRepository
	es       *email.EmailService
	redis    *redis.Client
	auditLog audit_logs.AuditLogService
}

func NewProductService(r ProductRepository, es *email.EmailService, redis *redis.Client, auditLog audit_logs.AuditLogService) *service {
	return &service{r: r, es: es, redis: redis, auditLog: auditLog}
}

func (s *service) GetProductByUUID(ctx context.Context, uuid string) (*Product, error) {
	return s.r.GetProductByUUID(ctx, uuid)
}

func (s *service) GetProducts(ctx context.Context, q dto.Query) ([]Product, int, error) {
	filters := q.SanitizeQuery([]string{"name", "barcode"})

	return s.r.GetProducts(ctx, filters)
}

func (s *service) GetProductsPublic(ctx context.Context, q dto.Query) ([]ProductPublicResponse, int, error) {
	filters := q.SanitizeQuery([]string{"name", "barcode"})

	products, count, err := s.r.GetProductsPublic(ctx, filters)
	items := make([]ProductPublicResponse, len(products))
	for i, p := range products {
		items[i] = ProductPublicResponse{
			Name:         p.Name,
			Slug:         p.Slug,
			Description:  p.Description,
			Barcode:      p.Barcode,
			DisplayPrice: float64(p.Price) / 100.00,
			Category:     &categories.CategoryResponse{},
		}

		if p.Category != nil {
			items[i].Category = &categories.CategoryResponse{
				Name: p.Category.Name,
				Slug: p.Category.Slug,
				UUID: p.Category.UUID,
			}
		}
	}

	return items, count, err
}

func (s *service) GetProductsBackoffice(ctx context.Context, q dto.Query) ([]ProductBackofficeResponse, int, error) {
	filters := q.SanitizeQuery([]string{"name", "barcode"})

	products, count, err := s.r.GetProductsBackoffice(ctx, filters)

	items := make([]ProductBackofficeResponse, len(products))
	for i, p := range products {
		items[i] = ProductBackofficeResponse{
			Name:         p.Name,
			Slug:         p.Slug,
			Description:  p.Description,
			SKU:          p.SKU,
			Barcode:      p.Barcode,
			DisplayPrice: float64(p.Price) / 100.00,
			Category:     &categories.CategoryResponse{},
		}

		if p.Category != nil {
			items[i].Category = &categories.CategoryResponse{
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

func (s *service) CreateProduct(ctx context.Context, req ProductRequest, audit models.AuditLogRequest) error {
	product := &Product{}
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
	s.auditLog.CreateAuditLog(s.auditLog.ParseAuditLog(audit, product.UUID, "PRODUCT", nil, *product, err))

	return err
}

func (s *service) UpdateProduct(ctx context.Context, uuid string, req ProductRequest, audit models.AuditLogRequest) error {
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
	s.auditLog.CreateAuditLog(s.auditLog.ParseAuditLog(audit, uuid, "PRODUCT", oldProduct, *product, err))

	return err
}

func (s *service) DeleteProduct(ctx context.Context, uuid string, audit models.AuditLogRequest) error {
	product, err := s.r.GetProductByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.DeleteProduct(ctx, uuid)
	s.auditLog.CreateAuditLog(s.auditLog.ParseAuditLog(audit, uuid, "PRODUCT", nil, *product, err))

	return err
}

func (s *service) UpdateProductStatus(ctx context.Context, uuid string, audit models.AuditLogRequest) error {
	product, err := s.r.GetProductByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.UpdateProductStatus(ctx, uuid)
	s.auditLog.CreateAuditLog(s.auditLog.ParseAuditLog(audit, uuid, "PRODUCT", nil, *product, err))

	return err
}
