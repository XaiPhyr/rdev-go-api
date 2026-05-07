package dto

type UserRequestUpdate struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type UserQuery struct {
	Limit  int    `form:"limit,default=10"`
	Offset int    `form:"offset,default=0"`
	Search string `form:"search"`
	Sort   string `form:"sort"`
}
