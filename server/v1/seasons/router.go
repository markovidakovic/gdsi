package seasons

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
		r.With(middleware.RequirePermission(permission.CreateSeason)).Post("/", hdl.postSeason)
		r.Get("/", hdl.getSeasons)
		r.Get("/{seasonId}", hdl.getSeason)
		r.With(middleware.RequirePermission(permission.UpdateSeason)).Put("/{seasonId}", hdl.putSeason)
		r.With(middleware.RequirePermission(permission.DeleteSeason)).Delete("/{seasonId}", hdl.deleteSeason)
	}
}
