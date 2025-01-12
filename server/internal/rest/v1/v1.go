package v1

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
	"github.com/markovidakovic/gdsi/server/internal/rest/v1/auth"
	"github.com/markovidakovic/gdsi/server/internal/rest/v1/courts"
	"github.com/markovidakovic/gdsi/server/internal/rest/v1/seasons"
)

func MountHandlers(cfg *config.Config, db *db.Conn, rtr *chi.Mux) {
	authHandler := auth.NewHandler(cfg, db)
	courtsHandler := courts.NewHandler(cfg, db)
	seasonsHandler := seasons.NewHandler(cfg, db)

	rtr.Route("/v1", func(r chi.Router) {
		// Public endpoints
		r.Group(func(r chi.Router) {
			r.Route("/auth", func(r chi.Router) {
				r.Post("/signup", authHandler.Signup)
				r.Post("/tokens/access", authHandler.Login)
			})
		})

		// Private endpoints
		r.Group(func(r chi.Router) {
			//  Seek, verify and validate JWT tokens
			r.Use(jwtauth.Verifier(cfg.JwtAuth))
			r.Use(jwtauth.Authenticator(cfg.JwtAuth))

			r.Route("/courts", func(r chi.Router) {
				r.Get("/", courtsHandler.Get)
			})

			r.Route("/seasons", func(r chi.Router) {
				r.Post("/", seasonsHandler.Create)
			})

		})
	})

	log.Println("v1 endpoints mounted")
}
