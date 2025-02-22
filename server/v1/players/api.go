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
	r.Get("/", a.hdl.getPlayers)
	r.With(middleware.URLPathParamUUIDs("playerId")).Get("/{playerId}", a.hdl.getPlayer)
	r.With(middleware.URLPathParamUUIDs("playerId")).With(middleware.RequirePermissionOrOwnership(permission.UpdatePlayer, a.hdl.store.checkPlayerOwnership, "account", "playerId")).Put("/{playerId}", a.hdl.updatePlayer)
}
