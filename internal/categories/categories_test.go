package categories_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/XaiPhyr/rdev-go-api/internal/categories"
	"github.com/XaiPhyr/rdev-go-api/internal/mocks"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
)

const UUID = "12345678-1234-5678-1234-567890123456"

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
		category.ID = 1
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

func TestCategory(t *testing.T) {
	testCategoryRepo := &CategoryTest{}
	emailSvc := mocks.NewTestEmailService()
	_, auditLogSvc := mocks.NewTestAuditService()

	testCategorySvc := categories.NewCategoryService(testCategoryRepo, emailSvc, nil, auditLogSvc)

	t.Run("Get Categories", func(t *testing.T) {
		testCategoryRepo.GetCategoriesFunc = func(ctx context.Context, q dto.BaseFilters) ([]models.Category, int, error) {
			CheckCategoryQuery(t, q)
			return []models.Category{{Name: "Test Category"}}, 1, nil
		}

		query := dto.Query{Search: "test", Limit: 10, Offset: 2, Sort: "name ASC"}
		_, _, err := testCategorySvc.GetCategories(context.Background(), query)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Get Category By UUID", func(t *testing.T) {
		testCategoryRepo.GetCategoryByUUIDFunc = func(ctx context.Context, uuid string) (*models.Category, error) {
			CheckUUID(t, uuid)
			return &models.Category{Name: "Test Category"}, nil
		}

		_, err := testCategorySvc.GetCategoryByUUID(context.Background(), UUID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Get Category Tree", func(t *testing.T) {
		testCategoryRepo.GetCategoryTreeFunc = func(ctx context.Context, q dto.BaseFilters) ([]categories.CategoryTree, error) {
			CheckCategoryQuery(t, q)
			return []categories.CategoryTree{{Name: "Test Category"}}, nil
		}

		query := dto.Query{Search: "test", Limit: 10, Offset: 2, Sort: "name ASC"}
		_, err := testCategorySvc.GetCategoryTree(context.Background(), query)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Create Category", func(t *testing.T) {
		testCategoryRepo.CreateCategoryFunc = func(ctx context.Context, category *models.Category) error {
			CheckCategory(t, category)
			return nil
		}

		numRequest := 50
		var wg sync.WaitGroup

		for i := range numRequest {
			wg.Go(func() {
				name := fmt.Sprintf("test-%d", i)
				slug := fmt.Sprintf("slug-test-%d", i)
				req := categories.CategoryRequest{Name: &name, Slug: &slug}
				err := testCategorySvc.CreateCategory(context.Background(), req, models.AuditLogRequest{})

				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			})
		}

		wg.Wait()
	})

	t.Run("Update Category", func(t *testing.T) {
		testCategoryRepo.UpdateCategoryFunc = func(ctx context.Context, category *models.Category) error {
			if category.ID == 0 {
				t.Error("Expected category ID to be populated")
			}
			CheckCategory(t, category)
			return nil
		}

		name := "test"
		slug := "slug-test"
		req := categories.CategoryRequest{Name: &name, Slug: &slug}
		err := testCategorySvc.UpdateCategory(context.Background(), CheckUUID(t, UUID), req, models.AuditLogRequest{})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Delete Category", func(t *testing.T) {
		testCategoryRepo.DeleteCategoryFunc = func(ctx context.Context, uuid string) error {
			CheckUUID(t, uuid)
			return nil
		}

		err := testCategorySvc.DeleteCategory(context.Background(), CheckUUID(t, UUID), models.AuditLogRequest{})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Update Category Status", func(t *testing.T) {
		testCategoryRepo.UpdateCategoryStatusFunc = func(ctx context.Context, uuid string) error {
			CheckUUID(t, uuid)
			return nil
		}

		err := testCategorySvc.UpdateCategoryStatus(context.Background(), CheckUUID(t, UUID), models.AuditLogRequest{})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func CheckUUID(t testing.TB, uuid string) string {
	t.Helper()

	if uuid == "" {
		t.Error("Expected UUID to be provided")
	}

	return uuid
}

func CheckCategory(t testing.TB, category *models.Category) {
	t.Helper()

	if category.Name == "" {
		t.Error("Expected category name to be populated")
	}
	if category.Slug == "" {
		t.Error("Expected category slug to be populated")
	}
}

func CheckCategoryQuery(t testing.TB, q dto.BaseFilters) {
	t.Helper()

	if q.Search != "test" {
		t.Errorf("Expected search filter to be 'test', got '%s'", q.Search)
	}
	if q.Page < 1 {
		t.Errorf("Expected page to be 1, got %d", q.Page)
	}
	if q.PageSize < 1 {
		t.Errorf("Expected page size to be 10, got %d", q.PageSize)
	}
	if q.Sort != "name ASC" {
		t.Errorf("Expected sort to be name ASC, got '%s'", q.Sort)
	}
}
