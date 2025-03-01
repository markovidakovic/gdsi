package params

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/markovidakovic/gdsi/server/pagination"
)

type Query struct {
	Page       int
	PerPage    int
	OrderBy    *OrderBy
	Additional map[string]string
}

type OrderBy struct {
	Field     string
	Direction string
}

func NewQuery(values url.Values) *Query {
	q := &Query{}

	pageStr := values.Get("page")
	perPageStr := values.Get("per_page")
	orderBy := values.Get("order_by")

	if pageStr != "" || perPageStr != "" {
		if pageStr == "" {
			q.Page = 1
		} else {
			q.Page, _ = strconv.Atoi(pageStr)
		}
		if perPageStr == "" {
			q.PerPage = 10
		} else {
			q.PerPage, _ = strconv.Atoi(perPageStr)
		}
	}

	if orderBy != "" {
		obSl := strings.Split(orderBy, " ")
		q.OrderBy = &OrderBy{
			Field:     obSl[0],
			Direction: obSl[1],
		}
	} else {
		q.OrderBy = nil
	}

	return q
}

func (q *Query) CalcLimitAndOffset(count int) (limit, offset int) {
	if q.Page > 0 && q.PerPage > 0 {
		pageCount := pagination.CalcPageCount(count, q.PerPage)
		if q.Page > pageCount && pageCount > 0 {
			// if the wanted page is grater than page count, set the page to the last page which is page count
			q.Page = pageCount
		}

		limit = q.PerPage
		offset = (q.Page - 1) * q.PerPage
	} else {
		limit = -1
		q.Page = 1
		q.PerPage = count
	}
	return
}

func IsValidSortingField(allowed map[string]string, field string) bool {
	for k := range allowed {
		if k == field {
			return true
		}
	}
	return false
}
