package orders

import (
	"context"
	"errors"

	"github.com/XaiPhyr/rdev-go-api/internal/audit_logs"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/email"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
	"github.com/redis/go-redis/v9"
)

type OrderRepository interface {
	GetOrderByUUID(ctx context.Context, uuid string) (*models.Order, error)
}

type OrderService interface {
	GetOrderByUUID(ctx context.Context, uuid string) (*models.Order, error)
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
