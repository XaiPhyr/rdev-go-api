package carts_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/XaiPhyr/rdev-go-api/internal/carts"
	"github.com/XaiPhyr/rdev-go-api/internal/mocks"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
)

type CartTest struct {
	AddToCartFunc func(ctx context.Context, cart *models.Cart) (*models.Cart, error)
}

func (m *CartTest) AddToCart(ctx context.Context, cart *models.Cart) (*models.Cart, error) {
	fmt.Println()
	fmt.Println("TEST", cart.CartItems)
	fmt.Println()

	if m.AddToCartFunc != nil {
		return m.AddToCartFunc(ctx, cart)
	}

	return cart, nil
}

func TestCart(t *testing.T) {
	testCartRepo := &CartTest{}
	emailSvc := mocks.NewTestEmailService()
	_, auditLogSvc := mocks.NewTestAuditService()

	testCartSvc := carts.NewCartService(testCartRepo, emailSvc, nil, auditLogSvc)

	t.Run("Add To Cart", func(t *testing.T) {
		testCartRepo.AddToCartFunc = func(ctx context.Context, cart *models.Cart) (*models.Cart, error) {
			return &models.Cart{}, nil
		}

		cartItem := []carts.CartItemRequest{
			{ProductName: new("Product 1")},
		}

		_, err := testCartSvc.AddToCart(context.Background(), carts.CartRequest{CartItem: cartItem}, models.AuditLogRequest{})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}
