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

func NewCategory(r *data.CategoryRepository, es *EmailService, redis *redis.Client) *CategoryService {
	return &CategoryService{r: r, es: es, redis: redis}
}

func (s *CategoryService) GetCategoryByUUID(ctx context.Context, uuid string) (*data.Category, error) {
	return s.r.GetCategoryByUUID(ctx, uuid)
}

func (s *CategoryService) GetCategories(ctx context.Context, q dto.Query) ([]data.Category, int, error) {
	sort := "id ASC"
	if q.Sort != "" {
		sort = q.Sort
	}

	filters := data.BaseFilters{
		PageSize: q.Limit,
		Page:     q.Offset,
		Sort:     sort,
		Search:   q.Search,
	}

	return s.r.GetCategories(ctx, filters)
}

func (s *CategoryService) GetCategoryTree(ctx context.Context, q dto.Query) ([]data.CategoryTree, error) {
	sort := "id ASC"
	if q.Sort != "" {
		sort = q.Sort
	}

	filters := data.BaseFilters{
		PageSize: q.Limit,
		Page:     q.Offset,
		Sort:     sort,
		Search:   q.Search,
	}

	return s.r.GetCategoryTree(ctx, filters)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, uuid string, req dto.CategoryRequestUpdate) error {
	category, err := s.r.GetCategoryByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	category.ParentID = req.ParentID
	category.Name = req.Name
	category.Slug = req.Slug

	return s.r.UpdateCategory(ctx, category)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, uuid string) error {
	return s.r.DeleteCategory(ctx, uuid)
}

func (s *CategoryService) UpdateCategoryStatus(ctx context.Context, uuid string) error {
	return s.r.UpdateCategoryStatus(ctx, uuid)
}
