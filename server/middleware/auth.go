package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/markovidakovic/gdsi/server/permissions"
	"github.com/markovidakovic/gdsi/server/response"
)

var (
	AccountIdCtxKey   = &contextKey{"account-id"}
	AccountRoleCtxKey = &contextKey{"account-role"}
)

// OwnershipChecker type is used so we can provide the RequireOwnershipOrPermission middleware with
// a store function to check if the authenticated requestor created the resource
type OwnershipChecker = func(ctx context.Context, resourceId, accountId string) (bool, error)

// AttachAccountId sets the authenticated account id to the context for easier access in the handlers
func AccountInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())

		accountId, ok := claims["sub"].(string)
		if !ok {
			response.WriteFailure(w, response.NewUnauthorizedFailure("account unauthorized"))
			return
		}
		role, ok := claims["role"].(string)
		if !ok {
			response.WriteFailure(w, response.NewUnauthorizedFailure("account unauthorized"))
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, AccountIdCtxKey, accountId)
		ctx = context.WithValue(ctx, AccountRoleCtxKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequirePermission checks if the current account role has permission to access a specific resource
func RequirePermission(perm permissions.Permission) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Context().Value(AccountRoleCtxKey).(string)

			// check permission
			if !permissions.Has(role, perm) {
				response.WriteFailure(w, response.NewForbiddenFailure("insufficient permissions"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireOwnershipOrPermission checks if the current authenticated account is the resource owner or if
// it has rbac permission to access the resource
func RequireOwnershipOrPermission(perm permissions.Permission, oc OwnershipChecker, resourceUrlPattern string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accountId := r.Context().Value(AccountIdCtxKey).(string)
			role := r.Context().Value(AccountRoleCtxKey).(string)

			// check permission
			if permissions.Has(role, perm) {
				next.ServeHTTP(w, r)
				return
			}

			// get resource id of the url param with the provided pattern
			resourceId := chi.URLParam(r, resourceUrlPattern)
			fmt.Printf("resourceId: %v\n", resourceId)

			// call the ownership checker func
			isOwner, err := oc(r.Context(), resourceId, accountId)
			if err != nil {
				response.WriteFailure(w, response.NewInternalFailure(err))
				return
			}
			if !isOwner {
				response.WriteFailure(w, response.NewForbiddenFailure("not resource owner"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
