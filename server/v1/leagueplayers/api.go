package leagueplayers

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/middleware"
	"github.com/markovidakovic/gdsi/server/permission"
	"github.com/markovidakovic/gdsi/server/router"
	"github.com/markovidakovic/gdsi/server/validation"
)

type API struct {
	hdl *handler
}

var _ router.Mounter = (*API)(nil)

func NewAPI(cfg *config.Config, db *db.Conn, validator *validation.Validator) (*API, error) {
	h, err := newHandler(cfg, db, validator)
	if err != nil {
		return nil, err
	}
	return &API{hdl: h}, nil
}

func (a *API) Mount(r chi.Router) {
	r.With(middleware.URLPathUUIDParams("season_id", "league_id")).With(middleware.URLQueryPaginationParams).Get("/", a.hdl.getLeaguePlayers)
	r.With(middleware.URLPathUUIDParams("season_id", "league_id", "player_id")).Get("/{player_id}", a.hdl.getLeaguePlayer)
	r.With(middleware.URLPathUUIDParams("season_id", "league_id", "player_id")).With(middleware.RequirePermission(permission.UpdatePlayer)).Post("/{player_id}/assign", a.hdl.assignPlayerToLeague)
	r.With(middleware.URLPathUUIDParams("season_id", "league_id", "player_id")).With(middleware.RequirePermission(permission.UpdatePlayer)).Delete("/{player_id}/assign", a.hdl.unassignPlayerFromLeague)
}
