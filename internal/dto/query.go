package dto

type Query struct {
	Limit  int    `form:"limit,default=10"`
	Offset int    `form:"offset,default=0"`
	Search string `form:"search"`
	Sort   string `form:"sort"`
}
