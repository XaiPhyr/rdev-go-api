package models

import (
	"github.com/XaiPhyr/rdev-go-api/internal/shared/fields"
	"github.com/uptrace/bun"
)

type Category struct {
	bun.BaseModel `bun:"table:categories,alias:c"`
	fields.BaseFields

	Name     string `bun:"name" json:"name"`
	Slug     string `bun:"slug" json:"slug"`
	ParentID *int   `bun:"parent_id,default:null" json:"parent_id,omitempty"`
	Depth    string `bun:",scanonly" json:"depth,omitempty"`
	FullPath string `bun:",scanonly" json:"full_path,omitempty"`
}
