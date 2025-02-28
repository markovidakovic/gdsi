package pagination

import (
	"net/url"
	"strconv"
)

type Paginated[T any] struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	PageCount int `json:"page_count"`
	ItemCount int `json:"item_count"`
	Items     []T `json:"items"`
}

func NewPaginated[T any](page, perPage, itemCount int, items []T) *Paginated[T] {
	pageCount := calcPageCount(itemCount, perPage)

	return &Paginated[T]{
		Page:      page,
		PerPage:   perPage,
		PageCount: pageCount,
		ItemCount: itemCount,
		Items:     items,
	}
}

type QueryParams struct {
	Page    int
	PerPage int
	OrderBy string
}

func (qp *QueryParams) Populate(query url.Values) {
	pageStr := query.Get("page")
	perPageStr := query.Get("per_page")
	// orderBy := query.Get("order_by")

	if pageStr != "" || perPageStr != "" {
		if pageStr == "" {
			qp.Page = 1
		} else {
			qp.Page, _ = strconv.Atoi(pageStr) // validation already done in the middleware
		}
		if perPageStr == "" {
			qp.PerPage = 10
		} else {
			qp.PerPage, _ = strconv.Atoi(perPageStr) // validation already done in the middleware
		}
	}
}

func (qp *QueryParams) CalcLimitAndOffset(count int) (limit, offset int) {
	if qp.Page > 0 && qp.PerPage > 0 {
		pageCount := calcPageCount(count, qp.PerPage)
		if qp.Page > pageCount && pageCount > 0 {
			// if the wanted page is greater then the page count, set the page to the last page which is page count
			qp.Page = pageCount
		}

		limit = qp.PerPage
		offset = (qp.Page - 1) * qp.PerPage

	} else {
		limit = -1
		qp.Page = 1
		qp.PerPage = count
	}
	return
}

func calcPageCount(count, perPage int) int {
	pageCount := count / perPage
	if count%perPage > 0 {
		pageCount++
	}
	return pageCount
}
