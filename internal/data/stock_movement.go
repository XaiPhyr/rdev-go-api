package data

import "github.com/uptrace/bun"

type StockMovement struct {
	bun.BaseModel `bun:"table:stock_movements,alias:sm"`
	BaseFields

	ProductID    int64  `bun:"product_id" json:"product_id"`
	ChangeAmount int64  `bun:"change_amount" json:"change_amount"`
	Reason       string `bun:"reason,default:'INITIAL_STOCK'" json:"reason"`
	ReferenceID  string `bun:"reference_id" json:"reference_id"`
}
