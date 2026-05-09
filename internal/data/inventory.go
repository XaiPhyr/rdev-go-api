package data

import "github.com/uptrace/bun"

type Inventory struct {
	bun.BaseModel `bun:"table:inventories,alias:i"`
	BaseFields

	ProductID         int64 `bun:"product_id" json:"product_id"`
	Quantity          int64 `bun:"quantity" json:"quantity"`
	LowStockThreshold int64 `bun:"low_stock_threshold" json:"low_stock_threshold"`
}
