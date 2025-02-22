package seasons

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/middleware"
	"github.com/markovidakovic/gdsi/server/permission"
	"github.com/markovidakovic/gdsi/server/router"
)

type api struct {
	hdl *handler
}

var _ router.Mounter = (*api)(nil)

func New(cfg *config.Config, db *db.Conn) *api {
	return &api{
		hdl: newHandler(cfg, db),
	}
}

func (a *api) Mount(r chi.Router) {
	r.With(middleware.RequirePermission(permission.CreateSeason)).Post("/", a.hdl.createSeason)
	r.Get("/", a.hdl.getSeasons)
	r.With(middleware.URLPathParamUUIDs("seasonId")).Get("/{seasonId}", a.hdl.getSeason)
	r.With(middleware.URLPathParamUUIDs("seasonId")).With(middleware.RequirePermission(permission.UpdateSeason)).Put("/{seasonId}", a.hdl.updateSeason)
	r.With(middleware.URLPathParamUUIDs("seasonId")).With(middleware.RequirePermission(permission.DeleteSeason)).Delete("/{seasonId}", a.hdl.deleteSeason)
}
