package audit_logs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
)

type AuditLogRepository interface {
	CreateAuditLog(ctx context.Context, auditLog models.AuditLog) error
}

type AuditLogService interface {
	CreateAuditLog(auditLog models.AuditLog) error
	QueAuditLog(ctx context.Context)
	ParseAndCreateAuditLog(audit models.AuditLogRequest, module_id, module string, beforeChange, afterChang interface{}, err error)
}

type service struct {
	r        AuditLogRepository
	auditQue chan models.AuditLog
}

func NewAuditLogService(r AuditLogRepository) AuditLogService {
	return &service{r: r, auditQue: make(chan models.AuditLog, 1000)}
}

func (s *service) QueAuditLog(ctx context.Context) {
	for {
		select {
		case log := <-s.auditQue:
			fmt.Println("QUE START...")
			if err := s.r.CreateAuditLog(ctx, log); err != nil {
				s.handleDeadLetter(log, err)
			}
		case <-ctx.Done():
			fmt.Println("QUE DONE...")
			return
		}
	}
}

func (s *service) CreateAuditLog(al models.AuditLog) error {
	s.auditQue <- al
	return nil
}

func (s *service) handleDeadLetter(al models.AuditLog, err error) {
	// catch data save to local file
	fmt.Printf("⚠️ DEAD LETTER TRIGGERED: %v\n", err)
}

func (s *service) ParseAndCreateAuditLog(audit models.AuditLogRequest, module_id, module string, beforeChange, afterChang interface{}, err error) {
	errMsg := ""
	status := http.StatusOK
	if err != nil {
		status = http.StatusBadRequest
		errMsg = err.Error()
	}

	al := models.AuditLog{
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

	s.auditQue <- al
}
