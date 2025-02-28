package middleware

import (
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
		query := r.URL.Query()
		val := []failure.InvalidField{}

		if query.Get("page") != "" {
			page, err := strconv.Atoi(query.Get("page"))
			if err != nil {
				val = append(val, failure.InvalidField{
					Field:    "page",
					Message:  "invalid value",
					Location: "query",
				})

			}
			if page < 1 {
				val = append(val, failure.InvalidField{
					Field:    "page",
					Message:  "invalid value",
					Location: "query",
				})
			}
		}
		if query.Get("per_page") != "" {
			perPage, err := strconv.Atoi(query.Get("per_page"))
			if err != nil {
				val = append(val, failure.InvalidField{
					Field:    "per_page",
					Message:  "invalid value",
					Location: "query",
				})
			}
			if perPage < 1 {
				val = append(val, failure.InvalidField{
					Field:    "per_page",
					Message:  "invalid value",
					Location: "query",
				})
			}
		}
		if query.Get("order_by") != "" {
			obSl := strings.Split(query.Get("order_by"), " ")
			if len(obSl) != 2 || (obSl[1] != "asc" && obSl[1] != "desc") {
				val = append(val, failure.InvalidField{
					Field:    "order_by",
					Message:  "invalid value",
					Location: "query",
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
