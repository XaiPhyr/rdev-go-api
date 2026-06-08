package carts_test

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/XaiPhyr/rdev-go-api/internal/carts"
	"github.com/XaiPhyr/rdev-go-api/internal/mocks"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
)

type CartTest struct {
	AddToCartFunc           func(ctx context.Context, cart *models.Cart) (*models.Cart, error)
	GetCartByCustomerIDFunc func(ctx context.Context, cart *carts.CartRequest) (*models.Cart, error)
}

func (m *CartTest) AddToCart(ctx context.Context, cart *models.Cart) (*models.Cart, error) {
	if m.AddToCartFunc != nil {
		return m.AddToCartFunc(ctx, cart)
	}

	return cart, nil
}

func (m *CartTest) GetCartByCustomerID(ctx context.Context, cart *carts.CartRequest) (*models.Cart, error) {
	if m.GetCartByCustomerIDFunc != nil {
		return m.GetCartByCustomerIDFunc(ctx, cart)
	}

	return &models.Cart{}, nil
}

func NewTestCart(testCartRepo *CartTest) carts.CartService {
	emailSvc := mocks.NewTestEmailService()
	_, auditLogSvc := mocks.NewTestAuditService()

	return carts.NewCartService(testCartRepo, emailSvc, nil, auditLogSvc)
}

func TestAddToCart(t *testing.T) {
	testCartRepo := &CartTest{}
	testCartSvc := NewTestCart(testCartRepo)

	t.Run("Add To Cart", func(t *testing.T) {
		testCartRepo.AddToCartFunc = func(ctx context.Context, cart *models.Cart) (*models.Cart, error) {
			return &models.Cart{}, nil
		}

		cartItem := []carts.CartItemRequest{
			{ProductName: new("Product 1")},
		}

		_, err := testCartSvc.AddToCart(context.Background(), &carts.CartRequest{CartItem: cartItem}, models.AuditLogRequest{})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestUserNoActiveCart(t *testing.T) {
	testCartRepo := &CartTest{}
	testCartSvc := NewTestCart(testCartRepo)

	t.Run("User no active cart", func(t *testing.T) {
		testCartRepo.AddToCartFunc = func(ctx context.Context, cart *models.Cart) (*models.Cart, error) {
			// check database if user has active cart
			// if not add cart then add items
			userCart := &models.Cart{}

			if userCart.CustomerID != cart.CustomerID {
				userCart.CustomerID = cart.CustomerID
				userCart.CartStatus = cart.CartStatus
				userCart.CartItems = cart.CartItems
			}

			return userCart, nil
		}

		cartItem := []carts.CartItemRequest{
			{ProductID: new(int64(1)), ProductName: new("Product 1"), Quantity: new(int64(0))},
		}

		cart := &carts.CartRequest{CustomerID: new(int64(1)), CartItem: cartItem, CartStatus: new("ACTIVE")}

		_, err := testCartSvc.AddToCart(context.Background(), cart, models.AuditLogRequest{})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestUpdateCartItemQuantity(t *testing.T) {
	testCartRepo := &CartTest{}
	testCartSvc := NewTestCart(testCartRepo)

	t.Run("Update cart item quantity", func(t *testing.T) {
		testCartRepo.AddToCartFunc = func(ctx context.Context, cart *models.Cart) (*models.Cart, error) {
			activeCarts := models.Cart{
				CustomerID: 2,
				CartStatus: "ACTIVE",
				CartItems: []*models.CartItem{
					{CartID: 1, ProductID: 2, Quantity: 1},
					{CartID: 1, ProductID: 3, Quantity: 1},
					{CartID: 1, ProductID: 1, Quantity: 1},
				},
			}

			for _, ci := range cart.CartItems {
				idx := slices.IndexFunc(activeCarts.CartItems, func(cartItem *models.CartItem) bool {
					return cartItem.ProductID == ci.ProductID
				})

				if ci.Quantity == 0 {
					activeCarts.CartItems = slices.Delete(activeCarts.CartItems, idx, idx+1)
				}

				if ci.Quantity > 0 {
					activeCarts.CartItems[idx].Quantity = ci.Quantity + activeCarts.CartItems[idx].Quantity
				}
			}

			return &activeCarts, nil
		}

		cartItem := []carts.CartItemRequest{
			{CartID: new(int64(1)), ProductID: new(int64(1)), ProductName: new("Product 1"), Quantity: new(int64(0))},
		}

		cart := &carts.CartRequest{CustomerID: new(int64(1)), CartItem: cartItem, CartStatus: new("SUSPENDED")}

		updatedCart, err := testCartSvc.AddToCart(context.Background(), cart, models.AuditLogRequest{})

		for _, v := range updatedCart.CartItems {
			fmt.Println()
			fmt.Println("UPDATED CART: ", v.ProductID, v.Quantity)
			fmt.Println()
		}

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestGetCartByCustomerID(t *testing.T) {
	testCartRepo := &CartTest{}
	testCartSvc := NewTestCart(testCartRepo)

	t.Run("Get cart by customer id", func(t *testing.T) {
		testCartRepo.GetCartByCustomerIDFunc = func(ctx context.Context, cart *carts.CartRequest) (*models.Cart, error) {
			return &models.Cart{}, nil
		}

		_, err := testCartSvc.GetCartByCustomerID(context.Background(), &carts.CartRequest{CustomerID: new(int64(0))})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}
