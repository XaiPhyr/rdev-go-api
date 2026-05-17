package categories

import (
	"context"
	"time"

	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewCategoryRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetCategoryByUUID(ctx context.Context, uuid string) (*Category, error) {
	category := new(Category)

	err := r.db.NewSelect().
		Model(category).
		Where("uuid = ?", uuid).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *Repository) GetCategories(ctx context.Context, q dto.BaseFilters) ([]Category, int, error) {
	var categories []Category

	count, err := r.db.NewSelect().
		Model(&categories).
		Limit(q.PageSize).
		Offset(q.Page).
		Order(q.Sort).
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}

func (r *Repository) GetCategoryTree(ctx context.Context, q dto.BaseFilters) ([]CategoryTree, error) {
	var categories []CategoryTree

	query := `WITH RECURSIVE category_tree AS (
		SELECT id, name, slug, parent_id, CAST(name as TEXT) as full_path, 1 AS depth 
		FROM categories
		UNION ALL 
		SELECT c.id, c.name, c.slug, c.parent_id, ct.full_path || ' > ' || c.name, ct.depth + 1 
		FROM categories c 
		INNER JOIN category_tree ct ON c.parent_id = ct.id
	)
	SELECT * FROM category_tree ORDER BY depth ASC`

	err := r.db.NewRaw(query).Scan(ctx, &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *Repository) CreateCategory(ctx context.Context, category *Category) error {
	_, err := r.db.NewInsert().Model(category).Exec(ctx)

	return err
}

func (r *Repository) UpdateCategory(ctx context.Context, category *Category) error {
	_, err := r.db.NewUpdate().
		Model(category).
		Column("parent_id", "name", "slug").
		Set("updated_at = ?", time.Now()).
		WherePK().
		Exec(ctx)

	return err
}

func (r *Repository) DeleteCategory(ctx context.Context, uuid string) error {
	_, err := r.db.NewDelete().
		Model((*Category)(nil)).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}

func (r *Repository) UpdateCategoryStatus(ctx context.Context, uuid string) error {
	_, err := r.db.NewUpdate().
		Model((*Category)(nil)).
		Set("status = CASE WHEN status = 'A' THEN 'I' ELSE 'A' END").
		Set("updated_at = ?", time.Now()).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}
