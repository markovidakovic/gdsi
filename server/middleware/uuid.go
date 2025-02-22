package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/response"
)

func URLPathParamUUIDs(params ...string) func(http.Handler) http.Handler {
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
