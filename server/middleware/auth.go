package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/markovidakovic/gdsi/server/permission"
	"github.com/markovidakovic/gdsi/server/response"
)

var (
	AccountIdCtxKey   = &contextKey{"account-id"}
	AccountRoleCtxKey = &contextKey{"account-role"}
	PlayerIdCtxKey    = &contextKey{"player-id"}
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
		playerId, ok := claims["player_id"].(string)
		if !ok {
			response.WriteFailure(w, response.NewUnauthorizedFailure("account unauthorized"))
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, AccountIdCtxKey, accountId)
		ctx = context.WithValue(ctx, AccountRoleCtxKey, role)
		ctx = context.WithValue(ctx, PlayerIdCtxKey, playerId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequirePermission checks if the current account role has permission to access a specific resource
func RequirePermission(perm permission.Permission) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Context().Value(AccountRoleCtxKey).(string)

			// check permission
			if !permission.Has(role, perm) {
				response.WriteFailure(w, response.NewForbiddenFailure("insufficient permissions"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireOwnerrship checks if the current authenticated account is the resource owner
func RequireOwnership(oc OwnershipChecker, owner, resourceUrlPattern string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ownerId string

			if owner == "account" {
				ownerId = r.Context().Value(AccountIdCtxKey).(string)
			} else if owner == "player" {
				ownerId = r.Context().Value(PlayerIdCtxKey).(string)
			}

			resourceId := chi.URLParam(r, resourceUrlPattern)

			isOwner, err := oc(r.Context(), resourceId, ownerId)
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

// RequirePermissionOrOwnership checks if the current authenticated account is the resource owner or if
// it has rbac permission to access the resource
func RequirePermissionOrOwnership(perm permission.Permission, oc OwnershipChecker, owner, resourceUrlPattern string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ownerId string

			if owner == "account" {
				ownerId = r.Context().Value(AccountIdCtxKey).(string)
			} else if owner == "player" {
				ownerId = r.Context().Value(PlayerIdCtxKey).(string)
			}

			role := r.Context().Value(AccountRoleCtxKey).(string)

			// check permission
			if permission.Has(role, perm) {
				next.ServeHTTP(w, r)
				return
			}

			// get resource id of the url param with the provided pattern
			resourceId := chi.URLParam(r, resourceUrlPattern)

			// call the ownership checker func
			isOwner, err := oc(r.Context(), resourceId, ownerId)
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
