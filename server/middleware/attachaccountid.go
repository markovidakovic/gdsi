package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/markovidakovic/gdsi/server/response"
)

type contextKey string

const (
	AccountIdKey contextKey = "account_id"
)

// AttachAccountId sets the authenticated account id to the context for easier access in the handlers
func AttachAccountId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())

		accountId, ok := claims["sub"].(string)
		if !ok {
			response.WriteFailure(w, response.NewUnauthorizedFailure("account unauthorized"))
			return
		}

		// add account_id to the context
		ctx := context.WithValue(r.Context(), AccountIdKey, accountId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
