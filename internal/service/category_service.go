package service

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api/internal/data"
	"github.com/XaiPhyr/rdev-go-api/internal/dto"

	"github.com/redis/go-redis/v9"
)

type CategoryRepository interface {
	GetCategoryByUUID(ctx context.Context, uuid string) (*data.Category, error)
	GetCategories(ctx context.Context, filters dto.BaseFilters) ([]data.Category, int, error)
	GetCategoryTree(ctx context.Context, filters dto.BaseFilters) ([]data.CategoryTree, error)
	CreateCategory(ctx context.Context, category *data.Category) error
	UpdateCategory(ctx context.Context, category *data.Category) error
	DeleteCategory(ctx context.Context, uuid string) error
	UpdateCategoryStatus(ctx context.Context, uuid string) error
}

type CategoryService interface {
	GetCategoryByUUID(ctx context.Context, uuid string) (*data.Category, error)
	GetCategories(ctx context.Context, q dto.Query) ([]data.Category, int, error)
	GetCategoryTree(ctx context.Context, q dto.Query) ([]data.CategoryTree, error)
	CreateCategory(ctx context.Context, req dto.CategoryRequest) error
	UpdateCategory(ctx context.Context, uuid string, req dto.CategoryRequest) error
	DeleteCategory(ctx context.Context, uuid string) error
	UpdateCategoryStatus(ctx context.Context, uuid string) error
}

type categoryService struct {
	r     CategoryRepository
	es    *EmailService
	redis *redis.Client
}

func NewCategoryService(r CategoryRepository, es *EmailService, redis *redis.Client) *categoryService {
	return &categoryService{r: r, es: es, redis: redis}
}

func (s *categoryService) GetCategoryByUUID(ctx context.Context, uuid string) (*data.Category, error) {
	return s.r.GetCategoryByUUID(ctx, uuid)
}

func (s *categoryService) GetCategories(ctx context.Context, q dto.Query) ([]data.Category, int, error) {
	filters := q.SanitizeQuery([]string{"name"})

	return s.r.GetCategories(ctx, filters)
}

func (s *categoryService) GetCategoryTree(ctx context.Context, q dto.Query) ([]data.CategoryTree, error) {
	filters := q.SanitizeQuery([]string{"name"})

	return s.r.GetCategoryTree(ctx, filters)
}

func (s *categoryService) CreateCategory(ctx context.Context, req dto.CategoryRequest) error {
	category := &data.Category{}

	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}

	return s.r.CreateCategory(ctx, category)
}

func (s *categoryService) UpdateCategory(ctx context.Context, uuid string, req dto.CategoryRequest) error {
	category, err := s.r.GetCategoryByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}

	return s.r.UpdateCategory(ctx, category)
}

func (s *categoryService) DeleteCategory(ctx context.Context, uuid string) error {
	return s.r.DeleteCategory(ctx, uuid)
}

func (s *categoryService) UpdateCategoryStatus(ctx context.Context, uuid string) error {
	return s.r.UpdateCategoryStatus(ctx, uuid)
}
