package orders

import (
	"context"
	"errors"

	"github.com/XaiPhyr/rdev-go-api/internal/audit_logs"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/email"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
	"github.com/redis/go-redis/v9"
)

type OrderRepository interface {
	GetOrderByUUID(ctx context.Context, uuid string) (*models.Order, error)
	GetOrders(ctx context.Context, q dto.BaseFilters) ([]models.Order, int, error)
	CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error)
	UpdateOrder(ctx context.Context, order *models.Order) (*models.Order, error)
	DeleteOrder(ctx context.Context, uuid string) error
	UpdateOrderStatus(ctx context.Context, uuid string) error
}

type OrderService interface {
	GetOrderByUUID(ctx context.Context, uuid string) (*models.Order, error)
	GetOrders(ctx context.Context, q dto.Query) ([]models.Order, int, error)
	CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error)
	UpdateOrder(ctx context.Context, uuid string, order *models.Order) (*models.Order, error)
	DeleteOrder(ctx context.Context, uuid string) error
	UpdateOrderStatus(ctx context.Context, uuid string) error
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

func (s *service) CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	if order.TotalAmount == 0 {
		return nil, errors.New("Order must have total amount")
	}

	return s.r.CreateOrder(ctx, order)
}

func (s *service) UpdateOrder(ctx context.Context, uuid string, req OrderRequest) (*models.Order, error) {
	if uuid == "" || len(uuid) != 36 {
		return nil, errors.New("Invalid UUID format")
	}

	if req.CustomerID == nil {
		return nil, errors.New("Must have customer id")
	}

	order, err := s.r.GetOrderByUUID(ctx, uuid)
	if err != nil {
		return nil, errors.New("Order not found")
	}

	if req.OrderStatus != nil {
		order.OrderStatus = *req.OrderStatus
	}

	return s.r.UpdateOrder(ctx, nil)
}

func (s *service) DeleteOrder(ctx context.Context, uuid string) error {
	if uuid == "" || len(uuid) != 36 {
		return errors.New("Invalid UUID format")
	}

	return s.r.DeleteOrder(ctx, uuid)
}

func (s *service) UpdateOrderStatus(ctx context.Context, uuid string) error {
	if uuid == "" || len(uuid) != 36 {
		return errors.New("Invalid UUID format")
	}

	return s.r.UpdateOrderStatus(ctx, uuid)
}
