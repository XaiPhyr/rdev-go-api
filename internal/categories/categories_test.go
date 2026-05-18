package categories_test

import (
	"context"
	"testing"

	"github.com/XaiPhyr/rdev-go-api/internal/categories"
	"github.com/XaiPhyr/rdev-go-api/internal/mocks"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
)

type CategoryTest struct {
	GetCategoryByUUIDFunc    func(ctx context.Context, uuid string) (*models.Category, error)
	GetCategoriesFunc        func(ctx context.Context, q dto.BaseFilters) ([]models.Category, int, error)
	GetCategoryTreeFunc      func(ctx context.Context, q dto.BaseFilters) ([]categories.CategoryTree, error)
	CreateCategoryFunc       func(ctx context.Context, category *models.Category) error
	UpdateCategoryFunc       func(ctx context.Context, category *models.Category) error
	DeleteCategoryFunc       func(ctx context.Context, uuid string) error
	UpdateCategoryStatusFunc func(ctx context.Context, uuid string) error
}

func (m *CategoryTest) CreateCategory(ctx context.Context, category *models.Category) error {
	if m.CreateCategoryFunc != nil {
		return m.CreateCategoryFunc(ctx, category)
	}

	return nil
}
func (m *CategoryTest) GetCategoryByUUID(ctx context.Context, uuid string) (*models.Category, error) {
	if m.GetCategoryByUUIDFunc != nil {
		return m.GetCategoryByUUIDFunc(ctx, uuid)
	}

	return nil, nil
}
func (m *CategoryTest) GetCategories(ctx context.Context, q dto.BaseFilters) ([]models.Category, int, error) {
	if m.GetCategoriesFunc != nil {
		return m.GetCategoriesFunc(ctx, q)
	}

	return nil, 0, nil
}
func (m *CategoryTest) GetCategoryTree(ctx context.Context, q dto.BaseFilters) ([]categories.CategoryTree, error) {
	if m.GetCategoryTreeFunc != nil {
		return m.GetCategoryTreeFunc(ctx, q)
	}

	return nil, nil
}
func (m *CategoryTest) UpdateCategory(ctx context.Context, category *models.Category) error {
	if m.UpdateCategoryFunc != nil {
		return m.UpdateCategoryFunc(ctx, category)
	}
	return nil
}
func (m *CategoryTest) DeleteCategory(ctx context.Context, uuid string) error {
	if m.DeleteCategoryFunc != nil {
		return m.DeleteCategoryFunc(ctx, uuid)
	}

	return nil
}
func (m *CategoryTest) UpdateCategoryStatus(ctx context.Context, uuid string) error {
	if m.UpdateCategoryStatusFunc != nil {
		return m.UpdateCategoryStatusFunc(ctx, uuid)
	}

	return nil
}

func TestGetCategories(t *testing.T) {
	testCategoryRepo := &CategoryTest{}
	emailSvc := mocks.NewTestEmailService()
	_, auditLogSvc := mocks.NewTestAuditService()

	testCategorySvc := categories.NewCategoryService(testCategoryRepo, emailSvc, nil, auditLogSvc)

	testCategoryRepo.GetCategoriesFunc = func(ctx context.Context, q dto.BaseFilters) ([]models.Category, int, error) {
		if q.Search != "test" {
			t.Errorf("Expected search filter to be 'test', got '%s'", q.Search)
		}
		return []models.Category{{Name: "Test Category"}}, 1, nil
	}

	query := dto.Query{Search: "test"}
	_, _, err := testCategorySvc.GetCategories(context.Background(), query)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCreateCategory(t *testing.T) {
	testCategoryRepo := &CategoryTest{}
	emailSvc := mocks.NewTestEmailService()
	_, auditLogSvc := mocks.NewTestAuditService()

	testCategorySvc := categories.NewCategoryService(testCategoryRepo, emailSvc, nil, auditLogSvc)

	testCategoryRepo.CreateCategoryFunc = func(ctx context.Context, category *models.Category) error {
		if category.Name == "" {
			t.Error("Expected category name to be populated")
		}
		return nil
	}

	name := "1"
	req := categories.CategoryRequest{Name: &name}
	err := testCategorySvc.CreateCategory(context.Background(), req, models.AuditLogRequest{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
