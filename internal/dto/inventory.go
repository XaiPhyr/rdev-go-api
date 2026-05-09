package dto

type InventoryRequest struct {
	ProductID         *int64 `json:"product_id"`
	Quantity          *int64 `json:"quantity"`
	LowStockThreshold *int64 `json:"low_stock_threshold"`
}
