package matches

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/middleware"
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
	r.With(middleware.URLPathUUIDParams("season_id", "league_id")).Post("/", a.hdl.createMatch)
	r.With(middleware.URLPathUUIDParams("season_id", "league_id")).With(middleware.URLQueryPaginationParams).Get("/", a.hdl.getMatches)
	r.With(middleware.URLPathUUIDParams("season_id", "league_id", "match_id")).Get("/{match_id}", a.hdl.getMatch)
	r.With(middleware.URLPathUUIDParams("season_id", "league_id", "match_id")).With(middleware.RequireOwnership(a.hdl.store.checkMatchOwnership, "player", "match_id")).Put("/{match_id}", a.hdl.updateMatch)
	r.With(middleware.URLPathUUIDParams("season_id", "league_id", "match_id")).With(middleware.RequireOwnership(a.hdl.store.checkMatchParticipation, "player", "match_id")).Post("/{match_id}/score", a.hdl.submitMatchScore)
}
