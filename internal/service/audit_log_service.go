package service

import (
	"context"
	"fmt"

	"github.com/XaiPhyr/rdev-go-api/internal/data"
)

type AuditLogRepository interface {
	CreateAuditLog(ctx context.Context, auditLog data.AuditLog) error
}

type AuditLogService interface {
	CreateAuditLog(auditLog data.AuditLog) error
	QueAuditLog(ctx context.Context)
}

type auditLogService struct {
	r        AuditLogRepository
	auditQue chan data.AuditLog
}

func NewAuditLogService(r AuditLogRepository) AuditLogService {
	return &auditLogService{r: r, auditQue: make(chan data.AuditLog, 1000)}
}

func (s *auditLogService) QueAuditLog(ctx context.Context) {
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

func (s *auditLogService) CreateAuditLog(al data.AuditLog) error {
	s.auditQue <- al
	return nil
}

func (s *auditLogService) handleDeadLetter(al data.AuditLog, err error) {
	// catch data save to local file
	fmt.Printf("⚠️ DEAD LETTER TRIGGERED: %v\n", err)
}
