package orders

import (
	"context"
	"errors"
	"fmt"

	"github.com/XaiPhyr/rdev-go-api/internal/audit_logs"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/email"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/helpers"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
	"github.com/redis/go-redis/v9"
)

type OrderRepository interface {
	GetOrderByUUID(ctx context.Context, uuid string) (*models.Order, error)
	GetOrders(ctx context.Context, q dto.BaseFilters) ([]models.Order, int, error)
	CreateOrder(ctx context.Context, order *models.Order) error
	UpdateOrder(ctx context.Context, order *models.Order) error
	DeleteOrder(ctx context.Context, uuid string) error
	UpdateOrderStatus(ctx context.Context, uuid string) error
}

type OrderService interface {
	GetOrderByUUID(ctx context.Context, uuid string) (*models.Order, error)
	GetOrders(ctx context.Context, q dto.Query) ([]models.Order, int, error)
	CreateOrder(ctx context.Context, req OrderRequest, audit models.AuditLogRequest) error
	UpdateOrder(ctx context.Context, uuid string, req OrderRequest, audit models.AuditLogRequest) error
	DeleteOrder(ctx context.Context, uuid string, audit models.AuditLogRequest) error
	UpdateOrderStatus(ctx context.Context, uuid string, audit models.AuditLogRequest) error
}

type service struct {
	r        OrderRepository
	email    email.EmailService
	redis    *redis.Client
	auditLog audit_logs.AuditLogService
}

func NewOrderService(r OrderRepository, email email.EmailService, redis *redis.Client, auditLog audit_logs.AuditLogService) *service {
	return &service{r: r, email: email, redis: redis, auditLog: auditLog}
}

func (s *service) GetOrderByUUID(ctx context.Context, uuid string) (*models.Order, error) {
	if uuid == "" || len(uuid) != 36 {
		return nil, errors.New("Invalid UUID format")
	}

	order, err := s.r.GetOrderByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if order.OrderNumber == "" {
		return nil, errors.New("Order must have a valid order number")
	}
	if len(order.OrderNumber) > 25 {
		return nil, errors.New("Order number exceeds maximum length of 25 characters")
	}

	return order, nil
}

func (s *service) GetOrders(ctx context.Context, q dto.Query) ([]models.Order, int, error) {
	filters := q.SanitizeQuery([]string{"order_number"})

	return s.r.GetOrders(ctx, filters)
}

func (s *service) CreateOrder(ctx context.Context, req OrderRequest, audit models.AuditLogRequest) error {
	order := &models.Order{}

	if req.CustomerID != nil {
		order.CustomerID = *req.CustomerID
	}

	if req.OrderStatus != nil {
		order.OrderStatus = *req.OrderStatus
	}

	if req.TotalAmount != nil {
		order.TotalAmount = *req.TotalAmount
	}

	if len(req.OrderItem) > 0 {
		for _, roi := range req.OrderItem {
			oi := models.OrderItem{}

			if roi.ProductID != nil {
				oi.ProductID = *roi.ProductID
			}

			if roi.TransactionPrice != nil {
				oi.TransactionPrice = *roi.TransactionPrice
			}

			if roi.Quantity != nil {
				oi.Quantity = *roi.Quantity
			}

			order.OrderItem = append(order.OrderItem, oi)
		}
	}

	if order.TotalAmount == 0 {
		return errors.New("Order must have total amount")
	}

	if len(order.OrderItem) == 0 {
		return errors.New("Order must have order items")
	}

	for _, oi := range order.OrderItem {
		if oi.ProductID == 0 || oi.TransactionPrice == 0 || oi.Quantity == 0 {
			return errors.New("Order ID/Product ID/Transaction Price/Quantity must not be 0")
		}
		if oi.Quantity < 0 {
			return errors.New("Quantity must not be negative")
		}
	}

	err := helpers.ValidateStruct(order)
	if err != nil {
		return fmt.Errorf("Validation error check field %v", err)
	}

	err = s.r.CreateOrder(ctx, order)

	return err
}

func (s *service) UpdateOrder(ctx context.Context, uuid string, req OrderRequest, audit models.AuditLogRequest) error {
	if uuid == "" || len(uuid) != 36 {
		return errors.New("Invalid UUID format")
	}

	if req.CustomerID == nil {
		return errors.New("Must have customer id")
	}

	order, err := s.r.GetOrderByUUID(ctx, uuid)
	if err != nil {
		return errors.New("Order not found")
	}

	if req.OrderStatus != nil {
		order.OrderStatus = *req.OrderStatus
	}

	err = s.r.UpdateOrder(ctx, nil)

	return err
}

func (s *service) DeleteOrder(ctx context.Context, uuid string, audit models.AuditLogRequest) error {
	if uuid == "" || len(uuid) != 36 {
		return errors.New("Invalid UUID format")
	}

	return s.r.DeleteOrder(ctx, uuid)
}

func (s *service) UpdateOrderStatus(ctx context.Context, uuid string, audit models.AuditLogRequest) error {
	if uuid == "" || len(uuid) != 36 {
		return errors.New("Invalid UUID format")
	}

	return s.r.UpdateOrderStatus(ctx, uuid)
}
