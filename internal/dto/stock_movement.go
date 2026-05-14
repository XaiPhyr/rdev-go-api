package dto

type StockMovementRequest struct {
	ProductID    *int64  `json:"product_id"`
	ChangeAmount *int64  `json:"change_amount"`
	Reason       *string `json:"reason"`
	ReferenceID  *string `json:"reference_id"`
}

type BulkUploadRequest struct {
	File string `json:"file"`
}
