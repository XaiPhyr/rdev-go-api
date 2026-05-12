package dto

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
