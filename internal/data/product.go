package data

import "github.com/uptrace/bun"

type Product struct {
	bun.BaseModel `bun:"table:products,alias:p"`
	BaseFields

	CategoryID   int64     `bun:"category_id" json:"category_id"`
	Category     *Category `bun:"rel:belongs-to,join:category_id=id" json:"category,omitempty"`
	Name         string    `bun:"name,notnull" json:"name"`
	Slug         string    `bun:"slug" json:"slug"`
	Description  string    `bun:"description" json:"description"`
	SKU          string    `bun:"sku" json:"sku"`
	Barcode      string    `bun:"barcode" json:"barcode"`
	Price        int64     `bun:"price,notnull" json:"price"`
	CostPrice    int64     `bun:"cost_price,notnull" json:"cost_price"`
	DisplayPrice float64   `bun:"column:display_price,scanonly" json:"display_price"`
}
