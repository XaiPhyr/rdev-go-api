package data

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	BaseFields

	FirstName string `bun:"first_name" json:"first_name"`
	LastName  string `bun:"last_name" json:"last_name"`
	Email     string `bun:"email" json:"email"`
	Username  string `bun:"username" json:"username"`
	Password  string `bun:"password" json:"-"`
}
