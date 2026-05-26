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
	AddToCart(ctx context.Context, cart *CartRequest, audit models.AuditLogRequest) (*models.Cart, error)
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

func (s *service) AddToCart(ctx context.Context, req *CartRequest, audit models.AuditLogRequest) (*models.Cart, error) {
	cart := &models.Cart{}

	if len(req.CartItem) == 0 {
		return nil, errors.New("Cannot proceed with empty cart")
	}

	if req.CustomerID != nil {
		cart.CustomerID = *req.CustomerID
	}

	if req.CartStatus != nil {
		cart.CartStatus = *req.CartStatus
	}

	for _, rci := range req.CartItem {
		cart_item := &models.CartItem{}

		if rci.CartID != nil {
			cart_item.CartID = *rci.CartID
		}

		if rci.ProductID != nil {
			cart_item.ProductID = *rci.ProductID
		}

		if rci.ProductName != nil {
			cart_item.ProductName = *rci.ProductName
		}

		if rci.Quantity != nil {
			cart_item.Quantity = *rci.Quantity
		}

		cart.CartItems = append(cart.CartItems, cart_item)
	}

	cart, err := s.r.AddToCart(ctx, cart)
	s.auditLog.ParseAndCreateAuditLog(audit, cart.UUID, "CART", nil, *cart, err)

	return cart, err
}
