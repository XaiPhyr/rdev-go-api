package stock_movements_test

import (
	"context"
	"sync"
	"testing"

	"github.com/XaiPhyr/rdev-go-api/internal/mocks"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
	"github.com/XaiPhyr/rdev-go-api/internal/stock_movements"
)

const UUID = "12345678-1234-5678-1234-567890123456"

type StockMovementTest struct {
	GetStockMovementByUUIDFunc    func(ctx context.Context, uuid string) (*models.StockMovement, error)
	GetStockMovementsFunc         func(ctx context.Context, filters dto.BaseFilters) ([]models.StockMovement, int, error)
	CreateStockMovementFunc       func(ctx context.Context, sm *models.StockMovement) error
	UpdateStockMovementFunc       func(ctx context.Context, sm *models.StockMovement) error
	DeleteStockMovementFunc       func(ctx context.Context, uuid string) error
	UpdateStockMovementStatusFunc func(ctx context.Context, uuid string) error
	ProcessBulkUploadFunc         func(ctx context.Context, row [][]string) error
}

func (m *StockMovementTest) GetStockMovementByUUID(ctx context.Context, uuid string) (*models.StockMovement, error) {
	if m.GetStockMovementByUUIDFunc != nil {
		return m.GetStockMovementByUUIDFunc(ctx, uuid)
	}

	return nil, nil
}
func (m *StockMovementTest) GetStockMovements(ctx context.Context, q dto.BaseFilters) ([]models.StockMovement, int, error) {
	if m.GetStockMovementsFunc != nil {
		return m.GetStockMovementsFunc(ctx, q)
	}

	return nil, 0, nil
}
func (m *StockMovementTest) CreateStockMovement(ctx context.Context, sm *models.StockMovement) error {
	if m.CreateStockMovementFunc != nil {
		return m.CreateStockMovementFunc(ctx, sm)
	}

	return nil
}
func (m *StockMovementTest) UpdateStockMovement(ctx context.Context, sm *models.StockMovement) error {
	if m.UpdateStockMovementFunc != nil {
		sm.ID = 1
		return m.UpdateStockMovementFunc(ctx, sm)
	}
	return nil
}
func (m *StockMovementTest) DeleteStockMovement(ctx context.Context, uuid string) error {
	if m.DeleteStockMovementFunc != nil {
		return m.DeleteStockMovementFunc(ctx, uuid)
	}

	return nil
}
func (m *StockMovementTest) UpdateStockMovementStatus(ctx context.Context, uuid string) error {
	if m.UpdateStockMovementStatusFunc != nil {
		return m.UpdateStockMovementStatusFunc(ctx, uuid)
	}

	return nil
}
func (m *StockMovementTest) ProcessBulkUpload(ctx context.Context, row [][]string) error {
	if m.ProcessBulkUploadFunc != nil {
		return m.ProcessBulkUploadFunc(ctx, row)
	}

	return nil
}

func TestStockMovement(t *testing.T) {
	testStockMovementRepo := &StockMovementTest{}
	emailSvc := mocks.NewTestEmailService()
	_, auditLogSvc := mocks.NewTestAuditService()

	testStockMovementSvc := stock_movements.NewStockMovementService(testStockMovementRepo, emailSvc, nil, auditLogSvc)

	t.Run("Get StockMovements", func(t *testing.T) {
		testStockMovementRepo.GetStockMovementsFunc = func(ctx context.Context, q dto.BaseFilters) ([]models.StockMovement, int, error) {
			CheckStockMovementQuery(t, q)
			return []models.StockMovement{{ProductID: 1}}, 1, nil
		}

		query := dto.Query{Search: "test", Limit: 10, Offset: 2, Sort: "change_amount ASC"}
		_, _, err := testStockMovementSvc.GetStockMovements(context.Background(), query)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Get StockMovement By UUID", func(t *testing.T) {
		testStockMovementRepo.GetStockMovementByUUIDFunc = func(ctx context.Context, uuid string) (*models.StockMovement, error) {
			CheckUUID(t, uuid)
			return &models.StockMovement{ProductID: 1}, nil
		}

		_, err := testStockMovementSvc.GetStockMovementByUUID(context.Background(), UUID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Create StockMovement", func(t *testing.T) {
		testStockMovementRepo.CreateStockMovementFunc = func(ctx context.Context, sm *models.StockMovement) error {
			CheckStockMovement(t, sm)
			return nil
		}

		numRequest := 50
		var wg sync.WaitGroup

		for i := range numRequest {
			wg.Go(func() {
				product_id := int64(i + 1)
				quantity := int64(i + 1*2)
				req := stock_movements.StockMovementRequest{ProductID: &product_id, ChangeAmount: &quantity}
				err := testStockMovementSvc.CreateStockMovement(context.Background(), req, models.AuditLogRequest{})

				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			})
		}

		wg.Wait()
	})

	t.Run("Update StockMovement", func(t *testing.T) {
		testStockMovementRepo.UpdateStockMovementFunc = func(ctx context.Context, sm *models.StockMovement) error {
			if sm.ID == 0 {
				t.Error("Expected sm ID to be populated")
			}
			CheckStockMovement(t, sm)
			return nil
		}

		product_id := int64(1)
		quantity := int64(2)
		req := stock_movements.StockMovementRequest{ProductID: &product_id, ChangeAmount: &quantity}
		err := testStockMovementSvc.UpdateStockMovement(context.Background(), CheckUUID(t, UUID), req, models.AuditLogRequest{})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Delete StockMovement", func(t *testing.T) {
		testStockMovementRepo.DeleteStockMovementFunc = func(ctx context.Context, uuid string) error {
			CheckUUID(t, uuid)
			return nil
		}

		err := testStockMovementSvc.DeleteStockMovement(context.Background(), CheckUUID(t, UUID), models.AuditLogRequest{})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Update StockMovement Status", func(t *testing.T) {
		testStockMovementRepo.UpdateStockMovementStatusFunc = func(ctx context.Context, uuid string) error {
			CheckUUID(t, uuid)
			return nil
		}

		err := testStockMovementSvc.UpdateStockMovementStatus(context.Background(), CheckUUID(t, UUID), models.AuditLogRequest{})
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

func CheckStockMovement(t testing.TB, sm *models.StockMovement) {
	t.Helper()

	if sm.ProductID == 0 {
		t.Error("Expected sm product_id to be populated")
	}
	if sm.ChangeAmount == 0 {
		t.Error("Expected sm change_amount to be populated")
	}
}

func CheckStockMovementQuery(t testing.TB, q dto.BaseFilters) {
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
	if q.Sort != "change_amount ASC" {
		t.Errorf("Expected sort to be change_amount ASC, got '%s'", q.Sort)
	}
}
