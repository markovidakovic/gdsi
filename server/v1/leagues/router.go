package leagues

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
		r.With(middleware.RequirePermission(permission.CreateLeague)).Post("/", hdl.postLeague)
		r.Get("/", hdl.getLeagues)
		r.Get("/{leagueId}", hdl.getLeague)
		r.With(middleware.RequirePermission(permission.UpdateLeague)).Put("/{leagueId}", hdl.putLeague)
		r.With(middleware.RequirePermission(permission.DeleteLeague)).Delete("/{leagueId}", hdl.deleteLeague)
	}
}
