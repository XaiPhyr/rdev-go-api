package carts

type Handler struct {
	svc CartService
}

func NewCartHandler(svc CartService) *Handler {
	return &Handler{svc: svc}
}
