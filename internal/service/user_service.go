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
	CreateUser(ctx context.Context, req dto.UserRequest) error
	UpdateUser(ctx context.Context, uuid string, req dto.UserRequest) error
	DeleteUser(ctx context.Context, uuid string) error
	UpdateUserStatus(ctx context.Context, uuid string) error
}

type userService struct {
	r     UserRepository
	es    *EmailService
	redis *redis.Client
}

func NewUserService(r UserRepository, es *EmailService, redis *redis.Client) *userService {
	return &userService{r: r, es: es, redis: redis}
}

func (s *userService) GetUserByUUID(ctx context.Context, uuid string) (*data.User, error) {
	return s.r.GetUserByUUID(ctx, uuid)
}

func (s *userService) GetUsers(ctx context.Context, q dto.Query) ([]data.User, int, error) {
	filters := q.SanitizeQuery([]string{"first_name", "last_name", "email", "username"})

	return s.r.GetUsers(ctx, filters)
}

func (s *userService) CreateUser(ctx context.Context, req dto.UserRequest) error {
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

	return s.r.CreateUser(ctx, user)
}

func (s *userService) UpdateUser(ctx context.Context, uuid string, req dto.UserRequest) error {
	user, err := s.r.GetUserByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Username != nil {
		user.Username = *req.Username
	}

	return s.r.UpdateUser(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, uuid string) error {
	return s.r.DeleteUser(ctx, uuid)
}

func (s *userService) UpdateUserStatus(ctx context.Context, uuid string) error {
	return s.r.UpdateUserStatus(ctx, uuid)
}
