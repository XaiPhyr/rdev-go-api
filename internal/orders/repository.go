package orders

import (
	"context"

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
