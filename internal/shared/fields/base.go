package fields

import "time"

type BaseFields struct {
	ID        int64      `bun:"id,pk,autoincrement" json:"id"`
	Status    string     `bun:"status,default:'A'" json:"status"`
	UUID      string     `bun:"uuid,notnull,unique,type:uuid,default:gen_random_uuid()" json:"uuid"`
	CreatedAt time.Time  `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"deleted_at"`
}
