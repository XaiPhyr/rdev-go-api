package categories

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api/internal/audit_logs"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/email"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"

	"github.com/redis/go-redis/v9"
)

type CategoryRepository interface {
	GetCategoryByUUID(ctx context.Context, uuid string) (*Category, error)
	GetCategories(ctx context.Context, filters dto.BaseFilters) ([]Category, int, error)
	GetCategoryTree(ctx context.Context, filters dto.BaseFilters) ([]CategoryTree, error)
	CreateCategory(ctx context.Context, category *Category) error
	UpdateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, uuid string) error
	UpdateCategoryStatus(ctx context.Context, uuid string) error
}

type CategoryService interface {
	GetCategoryByUUID(ctx context.Context, uuid string) (*Category, error)
	GetCategories(ctx context.Context, q dto.Query) ([]Category, int, error)
	GetCategoryTree(ctx context.Context, q dto.Query) ([]CategoryTree, error)
	CreateCategory(ctx context.Context, req CategoryRequest, auditLog models.AuditLogRequest) error
	UpdateCategory(ctx context.Context, uuid string, req CategoryRequest, auditLog models.AuditLogRequest) error
	DeleteCategory(ctx context.Context, uuid string, auditLog models.AuditLogRequest) error
	UpdateCategoryStatus(ctx context.Context, uuid string, auditLog models.AuditLogRequest) error
}

type service struct {
	r        CategoryRepository
	es       *email.EmailService
	redis    *redis.Client
	auditLog audit_logs.AuditLogService
}

func NewCategoryService(r CategoryRepository, es *email.EmailService, redis *redis.Client, auditLog audit_logs.AuditLogService) *service {
	return &service{r: r, es: es, redis: redis, auditLog: auditLog}
}

func (s *service) GetCategoryByUUID(ctx context.Context, uuid string) (*Category, error) {
	return s.r.GetCategoryByUUID(ctx, uuid)
}

func (s *service) GetCategories(ctx context.Context, q dto.Query) ([]Category, int, error) {
	filters := q.SanitizeQuery([]string{"name"})

	return s.r.GetCategories(ctx, filters)
}

func (s *service) GetCategoryTree(ctx context.Context, q dto.Query) ([]CategoryTree, error) {
	filters := q.SanitizeQuery([]string{"name"})

	return s.r.GetCategoryTree(ctx, filters)
}

func (s *service) CreateCategory(ctx context.Context, req CategoryRequest, audit models.AuditLogRequest) error {
	category := &Category{}

	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}

	err := s.r.CreateCategory(ctx, category)
	s.auditLog.CreateAuditLog(s.auditLog.ParseAuditLog(audit, category.UUID, "CATEGORY", nil, *category, err))

	return err
}

func (s *service) UpdateCategory(ctx context.Context, uuid string, req CategoryRequest, audit models.AuditLogRequest) error {
	category, err := s.r.GetCategoryByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	oldCategory := *category

	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}

	err = s.r.UpdateCategory(ctx, category)
	s.auditLog.CreateAuditLog(s.auditLog.ParseAuditLog(audit, uuid, "CATEGORY", oldCategory, *category, err))

	return err
}

func (s *service) DeleteCategory(ctx context.Context, uuid string, audit models.AuditLogRequest) error {
	category, err := s.r.GetCategoryByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.DeleteCategory(ctx, uuid)
	s.auditLog.CreateAuditLog(s.auditLog.ParseAuditLog(audit, uuid, "CATEGORY", nil, *category, err))

	return err
}

func (s *service) UpdateCategoryStatus(ctx context.Context, uuid string, audit models.AuditLogRequest) error {
	category, err := s.r.GetCategoryByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.UpdateCategoryStatus(ctx, uuid)
	s.auditLog.CreateAuditLog(s.auditLog.ParseAuditLog(audit, uuid, "CATEGORY", nil, *category, err))

	return err
}
