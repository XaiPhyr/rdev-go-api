package dto

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
	Category     *CategoryPublicResponse `json:"category"`
	Name         string                  `json:"name"`
	Slug         string                  `json:"slug"`
	Description  string                  `json:"description"`
	Barcode      string                  `json:"barcode"`
	DisplayPrice float64                 `json:"display_price"`
}

// Backoffice API response (for staff)
type ProductBackofficeResponse struct {
	Category     *CategoryPublicResponse `json:"category"`
	Name         string                  `json:"name"`
	Slug         string                  `json:"slug"`
	Description  string                  `json:"description"`
	SKU          string                  `json:"sku"`
	Barcode      string                  `json:"barcode"`
	DisplayPrice float64                 `json:"display_price"`
	Quantity     int64                   `json:"quantity"`
}
