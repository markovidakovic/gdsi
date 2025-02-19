package v1

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/middleware"
	"github.com/markovidakovic/gdsi/server/router"
	"github.com/markovidakovic/gdsi/server/v1/auth"
	"github.com/markovidakovic/gdsi/server/v1/courts"
	"github.com/markovidakovic/gdsi/server/v1/leagueplayers"
	"github.com/markovidakovic/gdsi/server/v1/leagues"
	"github.com/markovidakovic/gdsi/server/v1/matches"
	"github.com/markovidakovic/gdsi/server/v1/me"
	"github.com/markovidakovic/gdsi/server/v1/players"
	"github.com/markovidakovic/gdsi/server/v1/seasons"
	"github.com/markovidakovic/gdsi/server/v1/standings"
	"github.com/markovidakovic/gdsi/server/validation"
)

type api struct {
	cfg       *config.Config
	db        *db.Conn
	validator *validation.Validator
}

var _ router.Mounter = (*api)(nil)

func New(cfg *config.Config, db *db.Conn) *api {
	return &api{
		cfg:       cfg,
		db:        db,
		validator: validation.NewValidator(db),
	}
}

func (a *api) Mount(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Route("/auth", auth.New(a.cfg, a.db).Mount)
	})
	r.Group(func(r chi.Router) {
		// seek, verify and validate jwt
		r.Use(jwtauth.Verifier(a.cfg.JwtAuth))
		r.Use(jwtauth.Authenticator(a.cfg.JwtAuth))
		r.Use(middleware.AccountInfo)

		r.Route("/me", me.New(a.cfg, a.db).Mount)
		r.Route("/courts", courts.New(a.cfg, a.db).Mount)
		r.Route("/players", players.New(a.cfg, a.db).Mount)
		r.Route("/seasons", seasons.New(a.cfg, a.db).Mount)
		r.Route("/seasons/{seasonId}/leagues", leagues.New(a.cfg, a.db).Mount)
		r.Route("/seasons/{seasonId}/leagues/{leagueId}/players", leagueplayers.New(a.cfg, a.db).Mount)
		r.Route("/seasons/{seasonId}/leagues/{leagueId}/matches", matches.New(a.cfg, a.db, a.validator).Mount)
		r.Route("/seasons/{seasonId}/leagues/{leagueId}/standings", standings.New(a.cfg, a.db).Mount)
	})
	log.Println("v1 endpoints mounted")
}
