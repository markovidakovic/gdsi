package response

type Pagination struct {
	Page        int         `json:"page"`
	PageCurrent int         `json:"page_current"`
	PageCount   int         `json:"page_count"`
	Count       int         `json:"count"`
	Items       interface{} `json:"items"`
}

func NewPagination(page, pageCurrent, pageCount, count int, items interface{}) Pagination {
	return Pagination{
		Page:        page,
		PageCurrent: pageCurrent,
		PageCount:   pageCount,
		Count:       count,
		Items:       items,
	}
}
