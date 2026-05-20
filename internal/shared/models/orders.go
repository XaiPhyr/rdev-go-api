package models

import (
	"github.com/XaiPhyr/rdev-go-api/internal/shared/fields"
	"github.com/uptrace/bun"
)

type Order struct {
	bun.BaseModel `bun:"table:orders,alias:o"`
	fields.BaseFields

	ParentID    int64  `bun:"-" json:"parent_id,omitempty"`
	CustomerID  int64  `bun:"customer_id" json:"customer_id"`
	OrderNumber string `bun:"order_number" json:"order_number"`
	OrderStatus string `bun:"order_status" json:"order_status"`
	TotalAmount int64  `bun:"total_amount" json:"total_amount"`
}
