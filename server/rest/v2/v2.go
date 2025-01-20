package v2

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
)

func MountHandlers(cfg *config.Config, db db.Conn, rtr *chi.Mux) {
	rtr.Route("/v2", func(r chi.Router) {
	})

	log.Println("v2 endpoints mounted")
}
