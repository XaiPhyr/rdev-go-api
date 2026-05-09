package service

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api/internal/data"
	"github.com/XaiPhyr/rdev-go-api/internal/dto"

	"github.com/redis/go-redis/v9"
)

type InventoryService struct {
	r     *data.InventoryRepository
	es    *EmailService
	redis *redis.Client
}

func NewInventory(r *data.InventoryRepository, es *EmailService, redis *redis.Client) *InventoryService {
	return &InventoryService{r: r, es: es, redis: redis}
}

func (s *InventoryService) GetInventoryByUUID(ctx context.Context, uuid string) (*data.Inventory, error) {
	return s.r.GetInventoryByUUID(ctx, uuid)
}

func (s *InventoryService) GetInventories(ctx context.Context, q dto.Query) ([]data.Inventory, int, error) {
	filters := q.SanitizeQuery([]string{"quantity", "low_stock_threshold"})

	return s.r.GetInventories(ctx, filters)
}

func (s *InventoryService) UpdateInventory(ctx context.Context, uuid string, req dto.InventoryRequestUpdate) error {
	inventory, err := s.r.GetInventoryByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	if req.ProductID != nil {
		inventory.ProductID = *req.ProductID
	}
	if req.Quantity != nil {
		inventory.Quantity = *req.Quantity
	}
	if req.LowStockThreshold != nil {
		inventory.LowStockThreshold = *req.LowStockThreshold
	}

	return s.r.UpdateInventory(ctx, inventory)
}

func (s *InventoryService) DeleteInventory(ctx context.Context, uuid string) error {
	return s.r.DeleteInventory(ctx, uuid)
}

func (s *InventoryService) UpdateInventoryStatus(ctx context.Context, uuid string) error {
	return s.r.UpdateInventoryStatus(ctx, uuid)
}
