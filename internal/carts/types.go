package carts

type CartRequest struct {
	CustomerID *int64            `json:"customer_id"`
	CartStatus *string           `json:"cart_status"`
	CartItem   []CartItemRequest `json:"cart_item,omitempty"`
}

type CartItemRequest struct {
	CartID           *int64  `json:"cart_id"`
	ProductID        *int64  `json:"product_id"`
	ProductName      *string `json:"product_name"`
	TransactionPrice *int64  `json:"transaction_price"`
	Quantity         *int64  `json:"quantity"`
}
