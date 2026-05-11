package service

import (
	"context"
	"strings"

	"github.com/XaiPhyr/rdev-go-api/internal/data"
	"github.com/XaiPhyr/rdev-go-api/internal/dto"

	"github.com/redis/go-redis/v9"
)

type StockMovementService struct {
	r     *data.StockMovementRepository
	es    *EmailService
	redis *redis.Client
}

func NewStockMovement(r *data.StockMovementRepository, es *EmailService, redis *redis.Client) *StockMovementService {
	return &StockMovementService{r: r, es: es, redis: redis}
}

func (s *StockMovementService) GetStockMovementByUUID(ctx context.Context, uuid string) (*data.StockMovement, error) {
	return s.r.GetStockMovementByUUID(ctx, uuid)
}

func (s *StockMovementService) GetStockMovements(ctx context.Context, q dto.Query) ([]data.StockMovement, int, error) {
	filters := q.SanitizeQuery([]string{"change_amount", "reason", "reference_id"})

	return s.r.GetStockMovements(ctx, filters)
}

func (s *StockMovementService) CreateStockMovement(ctx context.Context, req dto.StockMovementRequest) error {
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

	return s.r.CreateStockMovement(ctx, sm)
}

func (s *StockMovementService) UpdateStockMovement(ctx context.Context, uuid string, req dto.StockMovementRequest) error {
	sm, err := s.r.GetStockMovementByUUID(ctx, uuid)
	if err != nil {
		return err
	}

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

	return s.r.UpdateStockMovement(ctx, sm)
}

func (s *StockMovementService) DeleteStockMovement(ctx context.Context, uuid string) error {
	return s.r.DeleteStockMovement(ctx, uuid)
}

func (s *StockMovementService) UpdateStockMovementStatus(ctx context.Context, uuid string) error {
	return s.r.UpdateStockMovementStatus(ctx, uuid)
}
