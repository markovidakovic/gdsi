package players

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/middleware"
	"github.com/markovidakovic/gdsi/server/permission"
	"github.com/markovidakovic/gdsi/server/router"
)

type API struct {
	hdl *handler
}

var _ router.Mounter = (*API)(nil)

func NewAPI(cfg *config.Config, db *db.Conn) (*API, error) {
	hdl, err := newHandler(cfg, db)
	if err != nil {
		return nil, err
	}
	return &API{hdl}, nil
}

func (a *API) Mount(r chi.Router) {
	r.With(middleware.URLQueryPaginationParams).Get("/", a.hdl.getPlayers)
	r.With(middleware.URLPathUUIDParams("player_id")).Get("/{player_id}", a.hdl.getPlayer)
	r.With(middleware.URLPathUUIDParams("player_id")).With(middleware.RequirePermissionOrOwnership(permission.UpdatePlayer, a.hdl.store.checkPlayerOwnership, "account", "player_id")).Put("/{player_id}", a.hdl.updatePlayer)
}
