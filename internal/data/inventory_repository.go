package data

import (
	"context"
	"time"

	"github.com/XaiPhyr/rdev-go-api/internal/dto"
	"github.com/uptrace/bun"
)

type InventoryRepository struct {
	db *bun.DB
}

func NewInventoryRepository(db *bun.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) GetInventoryByUUID(ctx context.Context, uuid string) (*Inventory, error) {
	inventory := new(Inventory)

	err := r.db.NewSelect().
		Model(inventory).
		Where("uuid = ?", uuid).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return inventory, nil
}

func (r *InventoryRepository) GetInventories(ctx context.Context, q dto.BaseFilters) ([]Inventory, int, error) {
	var inventories []Inventory

	count, err := r.db.NewSelect().
		Model(&inventories).
		Limit(q.PageSize).
		Offset(q.Page).
		Order(q.Sort).
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	return inventories, count, nil
}

func (r *InventoryRepository) CreateInventory(ctx context.Context, inventory *Inventory) error {
	_, err := r.db.NewInsert().Model(inventory).Exec(ctx)

	return err
}

func (r *InventoryRepository) UpdateInventory(ctx context.Context, inventory *Inventory) error {
	_, err := r.db.NewUpdate().
		Model(inventory).
		Column("product_id", "quantity", "low_stock_threshold").
		Set("updated_at = ?", time.Now()).
		WherePK().
		Exec(ctx)

	return err
}

func (r *InventoryRepository) DeleteInventory(ctx context.Context, uuid string) error {
	_, err := r.db.NewDelete().
		Model((*Inventory)(nil)).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}

func (r *InventoryRepository) UpdateInventoryStatus(ctx context.Context, uuid string) error {
	_, err := r.db.NewUpdate().
		Model((*Inventory)(nil)).
		Set("status = CASE WHEN status = 'A' THEN 'I' ELSE 'A' END").
		Set("updated_at = ?", time.Now()).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}
