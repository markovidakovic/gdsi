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

type v1 struct {
	cfg *config.Config
	// validator        *validation.Validator
	authAPI          *auth.API
	meAPI            *me.API
	courtsAPI        *courts.API
	playersAPI       *players.API
	seasonsAPI       *seasons.API
	leaguesAPI       *leagues.API
	leaguePlayersAPI *leagueplayers.API
	matchesAPI       *matches.API
	standingsAPI     *standings.API
}

var _ router.Mounter = (*v1)(nil)

func New(cfg *config.Config, db *db.Conn) (*v1, error) {
	validator := validation.NewValidator(db)

	authAPI, err := auth.NewAPI(cfg, db)
	if err != nil {
		return nil, err
	}
	meAPI, err := me.NewAPI(cfg, db)
	if err != nil {
		return nil, err
	}
	courtsAPI, err := courts.NewAPI(cfg, db)
	if err != nil {
		return nil, err
	}
	playersAPI, err := players.NewAPI(cfg, db)
	if err != nil {
		return nil, err
	}
	seasonsAPI, err := seasons.NewAPI(cfg, db)
	if err != nil {
		return nil, err
	}
	leaguesAPI, err := leagues.NewAPI(cfg, db, validator)
	if err != nil {
		return nil, err
	}
	leaguePlayersAPI, err := leagueplayers.NewAPI(cfg, db, validator)
	if err != nil {
		return nil, err
	}
	matchesAPI, err := matches.NewAPI(cfg, db, validator)
	if err != nil {
		return nil, err
	}
	standingsAPI, err := standings.NewAPI(cfg, db, validator)
	if err != nil {
		return nil, err
	}

	return &v1{
		cfg:              cfg,
		authAPI:          authAPI,
		meAPI:            meAPI,
		courtsAPI:        courtsAPI,
		playersAPI:       playersAPI,
		seasonsAPI:       seasonsAPI,
		leaguesAPI:       leaguesAPI,
		leaguePlayersAPI: leaguePlayersAPI,
		matchesAPI:       matchesAPI,
		standingsAPI:     standingsAPI,
	}, nil

}

func (a *v1) Mount(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Route("/auth", a.authAPI.Mount)
	})
	r.Group(func(r chi.Router) {
		// seek, verify and validate jwt
		r.Use(jwtauth.Verifier(a.cfg.JwtAuth))
		r.Use(jwtauth.Authenticator(a.cfg.JwtAuth))
		r.Use(middleware.AccountInfo)

		r.Route("/me", a.meAPI.Mount)
		r.Route("/courts", a.courtsAPI.Mount)
		r.Route("/players", a.playersAPI.Mount)
		r.Route("/seasons", a.seasonsAPI.Mount)
		r.Route("/seasons/{season_id}/leagues", a.leaguesAPI.Mount)
		r.Route("/seasons/{season_id}/leagues/{league_id}/players", a.leaguePlayersAPI.Mount)
		r.Route("/seasons/{season_id}/leagues/{league_id}/matches", a.matchesAPI.Mount)
		r.Route("/seasons/{season_id}/leagues/{league_id}/standings", a.standingsAPI.Mount)
	})
	log.Println("v1 endpoints mounted")
}
