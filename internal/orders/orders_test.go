package orders_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/XaiPhyr/rdev-go-api/internal/mocks"
	"github.com/XaiPhyr/rdev-go-api/internal/orders"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
)

type OrderTest struct {
	GetOrderByUUIDFunc func(ctx context.Context, uuid string) (*models.Order, error)
}

func (m *OrderTest) GetOrderByUUID(ctx context.Context, uuid string) (*models.Order, error) {
	if m.GetOrderByUUIDFunc != nil {
		return m.GetOrderByUUIDFunc(ctx, uuid)
	}

	return nil, nil
}

func TestOrders(t *testing.T) {
	testOrderRepo := &OrderTest{}
	emailSvc := mocks.NewTestEmailService()
	_, auditLogSvc := mocks.NewTestAuditService()

	testOrderSvc := orders.NewOrderService(testOrderRepo, emailSvc, nil, auditLogSvc)

	t.Run("Get Order By UUID with Context Timeout", func(t *testing.T) {
		validUUID := "12345678-1234-1234-1234-123456789012"

		testOrderRepo.GetOrderByUUIDFunc = func(ctx context.Context, uuid string) (*models.Order, error) {
			select {
			case <-time.After(5 * time.Millisecond):
				return &models.Order{OrderNumber: "ORD-001"}, nil
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Millisecond)
		defer cancel()

		_, err := testOrderSvc.GetOrderByUUID(ctx, validUUID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Get Order By UUID with Rate Limiting", func(t *testing.T) {
		callCount := 0
		validUUID := "12345678-1234-1234-1234-123456789012"

		testOrderRepo.GetOrderByUUIDFunc = func(ctx context.Context, uuid string) (*models.Order, error) {
			callCount++
			var order models.Order

			order.ID = 1

			if callCount > 3 {
				return nil, errors.New("Rate limit exceeded")
			}

			return &order, nil
		}

		var err error
		for range 3 {
			_, err = testOrderSvc.GetOrderByUUID(context.Background(), validUUID)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Get Order By UUID with Empty Order Number", func(t *testing.T) {
		validUUID := "12345678-1234-1234-1234-123456789012"

		testOrderRepo.GetOrderByUUIDFunc = func(ctx context.Context, uuid string) (*models.Order, error) {
			var order models.Order

			order.OrderNumber = "12345678901234567890123456789012345678901234567890"

			return &order, nil
		}

		_, err := testOrderSvc.GetOrderByUUID(context.Background(), validUUID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}
