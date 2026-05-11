package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/XaiPhyr/rdev-go-api/internal/dto"
	"github.com/uptrace/bun"
)

type StockMovementRepository struct {
	db *bun.DB
}

func NewStockMovementRepository(db *bun.DB) *StockMovementRepository {
	return &StockMovementRepository{db: db}
}

func (r *StockMovementRepository) GetStockMovementByUUID(ctx context.Context, uuid string) (*StockMovement, error) {
	sm := new(StockMovement)

	err := r.db.NewSelect().
		Model(sm).
		Where("uuid = ?", uuid).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return sm, nil
}

func (r *StockMovementRepository) GetStockMovements(ctx context.Context, q dto.BaseFilters) ([]StockMovement, int, error) {
	var sm []StockMovement

	count, err := r.db.NewSelect().
		Model(&sm).
		Limit(q.PageSize).
		Offset(q.Page).
		Order(q.Sort).
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	return sm, count, nil
}

func (r *StockMovementRepository) CreateStockMovement(ctx context.Context, sm *StockMovement) error {
	res, err := r.db.NewInsert().Model(sm).Exec(ctx)

	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return err
}

func (r *StockMovementRepository) UpdateStockMovement(ctx context.Context, sm *StockMovement) error {
	res, err := r.db.NewUpdate().
		Model(sm).
		Column("product_id", "quantity", "low_stock_threshold").
		Set("updated_at = ?", time.Now()).
		WherePK().
		Exec(ctx)

	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return err
}

func (r *StockMovementRepository) DeleteStockMovement(ctx context.Context, uuid string) error {
	res, err := r.db.NewDelete().
		Model((*StockMovement)(nil)).
		Where("uuid = ?", uuid).
		Exec(ctx)

	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return err
}

func (r *StockMovementRepository) UpdateStockMovementStatus(ctx context.Context, uuid string) error {
	res, err := r.db.NewUpdate().
		Model((*StockMovement)(nil)).
		Set("status = CASE WHEN status = 'A' THEN 'I' ELSE 'A' END").
		Set("updated_at = ?", time.Now()).
		Where("uuid = ?", uuid).
		Exec(ctx)

	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return err
}
