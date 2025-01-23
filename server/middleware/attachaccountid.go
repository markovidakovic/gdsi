package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/markovidakovic/gdsi/server/response"
)

var (
	AccountIdCtxKey = &contextKey{"account-id"}
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

		ctx := context.WithValue(r.Context(), AccountIdCtxKey, accountId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
