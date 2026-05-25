package carts

type CartRequest struct {
	CustomerID *int64            `json:"customer_id"`
	CartStatus *string           `json:"cart_status"`
	CartItem   []CartItemRequest `json:"cart_item,omitempty"`
}

type CartItemRequest struct {
	ProductName      *string `json:"product_name"`
	TransactionPrice *int64  `json:"transaction_price"`
	Quantity         *int64  `json:"quantity"`
}
