package dto

type CategoryRequest struct {
	ParentID *int    `json:"parent_id"`
	Name     *string `json:"name"`
	Slug     *string `json:"slug"`
}

type CategoryPublicResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
