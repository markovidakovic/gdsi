package matches

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/middleware"
)

func Route(cfg *config.Config, db *db.Conn) func(r chi.Router) {
	hdl := newHandler(cfg, db)

	return func(r chi.Router) {
		r.Post("/", hdl.postMatch)
		r.Get("/", hdl.getMatches)
		r.Get("/{matchId}", hdl.getMatch)
		r.With(middleware.RequireOwnership(hdl.store.checkMatchOwnership, "player", "matchId")).Put("/{matchId}", hdl.putMatch)
		r.With(middleware.RequireOwnership(hdl.store.checkMatchParticipation, "player", "matchId")).Post("/{matchId}/score", hdl.postMatchScore)
		r.With(middleware.RequireOwnership(hdl.store.checkMatchOwnership, "player", "matchId")).Delete("/{matchId}", hdl.deleteMatch)
	}
}
