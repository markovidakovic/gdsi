package me

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
)

func Route(cfg *config.Config, db *db.Conn) func(r chi.Router) {
	hdl := newHandler(cfg, db)
	return func(r chi.Router) {
		r.Get("/", hdl.getMe)
		r.Put("/", hdl.putMe)
		r.Delete("/", hdl.deleteMe)
	}
}
