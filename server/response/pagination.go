package response

type Pagination[T any] struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	PageCount int `json:"page_count"`
	ItemCount int `json:"item_count"`
	Items     []T `json:"items"`
}

type UrlQueryParams struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func NewPagination[T any](page, perPage, itemCount int, items []T) Pagination[T] {
	// pageCount := (itemCount + perPage - 1) / perPage
	pageCount := itemCount / perPage
	if itemCount%perPage > 0 {
		pageCount++
	}

	return Pagination[T]{
		Page:      page,
		PerPage:   perPage,
		PageCount: pageCount,
		ItemCount: itemCount,
		Items:     items,
	}
}
