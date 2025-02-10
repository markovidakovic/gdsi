package leagueplayers

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
	r.Get("/", a.hdl.getLeaguePlayers)
	r.Get("/{playerId}", a.hdl.getLeaguePlayer)
	r.With(middleware.RequirePermission(permission.UpdatePlayer)).Post("/{playerId}/assign", a.hdl.assignLeaguePlayer)
	r.With(middleware.RequirePermission(permission.UpdatePlayer)).Delete("/{playerId}/assign", a.hdl.removeLeaguePlayer)
}
