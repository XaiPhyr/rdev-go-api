package data

import (
	"context"

	"github.com/uptrace/bun"
)

type AuditLogRepository struct {
	db *bun.DB
}

func NewAuditLogRepository(db *bun.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) CreateAuditLog(ctx context.Context, auditLog AuditLog) error {
	_, err := r.db.NewInsert().Model(&auditLog).Exec(ctx)

	return err
}
