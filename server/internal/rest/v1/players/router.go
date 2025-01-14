package players

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
)

func Route(cfg *config.Config, db *db.Conn) func(r chi.Router) {
	hdl := newHandler(cfg, db)

	return func(r chi.Router) {
		r.Get("/", hdl.getPlayers)
		r.Get("/{playerId}", hdl.getPlayer)
		r.Put("/{playerId}", hdl.putPlayer)
	}
}
