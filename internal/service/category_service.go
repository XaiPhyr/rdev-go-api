package service

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api/internal/data"
	"github.com/XaiPhyr/rdev-go-api/internal/dto"

	"github.com/redis/go-redis/v9"
)

type CategoryService struct {
	r     *data.CategoryRepository
	es    *EmailService
	redis *redis.Client
}

func NewCategoryService(r *data.CategoryRepository, es *EmailService, redis *redis.Client) *CategoryService {
	return &CategoryService{r: r, es: es, redis: redis}
}

func (s *CategoryService) GetCategoryByUUID(ctx context.Context, uuid string) (*data.Category, error) {
	return s.r.GetCategoryByUUID(ctx, uuid)
}

func (s *CategoryService) GetCategories(ctx context.Context, q dto.Query) ([]data.Category, int, error) {
	filters := q.SanitizeQuery([]string{"name"})

	return s.r.GetCategories(ctx, filters)
}

func (s *CategoryService) GetCategoryTree(ctx context.Context, q dto.Query) ([]data.CategoryTree, error) {
	filters := q.SanitizeQuery([]string{"name"})

	return s.r.GetCategoryTree(ctx, filters)
}

func (s *CategoryService) CreateCategory(ctx context.Context, req dto.CategoryRequest) error {
	category := &data.Category{}

	if req.ParentID != nil {
		category.ParentID = *req.ParentID
	}
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}

	return s.r.CreateCategory(ctx, category)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, uuid string, req dto.CategoryRequest) error {
	category, err := s.r.GetCategoryByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	if req.ParentID != nil {
		category.ParentID = *req.ParentID
	}
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}

	return s.r.UpdateCategory(ctx, category)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, uuid string) error {
	return s.r.DeleteCategory(ctx, uuid)
}

func (s *CategoryService) UpdateCategoryStatus(ctx context.Context, uuid string) error {
	return s.r.UpdateCategoryStatus(ctx, uuid)
}
