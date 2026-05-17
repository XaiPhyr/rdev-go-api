package inventories

import (
	"context"
	"time"

	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewInventoryRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetInventoryByUUID(ctx context.Context, uuid string) (*models.Inventory, error) {
	inventory := new(models.Inventory)

	err := r.db.NewSelect().
		Model(inventory).
		Where("uuid = ?", uuid).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return inventory, nil
}

func (r *Repository) GetInventories(ctx context.Context, q dto.BaseFilters) ([]models.Inventory, int, error) {
	var inventories []models.Inventory

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

func (r *Repository) CreateInventory(ctx context.Context, inventory *models.Inventory) error {
	_, err := r.db.NewInsert().Model(inventory).Exec(ctx)

	return err
}

func (r *Repository) UpdateInventory(ctx context.Context, inventory *models.Inventory) error {
	_, err := r.db.NewUpdate().
		Model(inventory).
		Column("product_id", "quantity", "low_stock_threshold").
		Set("updated_at = ?", time.Now()).
		WherePK().
		Exec(ctx)

	return err
}

func (r *Repository) DeleteInventory(ctx context.Context, uuid string) error {
	_, err := r.db.NewDelete().
		Model((*models.Inventory)(nil)).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}

func (r *Repository) UpdateInventoryStatus(ctx context.Context, uuid string) error {
	_, err := r.db.NewUpdate().
		Model((*models.Inventory)(nil)).
		Set("status = CASE WHEN status = 'A' THEN 'I' ELSE 'A' END").
		Set("updated_at = ?", time.Now()).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}
