package service

import (
	"context"
	"rdev-go-api/internal/data"
	"rdev-go-api/internal/dto"
)

type UserService struct {
	r *data.UserRepository
}

func NewUserService(r *data.UserRepository) *UserService {
	return &UserService{r: r}
}

func (s *UserService) GetUserByUUID(ctx context.Context, uuid string) (*data.User, error) {
	return s.r.GetUserByUUID(ctx, uuid)
}

func (s *UserService) GetUsers(ctx context.Context, q dto.UserQuery) ([]data.User, int, error) {
	orderBy := "id ASC"
	if q.Sort != "" {
		orderBy = q.Sort
	}

	filters := data.UserFilters{
		Limit:  q.Limit,
		Offset: q.Offset,
		Order:  orderBy,
		Search: q.Search,
	}

	return s.r.GetUsers(ctx, filters)
}

func (s *UserService) UpdateUser(ctx context.Context, uuid string, req dto.UserRequestUpdate) error {
	user, err := s.r.GetUserByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Username = req.Username

	return s.r.UpdateUser(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, uuid string) error {
	return s.r.DeleteUser(ctx, uuid)
}
