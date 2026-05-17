package stock_movements

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/XaiPhyr/rdev-go-api/internal/audit_logs"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/email"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
	"github.com/xuri/excelize/v2"

	"github.com/redis/go-redis/v9"
)

type StockMovementRepository interface {
	GetStockMovementByUUID(ctx context.Context, uuid string) (*StockMovement, error)
	GetStockMovements(ctx context.Context, filters dto.BaseFilters) ([]StockMovement, int, error)
	CreateStockMovement(ctx context.Context, sm *StockMovement) error
	UpdateStockMovement(ctx context.Context, sm *StockMovement) error
	DeleteStockMovement(ctx context.Context, uuid string) error
	UpdateStockMovementStatus(ctx context.Context, uuid string) error
	ProcessBulkUpload(ctx context.Context, row [][]string) error
}

type StockMovementService interface {
	GetStockMovementByUUID(ctx context.Context, uuid string) (*StockMovement, error)
	GetStockMovements(ctx context.Context, q dto.Query) ([]StockMovement, int, error)
	CreateStockMovement(ctx context.Context, req StockMovementRequest, audit models.AuditLogRequest) error
	UpdateStockMovement(ctx context.Context, uuid string, req StockMovementRequest, audit models.AuditLogRequest) error
	DeleteStockMovement(ctx context.Context, uuid string, audit models.AuditLogRequest) error
	UpdateStockMovementStatus(ctx context.Context, uuid string, audit models.AuditLogRequest) error
	BulkUpload(ctx context.Context, fileHeader *multipart.FileHeader, audit models.AuditLogRequest) error
	ProcessBulkUpload(ctx context.Context, fileName string, audit models.AuditLogRequest) error
}

type service struct {
	r        StockMovementRepository
	es       *email.EmailService
	redis    *redis.Client
	auditLog audit_logs.AuditLogService
}

func NewStockMovementService(r StockMovementRepository, es *email.EmailService, redis *redis.Client, auditLog audit_logs.AuditLogService) *service {
	return &service{r: r, es: es, redis: redis, auditLog: auditLog}
}

func (s *service) GetStockMovementByUUID(ctx context.Context, uuid string) (*StockMovement, error) {
	return s.r.GetStockMovementByUUID(ctx, uuid)
}

func (s *service) GetStockMovements(ctx context.Context, q dto.Query) ([]StockMovement, int, error) {
	filters := q.SanitizeQuery([]string{"change_amount", "reason", "reference_id"})

	return s.r.GetStockMovements(ctx, filters)
}

func (s *service) CreateStockMovement(ctx context.Context, req StockMovementRequest, audit models.AuditLogRequest) error {
	sm := &StockMovement{}

	if req.ProductID != nil {
		sm.ProductID = *req.ProductID
	}
	if req.ChangeAmount != nil {
		sm.ChangeAmount = *req.ChangeAmount
	}
	if req.Reason != nil {
		sm.Reason = strings.ToUpper(*req.Reason)
	}
	if req.ReferenceID != nil {
		sm.ReferenceID = *req.ReferenceID
	}

	err := s.r.CreateStockMovement(ctx, sm)
	s.auditLog.CreateAuditLog(s.auditLog.ParseAuditLog(audit, sm.UUID, "STOCK_MOVEMENT", nil, *sm, err))

	return err
}

func (s *service) UpdateStockMovement(ctx context.Context, uuid string, req StockMovementRequest, audit models.AuditLogRequest) error {
	sm, err := s.r.GetStockMovementByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	oldStockMovement := *sm

	if req.ProductID != nil {
		sm.ProductID = *req.ProductID
	}
	if req.ChangeAmount != nil {
		sm.ChangeAmount = *req.ChangeAmount
	}
	if req.Reason != nil {
		sm.Reason = *req.Reason
	}
	if req.ReferenceID != nil {
		sm.ReferenceID = *req.ReferenceID
	}

	err = s.r.UpdateStockMovement(ctx, sm)
	s.auditLog.CreateAuditLog(s.auditLog.ParseAuditLog(audit, uuid, "STOCK_MOVEMENT", oldStockMovement, *sm, err))

	return err
}

func (s *service) DeleteStockMovement(ctx context.Context, uuid string, audit models.AuditLogRequest) error {
	sm, err := s.r.GetStockMovementByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.DeleteStockMovement(ctx, uuid)
	s.auditLog.CreateAuditLog(s.auditLog.ParseAuditLog(audit, uuid, "STOCK_MOVEMENT", nil, *sm, err))

	return err
}

func (s *service) UpdateStockMovementStatus(ctx context.Context, uuid string, audit models.AuditLogRequest) error {
	sm, err := s.r.GetStockMovementByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.UpdateStockMovementStatus(ctx, uuid)
	s.auditLog.CreateAuditLog(s.auditLog.ParseAuditLog(audit, uuid, "STOCK_MOVEMENT", nil, *sm, err))

	return err
}

func (s *service) BulkUpload(ctx context.Context, fileHeader *multipart.FileHeader, audit models.AuditLogRequest) error {
	// File saving should be separate e.g. AWS S3
	// Move config to internal/config
	// cfg, err := config.LoadDefaultConfig(context.TODO())
	// if err != nil {
	// 	log.Printf("error: %v", err)
	// 	return
	// }

	// client := s3.NewFromConfig(cfg)

	// uploader := manager.NewUploader(client)
	// result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
	// 	Bucket: aws.String("amzn-s3-demo-bucket"),
	// 	Key:    aws.String("my-object-key"),
	// 	Body:   uploadFile,
	// })

	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".xlsx" {
		return fmt.Errorf("could not proccess file format: %s", ext)
	}

	if _, err := os.Stat("./files/"); os.IsNotExist(err) {
		err := os.Mkdir("./files/", 0755)
		if err != nil {
			return fmt.Errorf("could not create directory: %w", err)
		}
	}

	source, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer source.Close()

	// should save filename as unique id
	path := filepath.Join("./files/", filepath.Base(fileHeader.Filename))
	destination, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not save file: %w", err)
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	// should have a db to save on

	return err
}

func (s *service) ProcessBulkUpload(ctx context.Context, fileName string, audit models.AuditLogRequest) error {
	f, err := excelize.OpenFile(filepath.Join("./files/", filepath.Base(fileName)))
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return fmt.Errorf("could not open sheet: %w", err)
	}

	err = s.r.ProcessBulkUpload(ctx, rows)

	return err
}
