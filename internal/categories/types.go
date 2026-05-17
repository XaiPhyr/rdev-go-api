package categories

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

type CategoryTree struct {
	bun.BaseModel `bun:"table:categories,alias:c"`

	ID       int64  `bun:"id,pk,autoincrement" json:"id"`
	Name     string `bun:"name" json:"name"`
	Slug     string `bun:"slug" json:"slug"`
	ParentID int    `bun:"parent_id" json:"parent_id"`
	Depth    string `bun:",scanonly" json:"depth"`
	FullPath string `bun:",scanonly" json:"full_path"`
}

type CategoryRequest struct {
	ParentID *int    `json:"parent_id"`
	Name     *string `json:"name"`
	Slug     *string `json:"slug"`
}

type CategoryResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
