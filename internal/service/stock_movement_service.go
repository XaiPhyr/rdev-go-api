package service

import (
	"context"
	"strings"

	"github.com/XaiPhyr/rdev-go-api/internal/data"
	"github.com/XaiPhyr/rdev-go-api/internal/dto"

	"github.com/redis/go-redis/v9"
)

type StockMovementRepository interface {
	GetStockMovementByUUID(ctx context.Context, uuid string) (*data.StockMovement, error)
	GetStockMovements(ctx context.Context, filters dto.BaseFilters) ([]data.StockMovement, int, error)
	CreateStockMovement(ctx context.Context, sm *data.StockMovement) error
	UpdateStockMovement(ctx context.Context, sm *data.StockMovement) error
	DeleteStockMovement(ctx context.Context, uuid string) error
	UpdateStockMovementStatus(ctx context.Context, uuid string) error
}

type StockMovementService interface {
	GetStockMovementByUUID(ctx context.Context, uuid string) (*data.StockMovement, error)
	GetStockMovements(ctx context.Context, q dto.Query) ([]data.StockMovement, int, error)
	CreateStockMovement(ctx context.Context, req dto.StockMovementRequest, audit dto.AuditLogRequest) error
	UpdateStockMovement(ctx context.Context, uuid string, req dto.StockMovementRequest, audit dto.AuditLogRequest) error
	DeleteStockMovement(ctx context.Context, uuid string, audit dto.AuditLogRequest) error
	UpdateStockMovementStatus(ctx context.Context, uuid string, audit dto.AuditLogRequest) error
}

type stockMovementService struct {
	r        StockMovementRepository
	es       *EmailService
	redis    *redis.Client
	auditLog AuditLogService
}

func NewStockMovementService(r StockMovementRepository, es *EmailService, redis *redis.Client, auditLog AuditLogService) *stockMovementService {
	return &stockMovementService{r: r, es: es, redis: redis, auditLog: auditLog}
}

func (s *stockMovementService) GetStockMovementByUUID(ctx context.Context, uuid string) (*data.StockMovement, error) {
	return s.r.GetStockMovementByUUID(ctx, uuid)
}

func (s *stockMovementService) GetStockMovements(ctx context.Context, q dto.Query) ([]data.StockMovement, int, error) {
	filters := q.SanitizeQuery([]string{"change_amount", "reason", "reference_id"})

	return s.r.GetStockMovements(ctx, filters)
}

func (s *stockMovementService) CreateStockMovement(ctx context.Context, req dto.StockMovementRequest, audit dto.AuditLogRequest) error {
	sm := &data.StockMovement{}

	if req.ProductID != nil {
		sm.ProductID = *req.ProductID
	}
	if req.ChangeAmount != nil {
		sm.ChangeAmount = *req.ChangeAmount
	}
	if req.Reason != nil {
		sm.Reason = strings.ToUpper(*req.Reason)
	}
	if req.ReferenceID != nil {
		sm.ReferenceID = *req.ReferenceID
	}

	err := s.r.CreateStockMovement(ctx, sm)
	s.auditLog.CreateAuditLog(parseAuditLog(audit, sm.UUID, "STOCK_MOVEMENT", nil, *sm, err))

	return err
}

func (s *stockMovementService) UpdateStockMovement(ctx context.Context, uuid string, req dto.StockMovementRequest, audit dto.AuditLogRequest) error {
	sm, err := s.r.GetStockMovementByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	oldStockMovement := *sm

	if req.ProductID != nil {
		sm.ProductID = *req.ProductID
	}
	if req.ChangeAmount != nil {
		sm.ChangeAmount = *req.ChangeAmount
	}
	if req.Reason != nil {
		sm.Reason = *req.Reason
	}
	if req.ReferenceID != nil {
		sm.ReferenceID = *req.ReferenceID
	}

	err = s.r.UpdateStockMovement(ctx, sm)
	s.auditLog.CreateAuditLog(parseAuditLog(audit, uuid, "STOCK_MOVEMENT", oldStockMovement, *sm, err))

	return err
}

func (s *stockMovementService) DeleteStockMovement(ctx context.Context, uuid string, audit dto.AuditLogRequest) error {
	sm, err := s.r.GetStockMovementByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.DeleteStockMovement(ctx, uuid)
	s.auditLog.CreateAuditLog(parseAuditLog(audit, uuid, "STOCK_MOVEMENT", nil, *sm, err))

	return err
}

func (s *stockMovementService) UpdateStockMovementStatus(ctx context.Context, uuid string, audit dto.AuditLogRequest) error {
	sm, err := s.r.GetStockMovementByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.UpdateStockMovementStatus(ctx, uuid)
	s.auditLog.CreateAuditLog(parseAuditLog(audit, uuid, "STOCK_MOVEMENT", nil, *sm, err))

	return err
}
