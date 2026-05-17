package audit_logs

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewAuditLogRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAuditLog(ctx context.Context, auditLog models.AuditLog) error {
	_, err := r.db.NewInsert().Model(&auditLog).Exec(ctx)

	return err
}
