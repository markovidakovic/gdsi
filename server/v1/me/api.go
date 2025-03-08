package me

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/router"
)

type API struct {
	hdl *handler
}

var _ router.Mounter = (*API)(nil)

func NewAPI(cfg *config.Config, db *db.Conn) (*API, error) {
	h, err := newHandler(cfg, db)
	if err != nil {
		return nil, err
	}
	return &API{hdl: h}, nil
}

func (a *API) Mount(r chi.Router) {
	r.Get("/", a.hdl.getMe)
	r.Put("/", a.hdl.updateMe)
}
