package products

import (
	"github.com/XaiPhyr/rdev-go-api/internal/categories"
	"github.com/XaiPhyr/rdev-go-api/internal/inventories"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/fields"
	"github.com/XaiPhyr/rdev-go-api/internal/stock_movements"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:products,alias:p"`
	fields.BaseFields

	CategoryID    int64                            `bun:"category_id" json:"category_id"`
	Category      *categories.Category             `bun:"rel:belongs-to,join:category_id=id" json:"category,omitempty"`
	Name          string                           `bun:"name,notnull" json:"name"`
	Slug          string                           `bun:"slug" json:"slug"`
	Description   string                           `bun:"description" json:"description"`
	SKU           string                           `bun:"sku" json:"sku"`
	Barcode       string                           `bun:"barcode" json:"barcode"`
	Price         int64                            `bun:"price,notnull" json:"price"`
	CostPrice     int64                            `bun:"cost_price,notnull" json:"cost_price"`
	DisplayPrice  float64                          `bun:"column:display_price,scanonly" json:"display_price"`
	Inventory     *inventories.Inventory           `bun:"rel:has-one,join:id=product_id" json:"inventory,omitempty"`
	StockMovement []*stock_movements.StockMovement `bun:"rel:has-many,join:id=product_id" json:"stock_movement,omitempty"`
}

type ProductRequest struct {
	CategoryID  *int64  `json:"category_id"`
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
	SKU         *string `json:"sku"`
	Barcode     *string `json:"barcode"`
	Price       *int64  `json:"price"`
	CostPrice   *int64  `json:"cost_price"`
	Quantity    *int64  `json:"quantity"`
}

// Public API response (for customer)
type ProductPublicResponse struct {
	Category     *categories.CategoryResponse `json:"category"`
	Name         string                       `json:"name"`
	Slug         string                       `json:"slug"`
	Description  string                       `json:"description"`
	Barcode      string                       `json:"barcode"`
	DisplayPrice float64                      `json:"display_price"`
}

// Backoffice API response (for staff)
type ProductBackofficeResponse struct {
	Category     *categories.CategoryResponse `json:"category"`
	Name         string                       `json:"name"`
	Slug         string                       `json:"slug"`
	Description  string                       `json:"description"`
	SKU          string                       `json:"sku"`
	Barcode      string                       `json:"barcode"`
	DisplayPrice float64                      `json:"display_price"`
	Quantity     int64                        `json:"quantity"`
}
