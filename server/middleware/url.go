package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/response"
)

func URLPathUUIDParams(params ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			val := []failure.InvalidField{}
			for _, param := range params {
				id := chi.URLParam(r, param)

				err := uuid.Validate(id)
				if err != nil {
					val = append(val, failure.InvalidField{
						Field:    param,
						Message:  "invalid uuid format",
						Location: "path",
					})
				}
			}

			if len(val) > 0 {
				response.WriteFailure(w, failure.NewValidation("validation failed", val))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// URLQueryPaginationParams checks for the existance and validates the
// standard pagination query parameters: page, per_page, order_by
func URLQueryPaginationParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := []failure.InvalidField{}
		query := r.URL.Query()

		page, err := strconv.Atoi(query.Get("page"))
		if err != nil {
			val = append(val, failure.InvalidField{
				Field:    "page",
				Message:  "invalid value",
				Location: "query",
			})
		}
		perPage, err := strconv.Atoi(query.Get("per_page"))
		if err != nil {
			val = append(val, failure.InvalidField{
				Field:    "per_page",
				Message:  "invalid value",
				Location: "query",
			})
		}
		orderBy := query.Get("order_by")
		obSl := strings.Split(orderBy, " ")
		fmt.Printf("obSl: %v\n", obSl)

		fmt.Printf("page: %v\n", page)
		fmt.Printf("perPage: %v\n", perPage)
		fmt.Printf("orderBy: %v\n", orderBy)

		if len(val) > 0 {
			response.WriteFailure(w, failure.NewValidation("validation failed", val))
			return
		}

		next.ServeHTTP(w, r)
	})
}
