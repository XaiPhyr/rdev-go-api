package users

import (
	"github.com/XaiPhyr/rdev-go-api/internal/shared/fields"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	fields.BaseFields

	FirstName string `bun:"first_name" json:"first_name"`
	LastName  string `bun:"last_name" json:"last_name"`
	Email     string `bun:"email" json:"email"`
	Username  string `bun:"username" json:"username"`
	Password  string `bun:"password" json:"-"`
}

type UserRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	Username  *string `json:"username"`
}
