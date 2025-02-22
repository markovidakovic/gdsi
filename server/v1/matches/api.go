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
	r.With(middleware.URLPathParamUUIDs("seasonId", "leagueId")).Post("/", a.hdl.createMatch)
	r.With(middleware.URLPathParamUUIDs("seasonId", "leagueId")).Get("/", a.hdl.getMatches)
	r.With(middleware.URLPathParamUUIDs("seasonId", "leagueId", "matchId")).Get("/{matchId}", a.hdl.getMatch)
	r.With(middleware.URLPathParamUUIDs("seasonId", "leagueId", "matchId")).With(middleware.RequireOwnership(a.hdl.store.checkMatchOwnership, "player", "matchId")).Put("/{matchId}", a.hdl.updateMatch)
	r.With(middleware.URLPathParamUUIDs("seasonId", "leagueId", "matchId")).With(middleware.RequireOwnership(a.hdl.store.checkMatchParticipation, "player", "matchId")).Post("/{matchId}/score", a.hdl.submitMatchScore)
}
