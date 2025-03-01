package pagination

type Paginated[T any] struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	PageCount int `json:"page_count"`
	ItemCount int `json:"item_count"`
	Items     []T `json:"items"`
}

func NewPaginated[T any](page, perPage, itemCount int, items []T) *Paginated[T] {
	pageCount := CalcPageCount(itemCount, perPage)

	return &Paginated[T]{
		Page:      page,
		PerPage:   perPage,
		PageCount: pageCount,
		ItemCount: itemCount,
		Items:     items,
	}
}

func CalcPageCount(count, perPage int) int {
	pageCount := count / perPage
	if count%perPage > 0 {
		pageCount++
	}
	return pageCount
}
