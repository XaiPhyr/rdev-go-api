package products_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/XaiPhyr/rdev-go-api/internal/mocks"
	"github.com/XaiPhyr/rdev-go-api/internal/products"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
)

const UUID = "12345678-1234-5678-1234-567890123456"

type ProductTest struct {
	GetProductByUUIDFunc      func(ctx context.Context, uuid string) (*models.Product, error)
	GetProductsFunc           func(ctx context.Context, filters dto.BaseFilters) ([]models.Product, int, error)
	GetProductsPublicFunc     func(ctx context.Context, q dto.BaseFilters) ([]models.Product, int, error)
	GetProductsBackofficeFunc func(ctx context.Context, q dto.BaseFilters) ([]models.Product, int, error)
	CreateProductFunc         func(ctx context.Context, product *models.Product, initQty int64) error
	UpdateProductFunc         func(ctx context.Context, product *models.Product) error
	DeleteProductFunc         func(ctx context.Context, uuid string) error
	UpdateProductStatusFunc   func(ctx context.Context, uuid string) error
}

func (m *ProductTest) GetProductByUUID(ctx context.Context, uuid string) (*models.Product, error) {
	if m.GetProductByUUIDFunc != nil {
		return m.GetProductByUUIDFunc(ctx, uuid)
	}

	return nil, nil
}
func (m *ProductTest) GetProducts(ctx context.Context, q dto.BaseFilters) ([]models.Product, int, error) {
	if m.GetProductsFunc != nil {
		return m.GetProductsFunc(ctx, q)
	}

	return nil, 0, nil
}
func (m *ProductTest) GetProductsBackoffice(ctx context.Context, q dto.BaseFilters) ([]models.Product, int, error) {
	if m.GetProductsBackofficeFunc != nil {
		return m.GetProductsBackofficeFunc(ctx, q)
	}

	return nil, 0, nil
}
func (m *ProductTest) GetProductsPublic(ctx context.Context, q dto.BaseFilters) ([]models.Product, int, error) {
	if m.GetProductsPublicFunc != nil {
		return m.GetProductsPublicFunc(ctx, q)
	}

	return nil, 0, nil
}

func (m *ProductTest) CreateProduct(ctx context.Context, product *models.Product, initQty int64) error {
	if m.CreateProductFunc != nil {
		return m.CreateProductFunc(ctx, product, initQty)
	}

	return nil
}
func (m *ProductTest) UpdateProduct(ctx context.Context, category *models.Product) error {
	if m.UpdateProductFunc != nil {
		category.ID = 1
		return m.UpdateProductFunc(ctx, category)
	}
	return nil
}
func (m *ProductTest) DeleteProduct(ctx context.Context, uuid string) error {
	if m.DeleteProductFunc != nil {
		return m.DeleteProductFunc(ctx, uuid)
	}

	return nil
}
func (m *ProductTest) UpdateProductStatus(ctx context.Context, uuid string) error {
	if m.UpdateProductStatusFunc != nil {
		return m.UpdateProductStatusFunc(ctx, uuid)
	}

	return nil
}

func TestProduct(t *testing.T) {
	testProductRepo := &ProductTest{}
	emailSvc := mocks.NewTestEmailService()
	_, auditLogSvc := mocks.NewTestAuditService()

	testProductSvc := products.NewProductService(testProductRepo, emailSvc, nil, auditLogSvc)

	t.Run("Get Products", func(t *testing.T) {
		testProductRepo.GetProductsFunc = func(ctx context.Context, q dto.BaseFilters) ([]models.Product, int, error) {
			CheckProductQuery(t, q)
			return []models.Product{{Name: "Test Product"}}, 1, nil
		}

		query := dto.Query{Search: "test", Limit: 10, Offset: 2, Sort: "name ASC"}
		_, _, err := testProductSvc.GetProducts(context.Background(), query)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Get Product By UUID", func(t *testing.T) {
		testProductRepo.GetProductByUUIDFunc = func(ctx context.Context, uuid string) (*models.Product, error) {
			CheckUUID(t, uuid)
			return &models.Product{Name: "Test Product"}, nil
		}

		_, err := testProductSvc.GetProductByUUID(context.Background(), UUID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Create Product", func(t *testing.T) {
		testProductRepo.CreateProductFunc = func(ctx context.Context, category *models.Product, initQty int64) error {
			CheckProduct(t, category)
			return nil
		}

		numRequest := 50
		var wg sync.WaitGroup

		for i := range numRequest {
			wg.Go(func() {
				name := fmt.Sprintf("test-%d", i)
				slug := fmt.Sprintf("slug-test-%d", i)
				req := products.ProductRequest{Name: &name, Slug: &slug}
				err := testProductSvc.CreateProduct(context.Background(), req, models.AuditLogRequest{})

				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			})
		}

		wg.Wait()
	})

	t.Run("Update Product", func(t *testing.T) {
		testProductRepo.UpdateProductFunc = func(ctx context.Context, category *models.Product) error {
			if category.ID == 0 {
				t.Error("Expected category ID to be populated")
			}
			CheckProduct(t, category)
			return nil
		}

		name := "test"
		slug := "slug-test"
		req := products.ProductRequest{Name: &name, Slug: &slug}
		err := testProductSvc.UpdateProduct(context.Background(), CheckUUID(t, UUID), req, models.AuditLogRequest{})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Delete Product", func(t *testing.T) {
		testProductRepo.DeleteProductFunc = func(ctx context.Context, uuid string) error {
			CheckUUID(t, uuid)
			return nil
		}

		err := testProductSvc.DeleteProduct(context.Background(), CheckUUID(t, UUID), models.AuditLogRequest{})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Update Product Status", func(t *testing.T) {
		testProductRepo.UpdateProductStatusFunc = func(ctx context.Context, uuid string) error {
			CheckUUID(t, uuid)
			return nil
		}

		err := testProductSvc.UpdateProductStatus(context.Background(), CheckUUID(t, UUID), models.AuditLogRequest{})
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

func CheckProduct(t testing.TB, category *models.Product) {
	t.Helper()

	if category.Name == "" {
		t.Error("Expected category name to be populated")
	}
	if category.Slug == "" {
		t.Error("Expected category slug to be populated")
	}
}

func CheckProductQuery(t testing.TB, q dto.BaseFilters) {
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
