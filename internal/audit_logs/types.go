package audit_logs

import (
	"github.com/XaiPhyr/rdev-go-api/internal/shared/fields"
	"github.com/uptrace/bun"
)

type AuditLog struct {
	bun.BaseModel `bun:"table:audit_logs,alias:al"`
	fields.BaseFields

	UserID         int64  `bun:"user_id" json:"user_id"`
	Path           string `bun:"path" json:"path"`
	Action         string `bun:"action" json:"action"`
	ResponseStatus int    `bun:"response_status" json:"response_status"`
	ModuleID       string `bun:"module_id" json:"module_id"`
	Module         string `bun:"module" json:"module"`
	BeforeChange   any    `bun:"before_change" json:"before_change"`
	AfterChange    any    `bun:"after_change" json:"after_change"`
	IPAddress      string `bun:"ip_address" json:"ip_address"`
	UserAgent      string `bun:"user_agent" json:"user_agent"`
	ErrorMessage   string `bun:"error_message" json:"error_message"`
}

type AuditLogRequest struct {
	UserID         int64  `json:"user_id"`
	Path           string `json:"path"`
	Action         string `json:"action"`
	ResponseStatus int    `json:"response_status"`
	ModuleID       string `json:"module_id"`
	Module         string `json:"module"`
	BeforeChange   any    `json:"before_change"`
	AfterChange    any    `json:"after_change"`
	IPAddress      string `json:"ip_address"`
	UserAgent      string `json:"user_agent"`
	ErrorMessage   string `json:"error_message"`
}
