package leagues

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/middleware"
	"github.com/markovidakovic/gdsi/server/permission"
	"github.com/markovidakovic/gdsi/server/router"
	"github.com/markovidakovic/gdsi/server/validation"
)

type api struct {
	hdl *handler
}

var _ router.Mounter = (*api)(nil)

func New(cfg *config.Config, db *db.Conn) *api {
	return &api{
		hdl: newHandler(cfg, db, validation.NewValidator(db)),
	}
}

func (a *api) Mount(r chi.Router) {
	r.With(middleware.URLPathParamUUIDs("season_id")).With(middleware.RequirePermission(permission.CreateLeague)).Post("/", a.hdl.createLeague)
	r.With(middleware.URLPathParamUUIDs("season_id")).Get("/", a.hdl.getLeagues)
	r.With(middleware.URLPathParamUUIDs("season_id", "league_id")).Get("/{league_id}", a.hdl.getLeague)
	r.With(middleware.URLPathParamUUIDs("season_id", "league_id")).With(middleware.RequirePermission(permission.UpdateLeague)).Put("/{league_id}", a.hdl.updateLeague)
	r.With(middleware.URLPathParamUUIDs("season_id", "league_id")).With(middleware.RequirePermission(permission.DeleteLeague)).Delete("/{league_id}", a.hdl.deleteLeague)
}
