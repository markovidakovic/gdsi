package matches

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
)

func Route(cfg *config.Config, db *db.Conn) func(r chi.Router) {
	hdl := newHandler(cfg, db)

	return func(r chi.Router) {
		r.Post("/", hdl.postMatch)
		r.Get("/", hdl.getMatches)
		r.Get("/{matchId}", hdl.getMatch)
		r.Put("/{matchId}", hdl.putMatch)
		r.Delete("/{matchId}", hdl.deleteMatch)
	}
}
