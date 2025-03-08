package seasons

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
	h, err := newHandler(cfg, db)
	if err != nil {
		return nil, err
	}
	return &API{hdl: h}, nil
}

func (a *API) Mount(r chi.Router) {
	r.With(middleware.RequirePermission(permission.CreateSeason)).Post("/", a.hdl.createSeason)
	r.With(middleware.URLQueryPaginationParams).Get("/", a.hdl.getSeasons)
	r.With(middleware.URLPathUUIDParams("season_id")).Get("/{season_id}", a.hdl.getSeason)
	r.With(middleware.URLPathUUIDParams("season_id")).With(middleware.RequirePermission(permission.UpdateSeason)).Put("/{season_id}", a.hdl.updateSeason)
	r.With(middleware.URLPathUUIDParams("season_id")).With(middleware.RequirePermission(permission.DeleteSeason)).Delete("/{season_id}", a.hdl.deleteSeason)
}
