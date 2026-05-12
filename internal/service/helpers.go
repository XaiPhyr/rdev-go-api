package service

import (
	"net/http"

	"github.com/XaiPhyr/rdev-go-api/internal/data"
	"github.com/XaiPhyr/rdev-go-api/internal/dto"
)

func parseAuditLog(audit dto.AuditLogRequest, module_id, module string, beforeChange, afterChang interface{}, err error) data.AuditLog {
	errMsg := ""
	status := http.StatusOK
	if err != nil {
		status = http.StatusBadRequest
		errMsg = err.Error()
	}

	return data.AuditLog{
		UserID:         audit.UserID,
		Path:           audit.Path,
		Action:         audit.Action,
		ResponseStatus: status,
		ModuleID:       module_id,
		Module:         module,
		BeforeChange:   beforeChange,
		AfterChange:    afterChang,
		IPAddress:      audit.IPAddress,
		UserAgent:      audit.UserAgent,
		ErrorMessage:   errMsg,
	}
}
