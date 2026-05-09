package service

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api/internal/data"
	"github.com/XaiPhyr/rdev-go-api/internal/dto"

	"github.com/redis/go-redis/v9"
)

type UserService struct {
	r     *data.UserRepository
	es    *EmailService
	redis *redis.Client
}

func NewUserService(r *data.UserRepository, es *EmailService, redis *redis.Client) *UserService {
	return &UserService{r: r, es: es, redis: redis}
}

func (s *UserService) GetUserByUUID(ctx context.Context, uuid string) (*data.User, error) {
	return s.r.GetUserByUUID(ctx, uuid)
}

func (s *UserService) GetUsers(ctx context.Context, q dto.Query) ([]data.User, int, error) {
	filters := q.SanitizeQuery([]string{"first_name", "last_name", "email", "username"})

	return s.r.GetUsers(ctx, filters)
}

func (s *UserService) UpdateUser(ctx context.Context, uuid string, req dto.UserRequestUpdate) error {
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

func (s *UserService) DeleteUser(ctx context.Context, uuid string) error {
	return s.r.DeleteUser(ctx, uuid)
}

func (s *UserService) UpdateUserStatus(ctx context.Context, uuid string) error {
	return s.r.UpdateUserStatus(ctx, uuid)
}
