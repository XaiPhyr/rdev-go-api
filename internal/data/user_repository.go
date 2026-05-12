package data

import (
	"context"
	"time"

	"github.com/XaiPhyr/rdev-go-api/internal/dto"
	"github.com/uptrace/bun"
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CheckUserPermission(ctx context.Context, userID int64, roleName string) (bool, error) {
	var exists bool

	// var allPerms []string
	// err = s.db.NewRaw(`
	//     WITH RECURSIVE role_hierarchy AS (
	//         SELECT role_id FROM user_roles WHERE user_id = ?
	//         UNION
	//         SELECT gr.role_id FROM user_groups ug
	//         JOIN group_roles gr ON ug.group_id = gr.group_id
	//         WHERE ug.user_id = ?
	//         UNION
	//         SELECT r.parent_id FROM roles r
	//         INNER JOIN role_hierarchy rh ON r.id = rh.role_id
	//         WHERE r.parent_id IS NOT NULL
	//     )
	//     SELECT DISTINCT p.slug
	//     FROM role_hierarchy rh
	//     JOIN role_permissions rp ON rp.role_id = rh.role_id
	//     JOIN permissions p ON p.id = rp.permission_id
	// `, userID, userID).Scan(ctx, &allPerms)

	query := `
		WITH RECURSIVE role_hierarchy AS (
			SELECT role_id FROM user_roles WHERE user_id = ?
			UNION
			SELECT gr.role_id FROM user_groups ug
			JOIN group_roles gr ON ug.group_id = gr.group_id
			WHERE ug.user_id = ?
			UNION
			SELECT r.parent_id
			FROM roles r
			INNER JOIN role_hierarchy rh ON r.id = rh.role_id
			WHERE r.parent_id IS NOT NULL
		)
		SELECT EXISTS (
			SELECT 1 
			FROM role_hierarchy rh
			JOIN roles r ON r.id = rh.role_id
			WHERE r.name = 'super_admin'
			UNION ALL
			SELECT 1 
			FROM role_hierarchy rh
			JOIN role_permissions rp ON rp.role_id = rh.role_id
			JOIN permissions p ON p.id = rp.permission_id
			WHERE p.slug = ?
		)`

	err := r.db.NewRaw(query, userID, userID, roleName).Scan(ctx, &exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *UserRepository) GetUserByUsernameOrEmail(ctx context.Context, identifier string) (*User, error) {
	var user = new(User)

	err := r.db.NewSelect().
		Model(user).
		Where("username = ?", identifier).
		WhereOr("email = ?", identifier).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByUUID(ctx context.Context, uuid string) (*User, error) {
	var user = new(User)

	err := r.db.NewSelect().
		Model(user).
		Where("uuid = ?", uuid).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUsers(ctx context.Context, q dto.BaseFilters) ([]User, int, error) {
	var users []User

	count, err := r.db.NewSelect().
		Model(&users).
		Limit(q.PageSize).
		Offset(q.Page).
		Order(q.Sort).
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)

	return err
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *User) error {
	_, err := r.db.NewUpdate().
		Model(user).
		Column("first_name", "last_name", "email", "username").
		Set("updated_at = ?", time.Now()).
		WherePK().
		Exec(ctx)

	return err
}

func (r *UserRepository) DeleteUser(ctx context.Context, uuid string) error {
	_, err := r.db.NewDelete().
		Model((*User)(nil)).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}

func (r *UserRepository) UpdateUserStatus(ctx context.Context, uuid string) error {
	_, err := r.db.NewUpdate().
		Model((*User)(nil)).
		Set("status = CASE WHEN status = 'A' THEN 'I' ELSE 'A' END").
		Set("updated_at = ?", time.Now()).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}
