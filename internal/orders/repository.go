package orders

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetOrderByUUID(ctx context.Context, uuid string) (*models.Order, error) {
	return nil, nil
}

func (r *Repository) GetOrders(ctx context.Context, q dto.BaseFilters) ([]models.Order, int, error) {
	return nil, 0, nil
}

func (r *Repository) CreateOrder(ctx context.Context, order *models.Order) error {
	return nil
}

func (r *Repository) UpdateOrder(ctx context.Context, order *models.Order) error {
	return nil
}

func (r *Repository) DeleteOrder(ctx context.Context, uuid string) error {
	return nil
}

func (r *Repository) UpdateOrderStatus(ctx context.Context, uuid string) error {
	return nil
}
