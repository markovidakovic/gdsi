package leagues

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
)

func Route(cfg *config.Config, db *db.Conn) func(r chi.Router) {
	hdl := newHandler(cfg, db)

	return func(r chi.Router) {
		r.Post("/", hdl.postLeague)
		r.Get("/", hdl.getLeagues)
		r.Get("/{leagueId}", hdl.getLeague)
		r.Put("/{leagueId}", hdl.putLeague)
		r.Delete("/{leagueId}", hdl.deleteLeague)
	}
}
