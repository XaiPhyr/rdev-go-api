package carts

import (
	"context"
	"errors"

	"github.com/XaiPhyr/rdev-go-api/internal/audit_logs"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/email"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
	"github.com/redis/go-redis/v9"
)

type CartRepository interface {
	AddToCart(ctx context.Context, cart *models.Cart) (*models.Cart, error)
}

type CartService interface {
	AddToCart(ctx context.Context, cart CartRequest, audit models.AuditLogRequest) (*models.Cart, error)
}

type service struct {
	r        CartRepository
	email    email.EmailService
	redis    *redis.Client
	auditLog audit_logs.AuditLogService
}

func NewCartService(r CartRepository, email email.EmailService, redis *redis.Client, auditLog audit_logs.AuditLogService) *service {
	return &service{r: r, email: email, redis: redis, auditLog: auditLog}
}

func (s *service) AddToCart(ctx context.Context, req CartRequest, audit models.AuditLogRequest) (*models.Cart, error) {
	cart := &models.Cart{}

	if len(req.CartItem) == 0 {
		return nil, errors.New("Cannot proceed with empty cart")
	}

	for _, rci := range req.CartItem {
		cart_item := &models.CartItem{
			ProductName: *rci.ProductName,
		}

		cart.CartItems = append(cart.CartItems, cart_item)
	}

	cart, err := s.r.AddToCart(ctx, cart)
	s.auditLog.ParseAndCreateAuditLog(audit, cart.UUID, "CART", nil, *cart, err)

	return cart, err
}
