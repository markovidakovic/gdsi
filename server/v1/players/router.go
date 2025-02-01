package players

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/middleware"
	"github.com/markovidakovic/gdsi/server/permission"
)

func Route(cfg *config.Config, db *db.Conn) func(r chi.Router) {
	hdl := newHandler(cfg, db)

	return func(r chi.Router) {
		r.Get("/", hdl.getPlayers)
		r.Get("/{playerId}", hdl.getPlayer)
		r.With(middleware.RequirePermissionOrOwnership(permission.UpdatePlayer, hdl.store.checkPlayerOwnership, "playerId")).Put("/{playerId}", hdl.putPlayer)
	}
}
