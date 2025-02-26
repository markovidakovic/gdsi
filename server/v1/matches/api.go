package matches

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/middleware"
	"github.com/markovidakovic/gdsi/server/router"
	"github.com/markovidakovic/gdsi/server/validation"
)

type api struct {
	hdl *handler
}

var _ router.Mounter = (*api)(nil)

func New(cfg *config.Config, db *db.Conn, validator *validation.Validator) *api {
	return &api{
		hdl: newHandler(cfg, db, validator),
	}
}

func (a *api) Mount(r chi.Router) {
	r.With(middleware.URLPathParamUUIDs("season_id", "league_id")).Post("/", a.hdl.createMatch)
	r.With(middleware.URLPathParamUUIDs("season_id", "league_id")).Get("/", a.hdl.getMatches)
	r.With(middleware.URLPathParamUUIDs("season_id", "league_id", "match_id")).Get("/{match_id}", a.hdl.getMatch)
	r.With(middleware.URLPathParamUUIDs("season_id", "league_id", "match_id")).With(middleware.RequireOwnership(a.hdl.store.checkMatchOwnership, "player", "match_id")).Put("/{match_id}", a.hdl.updateMatch)
	r.With(middleware.URLPathParamUUIDs("season_id", "league_id", "match_id")).With(middleware.RequireOwnership(a.hdl.store.checkMatchParticipation, "player", "match_id")).Post("/{match_id}/score", a.hdl.submitMatchScore)
}
