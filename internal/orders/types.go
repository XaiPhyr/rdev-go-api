package orders

type OrderRequest struct {
	CustomerID  *int64             `json:"customer_id"`
	OrderStatus *string            `json:"order_status"`
	TotalAmount *int64             `json:"total_amount"`
	OrderItem   []OrderItemRequest `json:"order_item"`
}

type OrderItemRequest struct {
	ProductID        *int64 `json:"product_id"`
	TransactionPrice *int64 `json:"transaction_price"`
	Quantity         *int64 `json:"quantity"`
}
