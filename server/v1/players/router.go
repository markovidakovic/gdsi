package players

import (
	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/middleware"
	"github.com/markovidakovic/gdsi/server/permission"
)

type Router struct {
	hdl *handler
}

func NewRouter(cfg *config.Config, db *db.Conn) *Router {
	return &Router{
		hdl: newHandler(cfg, db),
	}
}

func (r *Router) RouteGlobal() func(cr chi.Router) {
	return func(cr chi.Router) {
		cr.Get("/", r.hdl.getPlayers)
		cr.Get("/{playerId}", r.hdl.getPlayer)
		cr.With(middleware.RequirePermissionOrOwnership(permission.UpdatePlayer, r.hdl.store.checkPlayerOwnership, "playerId")).Put("/{playerId}", r.hdl.putPlayer)
	}
}

func (r *Router) RouteLeague() func(cr chi.Router) {
	return func(cr chi.Router) {
		cr.Get("/", r.hdl.getLeaguePlayers)
		cr.Get("/{playerId}", r.hdl.getLeaguePlayer)
		cr.With(middleware.RequirePermission(permission.UpdatePlayer)).Put("/", r.hdl.updateLeaguePlayer) // updates player.current_league_id (for now)
	}
}
