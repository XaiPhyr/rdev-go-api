package models

import (
	"github.com/XaiPhyr/rdev-go-api/internal/shared/fields"
	"github.com/uptrace/bun"
)

type Order struct {
	bun.BaseModel `bun:"table:orders,alias:o"`
	fields.BaseFields

	CustomerID    int64       `bun:"customer_id" json:"customer_id" validate:"required"`
	OrderNumber   string      `bun:"order_number" json:"order_number"`
	OrderStatus   string      `bun:"order_status" json:"order_status"`
	PaymentStatus string      `bun:"payment_status" json:"payment_status"`
	TotalAmount   int64       `bun:"total_amount" json:"total_amount"`
	OrderItem     []OrderItem `bun:"rel:has-many,join:id=order_id" json:"order_item,omitempty"`
}

type OrderItem struct {
	bun.BaseModel `bun:"table:order_items,alias:oi"`
	fields.BaseFields

	OrderID          int64 `bun:"order_id" json:"order_id"`
	ProductID        int64 `bun:"product_id" json:"product_id"`
	TransactionPrice int64 `bun:"transaction_price" json:"transaction_price"`
	Quantity         int64 `bun:"quantity" json:"quantity"`
}
