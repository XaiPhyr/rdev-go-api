package stock_movements

type StockMovementRequest struct {
	ProductID    *int64  `json:"product_id"`
	ChangeAmount *int64  `json:"change_amount"`
	Reason       *string `json:"reason"`
	ReferenceID  *string `json:"reference_id"`
}

type BulkUploadRequest struct {
	File string `json:"file"`
}

type BulkUploadErrResponse struct {
	Row         int    `json:"row"`
	Name        string `json:"name"`
	SKU         string `json:"sku"`
	ExistingSKU string `json:"existing_sku"`
}
