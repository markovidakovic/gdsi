package courts

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
		r.With(middleware.RequirePermission(permission.CreateCourt)).Post("/", hdl.postCourt)
		r.Get("/", hdl.getCourt)
		r.Get("/{courtId}", hdl.getCourtById)
		r.With(middleware.RequirePermission(permission.UpdateCourt)).Put("/{courtId}", hdl.putCourt)
		r.With(middleware.RequirePermission(permission.DeleteCourt)).Delete("/{courtId}", hdl.deleteCourt)
	}
}
