package inventories_test

import (
	"context"
	"sync"
	"testing"

	"github.com/XaiPhyr/rdev-go-api/internal/inventories"
	"github.com/XaiPhyr/rdev-go-api/internal/mocks"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
)

const UUID = "12345678-1234-5678-1234-567890123456"

type InventoryTest struct {
	GetInventoryByUUIDFunc    func(ctx context.Context, uuid string) (*models.Inventory, error)
	GetInventoriesFunc        func(ctx context.Context, filters dto.BaseFilters) ([]models.Inventory, int, error)
	CreateInventoryFunc       func(ctx context.Context, inventory *models.Inventory) error
	UpdateInventoryFunc       func(ctx context.Context, inventory *models.Inventory) error
	DeleteInventoryFunc       func(ctx context.Context, uuid string) error
	UpdateInventoryStatusFunc func(ctx context.Context, uuid string) error
}

func (m *InventoryTest) GetInventoryByUUID(ctx context.Context, uuid string) (*models.Inventory, error) {
	if m.GetInventoryByUUIDFunc != nil {
		return m.GetInventoryByUUIDFunc(ctx, uuid)
	}

	return nil, nil
}
func (m *InventoryTest) GetInventories(ctx context.Context, q dto.BaseFilters) ([]models.Inventory, int, error) {
	if m.GetInventoriesFunc != nil {
		return m.GetInventoriesFunc(ctx, q)
	}

	return nil, 0, nil
}
func (m *InventoryTest) CreateInventory(ctx context.Context, inventory *models.Inventory) error {
	if m.CreateInventoryFunc != nil {
		return m.CreateInventoryFunc(ctx, inventory)
	}

	return nil
}
func (m *InventoryTest) UpdateInventory(ctx context.Context, inventory *models.Inventory) error {
	if m.UpdateInventoryFunc != nil {
		inventory.ID = 1
		return m.UpdateInventoryFunc(ctx, inventory)
	}
	return nil
}
func (m *InventoryTest) DeleteInventory(ctx context.Context, uuid string) error {
	if m.DeleteInventoryFunc != nil {
		return m.DeleteInventoryFunc(ctx, uuid)
	}

	return nil
}
func (m *InventoryTest) UpdateInventoryStatus(ctx context.Context, uuid string) error {
	if m.UpdateInventoryStatusFunc != nil {
		return m.UpdateInventoryStatusFunc(ctx, uuid)
	}

	return nil
}

func TestInventory(t *testing.T) {
	testInventoryRepo := &InventoryTest{}
	emailSvc := mocks.NewTestEmailService()
	_, auditLogSvc := mocks.NewTestAuditService()

	testInventorySvc := inventories.NewInventoryService(testInventoryRepo, emailSvc, nil, auditLogSvc)

	t.Run("Get Inventories", func(t *testing.T) {
		testInventoryRepo.GetInventoriesFunc = func(ctx context.Context, q dto.BaseFilters) ([]models.Inventory, int, error) {
			CheckInventoryQuery(t, q)
			return []models.Inventory{{ProductID: 1}}, 1, nil
		}

		query := dto.Query{Search: "test", Limit: 10, Offset: 2, Sort: "quantity ASC"}
		_, _, err := testInventorySvc.GetInventories(context.Background(), query)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Get Inventory By UUID", func(t *testing.T) {
		testInventoryRepo.GetInventoryByUUIDFunc = func(ctx context.Context, uuid string) (*models.Inventory, error) {
			CheckUUID(t, uuid)
			return &models.Inventory{ProductID: 1}, nil
		}

		_, err := testInventorySvc.GetInventoryByUUID(context.Background(), UUID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Create Inventory", func(t *testing.T) {
		testInventoryRepo.CreateInventoryFunc = func(ctx context.Context, inventory *models.Inventory) error {
			CheckInventory(t, inventory)
			return nil
		}

		numRequest := 50
		var wg sync.WaitGroup

		for i := range numRequest {
			wg.Go(func() {
				product_id := int64(i + 1)
				quantity := int64(i + 1*2)
				req := inventories.InventoryRequest{ProductID: &product_id, Quantity: &quantity}
				err := testInventorySvc.CreateInventory(context.Background(), req, models.AuditLogRequest{})

				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			})
		}

		wg.Wait()
	})

	t.Run("Update Inventory", func(t *testing.T) {
		testInventoryRepo.UpdateInventoryFunc = func(ctx context.Context, inventory *models.Inventory) error {
			if inventory.ID == 0 {
				t.Error("Expected inventory ID to be populated")
			}
			CheckInventory(t, inventory)
			return nil
		}

		product_id := int64(1)
		quantity := int64(2)
		req := inventories.InventoryRequest{ProductID: &product_id, Quantity: &quantity}
		err := testInventorySvc.UpdateInventory(context.Background(), CheckUUID(t, UUID), req, models.AuditLogRequest{})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Delete Inventory", func(t *testing.T) {
		testInventoryRepo.DeleteInventoryFunc = func(ctx context.Context, uuid string) error {
			CheckUUID(t, uuid)
			return nil
		}

		err := testInventorySvc.DeleteInventory(context.Background(), CheckUUID(t, UUID), models.AuditLogRequest{})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Update Inventory Status", func(t *testing.T) {
		testInventoryRepo.UpdateInventoryStatusFunc = func(ctx context.Context, uuid string) error {
			CheckUUID(t, uuid)
			return nil
		}

		err := testInventorySvc.UpdateInventoryStatus(context.Background(), CheckUUID(t, UUID), models.AuditLogRequest{})
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

func CheckInventory(t testing.TB, inventory *models.Inventory) {
	t.Helper()

	if inventory.ProductID == 0 {
		t.Error("Expected inventory product_id to be populated")
	}
	if inventory.Quantity == 0 {
		t.Error("Expected inventory quantity to be populated")
	}
}

func CheckInventoryQuery(t testing.TB, q dto.BaseFilters) {
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
	if q.Sort != "quantity ASC" {
		t.Errorf("Expected sort to be quantity ASC, got '%s'", q.Sort)
	}
}
