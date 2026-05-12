package service

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api/internal/data"
	"github.com/XaiPhyr/rdev-go-api/internal/dto"

	"github.com/redis/go-redis/v9"
)

type UserRepository interface {
	GetUserByUUID(ctx context.Context, uuid string) (*data.User, error)
	GetUsers(ctx context.Context, filters dto.BaseFilters) ([]data.User, int, error)
	CreateUser(ctx context.Context, category *data.User) error
	UpdateUser(ctx context.Context, category *data.User) error
	DeleteUser(ctx context.Context, uuid string) error
	UpdateUserStatus(ctx context.Context, uuid string) error
}

type UserService interface {
	GetUserByUUID(ctx context.Context, uuid string) (*data.User, error)
	GetUsers(ctx context.Context, q dto.Query) ([]data.User, int, error)
	CreateUser(ctx context.Context, req dto.UserRequest, audit dto.AuditLogRequest) error
	UpdateUser(ctx context.Context, uuid string, req dto.UserRequest, audit dto.AuditLogRequest) error
	DeleteUser(ctx context.Context, uuid string, audit dto.AuditLogRequest) error
	UpdateUserStatus(ctx context.Context, uuid string, audit dto.AuditLogRequest) error
}

type userService struct {
	r        UserRepository
	es       *EmailService
	redis    *redis.Client
	auditLog AuditLogService
}

func NewUserService(r UserRepository, es *EmailService, redis *redis.Client, auditLog AuditLogService) *userService {
	return &userService{r: r, es: es, redis: redis, auditLog: auditLog}
}

func (s *userService) GetUserByUUID(ctx context.Context, uuid string) (*data.User, error) {
	return s.r.GetUserByUUID(ctx, uuid)
}

func (s *userService) GetUsers(ctx context.Context, q dto.Query) ([]data.User, int, error) {
	filters := q.SanitizeQuery([]string{"first_name", "last_name", "email", "username"})

	return s.r.GetUsers(ctx, filters)
}

func (s *userService) CreateUser(ctx context.Context, req dto.UserRequest, audit dto.AuditLogRequest) error {
	user := &data.User{}

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
	s.auditLog.CreateAuditLog(parseAuditLog(audit, user.UUID, "USER", nil, *user, err))

	return err
}

func (s *userService) UpdateUser(ctx context.Context, uuid string, req dto.UserRequest, audit dto.AuditLogRequest) error {
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
	s.auditLog.CreateAuditLog(parseAuditLog(audit, user.UUID, "USER", oldUser, *user, err))

	return err
}

func (s *userService) DeleteUser(ctx context.Context, uuid string, audit dto.AuditLogRequest) error {
	user, err := s.r.GetUserByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.DeleteUser(ctx, uuid)
	s.auditLog.CreateAuditLog(parseAuditLog(audit, uuid, "USER", nil, user, err))

	return err
}

func (s *userService) UpdateUserStatus(ctx context.Context, uuid string, audit dto.AuditLogRequest) error {
	user, err := s.r.GetUserByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.r.UpdateUserStatus(ctx, uuid)
	s.auditLog.CreateAuditLog(parseAuditLog(audit, uuid, "USER", nil, user, err))

	return err
}
