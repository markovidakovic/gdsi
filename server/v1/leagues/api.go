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
	r.With(middleware.URLPathUUIDParams("season_id")).With(middleware.RequirePermission(permission.CreateLeague)).Post("/", a.hdl.createLeague)
	r.With(middleware.URLPathUUIDParams("season_id")).With(middleware.URLQueryPaginationParams).Get("/", a.hdl.getLeagues)
	r.With(middleware.URLPathUUIDParams("season_id", "league_id")).Get("/{league_id}", a.hdl.getLeague)
	r.With(middleware.URLPathUUIDParams("season_id", "league_id")).With(middleware.RequirePermission(permission.UpdateLeague)).Put("/{league_id}", a.hdl.updateLeague)
	r.With(middleware.URLPathUUIDParams("season_id", "league_id")).With(middleware.RequirePermission(permission.DeleteLeague)).Delete("/{league_id}", a.hdl.deleteLeague)
}
