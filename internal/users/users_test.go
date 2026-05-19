package users_test

import (
	"context"
	"sync"
	"testing"

	"github.com/XaiPhyr/rdev-go-api/internal/mocks"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
	"github.com/XaiPhyr/rdev-go-api/internal/users"
)

const UUID = "12345678-1234-5678-1234-567890123456"

type UserTest struct {
	GetUserByUUIDFunc            func(ctx context.Context, uuid string) (*users.User, error)
	GetUsersFunc                 func(ctx context.Context, filters dto.BaseFilters) ([]users.User, int, error)
	CreateUserFunc               func(ctx context.Context, category *users.User) error
	UpdateUserFunc               func(ctx context.Context, category *users.User) error
	DeleteUserFunc               func(ctx context.Context, uuid string) error
	UpdateUserStatusFunc         func(ctx context.Context, uuid string) error
	GetUserByUsernameOrEmailFunc func(ctx context.Context, identifier string) (*users.User, error)
	CheckUserPermissionFunc      func(ctx context.Context, userID int64, roleName string) ([]string, error)
}

// CheckUserPermission implements [users.UserRepository].
func (m *UserTest) CheckUserPermission(ctx context.Context, userID int64, roleName string) ([]string, error) {
	if m.CheckUserPermissionFunc != nil {
		return m.CheckUserPermissionFunc(ctx, userID, roleName)
	}

	return nil, nil
}
func (m *UserTest) GetUserByUsernameOrEmail(ctx context.Context, identifier string) (*users.User, error) {
	if m.GetUserByUsernameOrEmailFunc != nil {
		return m.GetUserByUsernameOrEmailFunc(ctx, identifier)
	}

	return nil, nil
}
func (m *UserTest) GetUserByUUID(ctx context.Context, uuid string) (*users.User, error) {
	if m.GetUserByUUIDFunc != nil {
		return m.GetUserByUUIDFunc(ctx, uuid)
	}

	return nil, nil
}
func (m *UserTest) GetUsers(ctx context.Context, q dto.BaseFilters) ([]users.User, int, error) {
	if m.GetUsersFunc != nil {
		return m.GetUsersFunc(ctx, q)
	}

	return nil, 0, nil
}
func (m *UserTest) CreateUser(ctx context.Context, sm *users.User) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, sm)
	}

	return nil
}
func (m *UserTest) UpdateUser(ctx context.Context, sm *users.User) error {
	if m.UpdateUserFunc != nil {
		sm.ID = 1
		return m.UpdateUserFunc(ctx, sm)
	}
	return nil
}
func (m *UserTest) DeleteUser(ctx context.Context, uuid string) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(ctx, uuid)
	}

	return nil
}
func (m *UserTest) UpdateUserStatus(ctx context.Context, uuid string) error {
	if m.UpdateUserStatusFunc != nil {
		return m.UpdateUserStatusFunc(ctx, uuid)
	}

	return nil
}

func TestUser(t *testing.T) {
	testUserRepo := &UserTest{}
	emailSvc := mocks.NewTestEmailService()
	_, auditLogSvc := mocks.NewTestAuditService()

	testUserSvc := users.NewUserService(testUserRepo, emailSvc, nil, auditLogSvc)

	t.Run("Get Users", func(t *testing.T) {
		testUserRepo.GetUsersFunc = func(ctx context.Context, q dto.BaseFilters) ([]users.User, int, error) {
			CheckUserQuery(t, q)
			return []users.User{{FirstName: "John"}}, 1, nil
		}

		query := dto.Query{Search: "test", Limit: 10, Offset: 2, Sort: "first_name ASC"}
		_, _, err := testUserSvc.GetUsers(context.Background(), query)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Get User By UUID", func(t *testing.T) {
		testUserRepo.GetUserByUUIDFunc = func(ctx context.Context, uuid string) (*users.User, error) {
			CheckUUID(t, uuid)
			return &users.User{FirstName: "John"}, nil
		}

		_, err := testUserSvc.GetUserByUUID(context.Background(), UUID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Create User", func(t *testing.T) {
		testUserRepo.CreateUserFunc = func(ctx context.Context, sm *users.User) error {
			CheckUser(t, sm)
			return nil
		}

		numRequest := 50
		var wg sync.WaitGroup

		for range numRequest {
			wg.Go(func() {
				first_name := "John"
				last_name := "Doe"
				req := users.UserRequest{FirstName: &first_name, LastName: &last_name}
				err := testUserSvc.CreateUser(context.Background(), req, models.AuditLogRequest{})

				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			})
		}

		wg.Wait()
	})

	t.Run("Update User", func(t *testing.T) {
		testUserRepo.UpdateUserFunc = func(ctx context.Context, sm *users.User) error {
			if sm.ID == 0 {
				t.Error("Expected sm ID to be populated")
			}
			CheckUser(t, sm)
			return nil
		}

		first_name := "John"
		last_name := "Doe"
		req := users.UserRequest{FirstName: &first_name, LastName: &last_name}
		err := testUserSvc.UpdateUser(context.Background(), CheckUUID(t, UUID), req, models.AuditLogRequest{})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Delete User", func(t *testing.T) {
		testUserRepo.DeleteUserFunc = func(ctx context.Context, uuid string) error {
			CheckUUID(t, uuid)
			return nil
		}

		err := testUserSvc.DeleteUser(context.Background(), CheckUUID(t, UUID), models.AuditLogRequest{})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Update User Status", func(t *testing.T) {
		testUserRepo.UpdateUserStatusFunc = func(ctx context.Context, uuid string) error {
			CheckUUID(t, uuid)
			return nil
		}

		err := testUserSvc.UpdateUserStatus(context.Background(), CheckUUID(t, UUID), models.AuditLogRequest{})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func CheckUUID(t testing.TB, uuid string) string {
	t.Helper()

	if uuid == "" {
		t.Error("Expected UUID to be provided")
	}

	return uuid
}

func CheckUser(t testing.TB, user *users.User) {
	t.Helper()

	if user.FirstName == "" {
		t.Error("Expected user first_name to be populated")
	}
	if user.LastName == "" {
		t.Error("Expected user last_name to be populated")
	}
}

func CheckUserQuery(t testing.TB, q dto.BaseFilters) {
	t.Helper()

	if q.Search != "test" {
		t.Errorf("Expected search filter to be 'test', got '%s'", q.Search)
	}
	if q.Page < 1 {
		t.Errorf("Expected page to be 1, got %d", q.Page)
	}
	if q.PageSize < 1 {
		t.Errorf("Expected page size to be 10, got %d", q.PageSize)
	}
	if q.Sort != "first_name ASC" {
		t.Errorf("Expected sort to be first_name ASC, got '%s'", q.Sort)
	}
}
