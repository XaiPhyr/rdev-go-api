package models

import (
	"github.com/XaiPhyr/rdev-go-api/internal/shared/fields"
	"github.com/uptrace/bun"
)

type Cart struct {
	bun.BaseModel `bun:"table:carts,alias:c"`
	fields.BaseFields

	CustomerID int64       `bun:"customer_id" json:"customer_id"`
	CartStatus string      `bun:"cart_status" json:"cart_status"`
	CartItems  []*CartItem `bun:"rel:has-many,join:id=cart_id" json:"cart_item,omitempty"`
}

type CartItem struct {
	bun.BaseModel `bun:"table:cart_items,alias:ci"`
	fields.BaseFields

	CartID           int64  `bun:"cart_id" json:"cart_id"`
	ProductID        int64  `bun:"product_id" json:"product_id"`
	ProductName      string `bun:"product_name" json:"product_name"`
	TransactionPrice int64  `bun:"transaction_price" json:"transaction_price"`
	Quantity         int64  `bun:"quantity" json:"quantity"`
}
