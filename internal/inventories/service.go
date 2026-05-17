package inventories

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api/internal/audit_logs"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/email"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"

	"github.com/redis/go-redis/v9"
)

type InventoryRepository interface {
	GetInventoryByUUID(ctx context.Context, uuid string) (*models.Inventory, error)
	GetInventories(ctx context.Context, filters dto.BaseFilters) ([]models.Inventory, int, error)
	CreateInventory(ctx context.Context, category *models.Inventory) error
	UpdateInventory(ctx context.Context, category *models.Inventory) error
	DeleteInventory(ctx context.Context, uuid string) error
	UpdateInventoryStatus(ctx context.Context, uuid string) error
}

type InventoryService interface {
	GetInventoryByUUID(ctx context.Context, uuid string) (*models.Inventory, error)
	GetInventories(ctx context.Context, q dto.Query) ([]models.Inventory, int, error)
	CreateInventory(ctx context.Context, req InventoryRequest, audit models.AuditLogRequest) error
	UpdateInventory(ctx context.Context, uuid string, req InventoryRequest, audit models.AuditLogRequest) error
	DeleteInventory(ctx context.Context, uuid string, audit models.AuditLogRequest) error
	UpdateInventoryStatus(ctx context.Context, uuid string, audit models.AuditLogRequest) error
}

type service struct {
	r        InventoryRepository
	es       *email.EmailService
	redis    *redis.Client
	auditLog audit_logs.AuditLogService
}

func NewInventoryService(r InventoryRepository, es *email.EmailService, redis *redis.Client, auditLog audit_logs.AuditLogService) *service {
	return &service{r: r, es: es, redis: redis, auditLog: auditLog}
}

func (s *service) GetInventoryByUUID(ctx context.Context, uuid string) (*models.Inventory, error) {
	return s.r.GetInventoryByUUID(ctx, uuid)
}

func (s *service) GetInventories(ctx context.Context, q dto.Query) ([]models.Inventory, int, error) {
	filters := q.SanitizeQuery([]string{"quantity", "low_stock_threshold"})

	return s.r.GetInventories(ctx, filters)
}

func (s *service) CreateInventory(ctx context.Context, req InventoryRequest, audit models.AuditLogRequest) error {
	inventory := &models.Inventory{}

	if req.ProductID != nil {
		inventory.ProductID = *req.ProductID
	}
	if req.Quantity != nil {
		inventory.Quantity = *req.Quantity
	}
	if req.LowStockThreshold != nil {
		inventory.LowStockThreshold = *req.LowStockThreshold
	}

	err := s.r.CreateInventory(ctx, inventory)
	s.auditLog.ParseAndCreateAuditLog(audit, inventory.UUID, "INVENTORY", nil, *inventory, err)

	return err
}

func (s *service) UpdateInventory(ctx context.Context, uuid string, req InventoryRequest, audit models.AuditLogRequest) error {
	inventory, err := s.r.GetInventoryByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	oldInventory := *inventory

	if req.ProductID != nil {
		inventory.ProductID = *req.ProductID
	}
	if req.Quantity != nil {
		inventory.Quantity = *req.Quantity
	}
	if req.LowStockThreshold != nil {
		inventory.LowStockThreshold = *req.LowStockThreshold
	}

	err = s.r.UpdateInventory(ctx, inventory)
	s.auditLog.ParseAndCreateAuditLog(audit, uuid, "INVENTORY", oldInventory, *inventory, err)

	return err
}

func (s *service) DeleteInventory(ctx context.Context, uuid string, audit models.AuditLogRequest) error {
	inventory, err := s.r.GetInventoryByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.DeleteInventory(ctx, uuid)
	s.auditLog.ParseAndCreateAuditLog(audit, uuid, "INVENTORY", nil, *inventory, err)

	return err
}

func (s *service) UpdateInventoryStatus(ctx context.Context, uuid string, audit models.AuditLogRequest) error {
	inventory, err := s.r.GetInventoryByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.UpdateInventoryStatus(ctx, uuid)
	s.auditLog.ParseAndCreateAuditLog(audit, uuid, "INVENTORY", nil, *inventory, err)

	return err
}
