package v1

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
	"github.com/markovidakovic/gdsi/server/internal/rest/v1/auth"
	"github.com/markovidakovic/gdsi/server/internal/rest/v1/courts"
	"github.com/markovidakovic/gdsi/server/internal/rest/v1/leagues"
	"github.com/markovidakovic/gdsi/server/internal/rest/v1/matches"
	"github.com/markovidakovic/gdsi/server/internal/rest/v1/me"
	"github.com/markovidakovic/gdsi/server/internal/rest/v1/players"
	"github.com/markovidakovic/gdsi/server/internal/rest/v1/seasons"
	"github.com/markovidakovic/gdsi/server/internal/rest/v1/standings"
	"github.com/markovidakovic/gdsi/server/pkg/middleware"
)

func MountHandlers(cfg *config.Config, db *db.Conn) func(r chi.Router) {
	return func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Route("/auth", auth.Route(cfg, db))
		})
		r.Group(func(r chi.Router) {
			// seek, verify and validate jwt
			r.Use(jwtauth.Verifier(cfg.JwtAuth))
			r.Use(jwtauth.Authenticator(cfg.JwtAuth))
			r.Use(middleware.AttachAccountId)

			r.Route("/courts", courts.Route(cfg, db))
			r.Route("/me", me.Route(cfg, db))
			r.Route("/players", players.Route(cfg, db))
			r.Route("/seasons", seasons.Route(cfg, db))
			r.Route("/seasons/{seasonId}/leagues", leagues.Route(cfg, db))
			r.Route("/seasons/{seasonId}/leagues/{leagueId}/matches", matches.Route(cfg, db))
			r.Route("/seasons/{seasonId}/leagues/{leagueId}/standings", standings.Route(cfg, db))
		})

		log.Println("v1 endpoints mounted")
	}
}
