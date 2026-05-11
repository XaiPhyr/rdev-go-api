package service

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api/internal/data"
	"github.com/XaiPhyr/rdev-go-api/internal/dto"

	"github.com/redis/go-redis/v9"
)

type InventoryRepository interface {
	GetInventoryByUUID(ctx context.Context, uuid string) (*data.Inventory, error)
	GetInventories(ctx context.Context, filters dto.BaseFilters) ([]data.Inventory, int, error)
	CreateInventory(ctx context.Context, category *data.Inventory) error
	UpdateInventory(ctx context.Context, category *data.Inventory) error
	DeleteInventory(ctx context.Context, uuid string) error
	UpdateInventoryStatus(ctx context.Context, uuid string) error
}

type InventoryService interface {
	GetInventoryByUUID(ctx context.Context, uuid string) (*data.Inventory, error)
	GetInventories(ctx context.Context, q dto.Query) ([]data.Inventory, int, error)
	CreateInventory(ctx context.Context, req dto.InventoryRequest) error
	UpdateInventory(ctx context.Context, uuid string, req dto.InventoryRequest) error
	DeleteInventory(ctx context.Context, uuid string) error
	UpdateInventoryStatus(ctx context.Context, uuid string) error
}

type inventoryService struct {
	r     InventoryRepository
	es    *EmailService
	redis *redis.Client
}

func NewInventoryService(r InventoryRepository, es *EmailService, redis *redis.Client) *inventoryService {
	return &inventoryService{r: r, es: es, redis: redis}
}

func (s *inventoryService) GetInventoryByUUID(ctx context.Context, uuid string) (*data.Inventory, error) {
	return s.r.GetInventoryByUUID(ctx, uuid)
}

func (s *inventoryService) GetInventories(ctx context.Context, q dto.Query) ([]data.Inventory, int, error) {
	filters := q.SanitizeQuery([]string{"quantity", "low_stock_threshold"})

	return s.r.GetInventories(ctx, filters)
}

func (s *inventoryService) CreateInventory(ctx context.Context, req dto.InventoryRequest) error {
	inventory := &data.Inventory{}

	if req.ProductID != nil {
		inventory.ProductID = *req.ProductID
	}
	if req.Quantity != nil {
		inventory.Quantity = *req.Quantity
	}
	if req.LowStockThreshold != nil {
		inventory.LowStockThreshold = *req.LowStockThreshold
	}

	return s.r.CreateInventory(ctx, inventory)
}

func (s *inventoryService) UpdateInventory(ctx context.Context, uuid string, req dto.InventoryRequest) error {
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

func (s *inventoryService) DeleteInventory(ctx context.Context, uuid string) error {
	return s.r.DeleteInventory(ctx, uuid)
}

func (s *inventoryService) UpdateInventoryStatus(ctx context.Context, uuid string) error {
	return s.r.UpdateInventoryStatus(ctx, uuid)
}
