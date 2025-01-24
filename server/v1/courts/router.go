package courts

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
)

func Route(cfg *config.Config, db *db.Conn) func(r chi.Router) {
	hdl := newHandler(cfg, db)
	return func(r chi.Router) {
		r.Post("/", hdl.postCourt)
		r.Get("/", hdl.getCourt)
		r.Get("/{courtId}", hdl.getCourtById)
		r.Put("/{courtId}", hdl.putCourt)
		r.Delete("/{courtId}", hdl.deleteCourt)
	}
}
