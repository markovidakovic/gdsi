package players

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
	r.With(middleware.URLQueryPaginationParams).Get("/", a.hdl.getPlayers)
	r.With(middleware.URLPathUUIDParams("player_id")).Get("/{player_id}", a.hdl.getPlayer)
	r.With(middleware.URLPathUUIDParams("player_id")).With(middleware.RequirePermissionOrOwnership(permission.UpdatePlayer, a.hdl.store.checkPlayerOwnership, "account", "player_id")).Put("/{player_id}", a.hdl.updatePlayer)
}
