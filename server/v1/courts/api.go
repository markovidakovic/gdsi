package courts

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
	r.With(middleware.RequirePermission(permission.CreateCourt)).Post("/", a.hdl.createCourt)
	r.With(middleware.URLQueryPaginationParams).Get("/", a.hdl.getCourts)
	r.With(middleware.URLPathUUIDParams("court_id")).Get("/{court_id}", a.hdl.getCourt)
	r.With(middleware.URLPathUUIDParams("court_id")).With(middleware.RequirePermission(permission.UpdateCourt)).Put("/{court_id}", a.hdl.updateCourt)
	r.With(middleware.URLPathUUIDParams("court_id")).With(middleware.RequirePermission(permission.DeleteCourt)).Delete("/{court_id}", a.hdl.deleteCourt)
}
