package models

import "github.com/uptrace/bun"

type Product struct {
	bun.BaseModel `bun:"table:products,alias:p"`

	ID           int64   `bun:"id,pk,autoincrement" json:"id"`
	CategoryID   int64   `bun:"category_id" json:"category_id"`
	Name         string  `bun:"name,notnull" json:"name"`
	Slug         string  `bun:"slug" json:"slug"`
	Description  string  `bun:"description" json:"description"`
	SKU          string  `bun:"sku" json:"sku"`
	Barcode      string  `bun:"barcode" json:"barcode"`
	Price        int64   `bun:"price,notnull" json:"price"`
	CostPrice    int64   `bun:"cost_price,notnull" json:"cost_price"`
	DisplayPrice float64 `bun:"column:display_price,scanonly" json:"display_price"`
}
