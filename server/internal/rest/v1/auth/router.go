package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
)

func Route(cfg *config.Config, db *db.Conn) func(r chi.Router) {
	hdl := newHandler(cfg, db)

	return func(r chi.Router) {
		r.Post("/signup", hdl.postSignup)
		r.Post("/tokens/access", hdl.postLogin)
		r.Get("/tokens/refresh", hdl.getRefreshToken)
		r.Post("/passwords/forgotten", hdl.postForgottenPassword)
		r.Put("/passwords/forgotten", hdl.putForgottenPassword)
	}
}
