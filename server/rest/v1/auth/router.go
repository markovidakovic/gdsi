package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
)

func Route(cfg *config.Config, db *db.Conn) func(r chi.Router) {
	hdl := newHandler(cfg, db)

	return func(r chi.Router) {
		r.Post("/signup", hdl.postSignup)
		r.Post("/tokens/access", hdl.postLogin)
		r.Post("/tokens/refresh", hdl.postRefreshToken)
		r.Post("/passwords/forgotten", hdl.postForgottenPassword)
		r.Put("/passwords/forgotten", hdl.putForgottenPassword)
	}
}
