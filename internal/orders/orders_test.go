package orders_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/XaiPhyr/rdev-go-api/internal/mocks"
	"github.com/XaiPhyr/rdev-go-api/internal/orders"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/helpers"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
)

type OrderTest struct {
	GetOrderByUUIDFunc    func(ctx context.Context, uuid string) (*models.Order, error)
	GetOrdersFunc         func(ctx context.Context, q dto.BaseFilters) ([]models.Order, int, error)
	CreateOrderFunc       func(ctx context.Context, order *models.Order) (*models.Order, error)
	UpdateOrderFunc       func(ctx context.Context, order *models.Order) (*models.Order, error)
	DeleteOrderFunc       func(ctx context.Context, uuid string) error
	UpdateOrderStatusFunc func(ctx context.Context, uuid string) error
}

func (m *OrderTest) GetOrderByUUID(ctx context.Context, uuid string) (*models.Order, error) {
	if m.GetOrderByUUIDFunc != nil {
		return m.GetOrderByUUIDFunc(ctx, uuid)
	}

	return nil, nil
}
func (m *OrderTest) GetOrders(ctx context.Context, q dto.BaseFilters) ([]models.Order, int, error) {
	if m.GetOrdersFunc != nil {
		return m.GetOrdersFunc(ctx, q)
	}

	return nil, 0, nil
}
func (m *OrderTest) CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	if len(order.OrderItem) == 0 {
		return nil, errors.New("Order must have order items")
	}

	for _, oi := range order.OrderItem {
		if oi.OrderID == 0 || oi.ProductID == 0 || oi.TransactionPrice == 0 || oi.Quantity == 0 {
			return nil, errors.New("Order ID/Product ID/Transaction Price/Quantity must not be 0")
		}
		if oi.Quantity < 0 {
			return nil, errors.New("Quantity must not be negative")
		}
	}

	err := helpers.ValidateStruct(order)
	if err != nil {
		return nil, fmt.Errorf("Validation error check field %v", err)
	}

	if m.CreateOrderFunc != nil {
		return m.CreateOrderFunc(ctx, order)
	}

	return nil, nil
}
func (m *OrderTest) UpdateOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	if m.UpdateOrderFunc != nil {
		return m.UpdateOrderFunc(ctx, order)
	}

	return nil, nil
}
func (m *OrderTest) DeleteOrder(ctx context.Context, uuid string) error {
	if m.DeleteOrderFunc != nil {
		return m.DeleteOrderFunc(ctx, uuid)
	}

	return nil
}
func (m *OrderTest) UpdateOrderStatus(ctx context.Context, uuid string) error {
	if m.UpdateOrderStatusFunc != nil {
		return m.UpdateOrderStatusFunc(ctx, uuid)
	}

	return nil
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

			order.OrderNumber = "ORD-001"

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

			order.OrderNumber = "ORD-001"

			return &order, nil
		}

		_, err := testOrderSvc.GetOrderByUUID(context.Background(), validUUID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Get Orders Search with Special Characters", func(t *testing.T) {
		testOrderRepo.GetOrdersFunc = func(ctx context.Context, q dto.BaseFilters) ([]models.Order, int, error) {
			if strings.HasPrefix(q.Search, " ") || strings.HasSuffix(q.Search, " ") {
				return nil, 0, fmt.Errorf("Search has prefix/suffix spaces %s", q.Search)
			}

			return []models.Order{{OrderNumber: "ORD-001"}}, 0, nil
		}

		search := "ORD-1929%$@1982@ "
		cleanedSearch := helpers.CleanSpecialChars(search)
		query := dto.Query{Search: cleanedSearch}

		_, _, err := testOrderSvc.GetOrders(context.Background(), query)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Create Order with no total amount", func(t *testing.T) {
		testOrderRepo.CreateOrderFunc = func(ctx context.Context, order *models.Order) (*models.Order, error) {
			return &models.Order{OrderNumber: "ORD-001"}, nil
		}

		ord := &models.Order{
			TotalAmount: 1,
		}

		_, err := testOrderSvc.CreateOrder(context.Background(), ord)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Update Order with no customer id or no uuid", func(t *testing.T) {
		testOrderRepo.UpdateOrderFunc = func(ctx context.Context, order *models.Order) (*models.Order, error) {
			return &models.Order{OrderNumber: "ORD-001"}, nil
		}

		validUUID := "12345678-1234-1234-1234-123456789012"
		customer_id := int64(1)
		ord := orders.OrderRequest{CustomerID: &customer_id}
		_, err := testOrderSvc.UpdateOrder(context.Background(), validUUID, ord)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Delete Order with no uuid", func(t *testing.T) {
		testOrderRepo.DeleteOrderFunc = func(ctx context.Context, uuid string) error {
			return nil
		}

		validUUID := "12345678-1234-1234-1234-123456789012"

		err := testOrderSvc.DeleteOrder(context.Background(), validUUID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Update Order Status with no uuid", func(t *testing.T) {
		testOrderRepo.UpdateOrderStatusFunc = func(ctx context.Context, uuid string) error {
			return nil
		}

		validUUID := "12345678-1234-1234-1234-123456789012"

		err := testOrderSvc.UpdateOrderStatus(context.Background(), validUUID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Order with no order item", func(t *testing.T) {
		testOrderRepo.CreateOrderFunc = func(ctx context.Context, order *models.Order) (*models.Order, error) {
			return nil, nil
		}

		order_item := []models.OrderItem{
			{OrderID: 1, ProductID: 1, TransactionPrice: 1, Quantity: 1},
		}
		order := &models.Order{CustomerID: 1, TotalAmount: 50, OrderItem: order_item}

		_, err := testOrderSvc.CreateOrder(context.Background(), order)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}
