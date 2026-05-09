package dto

type BaseFilters struct {
	Page     int    `json:"page" query:"page"`
	PageSize int    `json:"page_size" query:"page_size"`
	Sort     string `json:"sort" query:"sort"`
	Status   string `json:"status" query:"status"`
	Search   string `json:"search" query:"search"`
}

type Query struct {
	Limit  int    `form:"limit,default=10"`
	Offset int    `form:"offset,default=0"`
	Search string `form:"search"`
	Sort   string `form:"sort"`
}

func (q Query) SanitizeQuery(allowedSortColumns []string) BaseFilters {
	pageSize := q.Limit
	if pageSize <= 0 {
		pageSize = 10
	}

	return BaseFilters{
		PageSize: pageSize,
		Page:     (q.Offset - 1) * pageSize,
		Sort:     q.validateSort(allowedSortColumns),
		Search:   q.Search,
	}
}

func (q Query) validateSort(allowedSortColumns []string) string {
	finalSort := "id ASC"
	if q.Sort != "" {
		for _, col := range allowedSortColumns {
			if q.Sort == col || q.Sort == col+" ASC" || q.Sort == col+" DESC" {
				finalSort = q.Sort
				break
			}
		}
	}

	return finalSort
}
