package response

type Pagination[T any] struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalItems   int `json:"total_items"`
	ItemsPerPage int `json:"items_per_page"`
	Items        []T `json:"items"`
}

func NewPagination[T any](currentPage, totalItems, itemsPerPage int, items []T) Pagination[T] {
	// totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
	totalPages := totalItems / itemsPerPage
	if totalItems%itemsPerPage > 0 {
		totalPages++
	}

	return Pagination[T]{
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
		TotalItems:   totalItems,
		ItemsPerPage: itemsPerPage,
		Items:        items,
	}
}
