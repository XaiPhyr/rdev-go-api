package data

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	UUID      string    `bun:"uuid,notnull,unique,type:uuid,default:gen_random_uuid()" json:"uuid"`
	FirstName string    `bun:"first_name" json:"first_name"`
	LastName  string    `bun:"last_name" json:"last_name"`
	Email     string    `bun:"email" json:"email"`
	Username  string    `bun:"username" json:"username"`
	Password  string    `bun:"password" json:"-"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type UserFilters struct {
	Limit  int
	Offset int
	Search string
	Order  string // This will be the actual SQL "created_at DESC"
}
