package leagues

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
	r.With(middleware.RequirePermission(permission.CreateLeague)).Post("/", a.hdl.createLeague)
	r.Get("/", a.hdl.getLeagues)
	r.Get("/{leagueId}", a.hdl.getLeague)
	r.With(middleware.RequirePermission(permission.UpdateLeague)).Put("/{leagueId}", a.hdl.updateLeague)
	r.With(middleware.RequirePermission(permission.DeleteLeague)).Delete("/{leagueId}", a.hdl.deleteLeague)
}
