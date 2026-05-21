package orders

type OrderRequest struct {
	CustomerID  *int64  `json:"customer_id"`
	OrderStatus *string `json:"order_status"`
}
