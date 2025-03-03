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

func NewQuery(values url.Values) *Query {
	q := &Query{
		Additional: make(map[string]string),
	}

	pageStr := values.Get("page")
	perPageStr := values.Get("per_page")
	orderBy := values.Get("order_by")

	if pageStr != "" || perPageStr != "" {
		if pageStr == "" {
			q.Page = 1
		} else {
			q.Page, _ = strconv.Atoi(pageStr) // already validated in middleware
		}
		if perPageStr == "" {
			q.PerPage = 10
		} else {
			q.PerPage, _ = strconv.Atoi(perPageStr) // already validated in middleware
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

	// check for additional query params
	for k, v := range values {
		if k == "page" || k == "per_page" || k == "order_by" {
			continue
		}
		if len(v) > 0 {
			q.Additional[k] = v[0]
		}
	}

	return q
}

func (q *Query) GetBool(key string, defaultValue bool) bool {
	val, ok := q.Additional[key]
	if ok {
		if val == "true" {
			return true
		}
		if val == "false" {
			return false
		}
	}
	return defaultValue
}

type OrderBy struct {
	Field     string
	Direction string
}

func (ob *OrderBy) IsValid(allowed map[string]string) bool {
	if ob == nil {
		return true
	}
	valid := false
	for k := range allowed {
		if k == ob.Field {
			valid = true
			break
		}
	}
	if valid {
		valid = ob.Direction == "asc" || ob.Direction == "desc"
	}
	return valid
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
