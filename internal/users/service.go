package users

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api/internal/audit_logs"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/email"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"

	"github.com/redis/go-redis/v9"
)

type UserRepository interface {
	GetUserByUUID(ctx context.Context, uuid string) (*User, error)
	GetUsers(ctx context.Context, filters dto.BaseFilters) ([]User, int, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, uuid string) error
	UpdateUserStatus(ctx context.Context, uuid string) error
	GetUserByUsernameOrEmail(ctx context.Context, identifier string) (*User, error)
	CheckUserPermission(ctx context.Context, userID int64, roleName string) ([]string, error)
}

type UserService interface {
	GetUserByUUID(ctx context.Context, uuid string) (*User, error)
	GetUsers(ctx context.Context, q dto.Query) ([]User, int, error)
	CreateUser(ctx context.Context, req UserRequest, audit models.AuditLogRequest) error
	UpdateUser(ctx context.Context, uuid string, req UserRequest, audit models.AuditLogRequest) error
	DeleteUser(ctx context.Context, uuid string, audit models.AuditLogRequest) error
	UpdateUserStatus(ctx context.Context, uuid string, audit models.AuditLogRequest) error
}

type service struct {
	r        UserRepository
	es       email.EmailService
	redis    *redis.Client
	auditLog audit_logs.AuditLogService
}

func NewUserService(r UserRepository, es email.EmailService, redis *redis.Client, auditLog audit_logs.AuditLogService) *service {
	return &service{r: r, es: es, redis: redis, auditLog: auditLog}
}

func (s *service) GetUserByUUID(ctx context.Context, uuid string) (*User, error) {
	return s.r.GetUserByUUID(ctx, uuid)
}

func (s *service) GetUsers(ctx context.Context, q dto.Query) ([]User, int, error) {
	filters := q.SanitizeQuery([]string{"first_name", "last_name", "email", "username"})

	return s.r.GetUsers(ctx, filters)
}

func (s *service) CreateUser(ctx context.Context, req UserRequest, audit models.AuditLogRequest) error {
	user := &User{}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}

	err := s.r.CreateUser(ctx, user)
	s.auditLog.ParseAndCreateAuditLog(audit, user.UUID, "USER", nil, *user, err)

	return err
}

func (s *service) UpdateUser(ctx context.Context, uuid string, req UserRequest, audit models.AuditLogRequest) error {
	user, err := s.r.GetUserByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	oldUser := *user

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Username != nil {
		user.Username = *req.Username
	}

	err = s.r.UpdateUser(ctx, user)
	s.auditLog.ParseAndCreateAuditLog(audit, user.UUID, "USER", oldUser, *user, err)

	return err
}

func (s *service) DeleteUser(ctx context.Context, uuid string, audit models.AuditLogRequest) error {
	user, err := s.r.GetUserByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.DeleteUser(ctx, uuid)
	s.auditLog.ParseAndCreateAuditLog(audit, uuid, "USER", nil, user, err)

	return err
}

func (s *service) UpdateUserStatus(ctx context.Context, uuid string, audit models.AuditLogRequest) error {
	user, err := s.r.GetUserByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.UpdateUserStatus(ctx, uuid)
	s.auditLog.ParseAndCreateAuditLog(audit, uuid, "USER", nil, user, err)

	return err
}
