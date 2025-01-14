package seasons

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
)

func Route(cfg *config.Config, db *db.Conn) func(r chi.Router) {
	hdl := newHandler(cfg, db)

	return func(r chi.Router) {
		r.Post("/", hdl.postSeason)
		r.Get("/", hdl.getSeasons)
		r.Get("/{seasonId}", hdl.getSeason)
		r.Put("/{seasonId}", hdl.putSeason)
		r.Delete("/{seasonId}", hdl.deleteSeason)
	}
}
