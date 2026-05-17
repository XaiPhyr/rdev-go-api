package stock_movements

import (
	"github.com/XaiPhyr/rdev-go-api/internal/shared/fields"
	"github.com/uptrace/bun"
)

type StockMovement struct {
	bun.BaseModel `bun:"table:stock_movements,alias:sm"`
	fields.BaseFields

	ProductID    int64  `bun:"product_id" json:"product_id"`
	ChangeAmount int64  `bun:"change_amount" json:"change_amount"`
	Reason       string `bun:"reason,default:'INITIAL_STOCK'" json:"reason"`
	ReferenceID  string `bun:"reference_id" json:"reference_id"`
}

type StockMovementRequest struct {
	ProductID    *int64  `json:"product_id"`
	ChangeAmount *int64  `json:"change_amount"`
	Reason       *string `json:"reason"`
	ReferenceID  *string `json:"reference_id"`
}

type BulkUploadRequest struct {
	File string `json:"file"`
}
