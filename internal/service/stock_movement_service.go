package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/XaiPhyr/rdev-go-api/internal/data"
	"github.com/XaiPhyr/rdev-go-api/internal/dto"
	"github.com/xuri/excelize/v2"

	"github.com/redis/go-redis/v9"
)

type StockMovementRepository interface {
	GetStockMovementByUUID(ctx context.Context, uuid string) (*data.StockMovement, error)
	GetStockMovements(ctx context.Context, filters dto.BaseFilters) ([]data.StockMovement, int, error)
	CreateStockMovement(ctx context.Context, sm *data.StockMovement) error
	UpdateStockMovement(ctx context.Context, sm *data.StockMovement) error
	DeleteStockMovement(ctx context.Context, uuid string) error
	UpdateStockMovementStatus(ctx context.Context, uuid string) error
	ProcessBulkUpload(ctx context.Context, row [][]string) error
}

type StockMovementService interface {
	GetStockMovementByUUID(ctx context.Context, uuid string) (*data.StockMovement, error)
	GetStockMovements(ctx context.Context, q dto.Query) ([]data.StockMovement, int, error)
	CreateStockMovement(ctx context.Context, req dto.StockMovementRequest, audit dto.AuditLogRequest) error
	UpdateStockMovement(ctx context.Context, uuid string, req dto.StockMovementRequest, audit dto.AuditLogRequest) error
	DeleteStockMovement(ctx context.Context, uuid string, audit dto.AuditLogRequest) error
	UpdateStockMovementStatus(ctx context.Context, uuid string, audit dto.AuditLogRequest) error
	BulkUpload(ctx context.Context, fileHeader *multipart.FileHeader) error
	ProcessBulkUpload(ctx context.Context, fileName string) error
}

type stockMovementService struct {
	r        StockMovementRepository
	es       *EmailService
	redis    *redis.Client
	auditLog AuditLogService
}

func NewStockMovementService(r StockMovementRepository, es *EmailService, redis *redis.Client, auditLog AuditLogService) *stockMovementService {
	return &stockMovementService{r: r, es: es, redis: redis, auditLog: auditLog}
}

func (s *stockMovementService) GetStockMovementByUUID(ctx context.Context, uuid string) (*data.StockMovement, error) {
	return s.r.GetStockMovementByUUID(ctx, uuid)
}

func (s *stockMovementService) GetStockMovements(ctx context.Context, q dto.Query) ([]data.StockMovement, int, error) {
	filters := q.SanitizeQuery([]string{"change_amount", "reason", "reference_id"})

	return s.r.GetStockMovements(ctx, filters)
}

func (s *stockMovementService) CreateStockMovement(ctx context.Context, req dto.StockMovementRequest, audit dto.AuditLogRequest) error {
	sm := &data.StockMovement{}

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
	s.auditLog.CreateAuditLog(parseAuditLog(audit, sm.UUID, "STOCK_MOVEMENT", nil, *sm, err))

	return err
}

func (s *stockMovementService) UpdateStockMovement(ctx context.Context, uuid string, req dto.StockMovementRequest, audit dto.AuditLogRequest) error {
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
	s.auditLog.CreateAuditLog(parseAuditLog(audit, uuid, "STOCK_MOVEMENT", oldStockMovement, *sm, err))

	return err
}

func (s *stockMovementService) DeleteStockMovement(ctx context.Context, uuid string, audit dto.AuditLogRequest) error {
	sm, err := s.r.GetStockMovementByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.DeleteStockMovement(ctx, uuid)
	s.auditLog.CreateAuditLog(parseAuditLog(audit, uuid, "STOCK_MOVEMENT", nil, *sm, err))

	return err
}

func (s *stockMovementService) UpdateStockMovementStatus(ctx context.Context, uuid string, audit dto.AuditLogRequest) error {
	sm, err := s.r.GetStockMovementByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.UpdateStockMovementStatus(ctx, uuid)
	s.auditLog.CreateAuditLog(parseAuditLog(audit, uuid, "STOCK_MOVEMENT", nil, *sm, err))

	return err
}

func (s *stockMovementService) BulkUpload(ctx context.Context, fileHeader *multipart.FileHeader) error {
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

func (s *stockMovementService) ProcessBulkUpload(ctx context.Context, fileName string) error {
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
