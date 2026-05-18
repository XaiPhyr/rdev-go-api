package mocks

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api/internal/audit_logs"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
)

type AuditLogTest struct{}

func (m AuditLogTest) CreateAuditLog(ctx context.Context, auditLog models.AuditLog) error {
	return nil
}

func NewTestAuditService() (AuditLogTest, audit_logs.AuditLogService) {
	mockRepo := AuditLogTest{}
	service := audit_logs.NewAuditLogService(mockRepo)
	return mockRepo, service
}
