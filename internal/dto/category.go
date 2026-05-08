package dto

type CategoryRequestUpdate struct {
	ParentID int    `json:"parent_id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
}
