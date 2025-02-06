package leagueplayers

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
		r.Get("/", hdl.getLeaguePlayers)
		r.Get("/{playerId}", hdl.getLeaguePlayer)
		r.With(middleware.RequirePermission(permission.UpdatePlayer)).Post("/{playerId}/assign", hdl.assignLeaguePlayer)
		r.With(middleware.RequirePermission(permission.UpdatePlayer)).Delete("/{playerId}/assign", hdl.removeLeaguePlayer)
	}
}
