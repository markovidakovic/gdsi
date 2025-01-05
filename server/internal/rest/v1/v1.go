package v1

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
	"github.com/markovidakovic/gdsi/server/internal/rest/v1/auth"
)

func MountHandlers(cfg *config.Config, db db.Conn, rtr *chi.Mux) {
	authHandler := auth.NewHandler(cfg, db)

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

		})
	})

	log.Println("v1 endpoints mounted")
}