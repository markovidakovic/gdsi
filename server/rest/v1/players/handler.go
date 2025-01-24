package players

import (
	"net/http"

	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/response"
)

type handler struct {
	service *service
}

func newHandler(cfg *config.Config, db *db.Conn) *handler {
	h := &handler{}
	store := newStore(db)
	h.service = newService(cfg, store)
	return h
}

// @Summary Get
// @Description Get players
// @Tags players
// @Produce json
// @Success 200 {array} players.PlayerModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/players [get]
func (h *handler) getPlayers(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "get players")
}

// @Summary Get by id
// @Description Get player by id
// @Tags players
// @Produce json
// @Param playerId path string true "Player id"
// @Success 200 {object} matches.MatchModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/players/{playerId} [get]
func (h *handler) getPlayer(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "get player by id")
}

// @Summary Update
// @Description Update an existing player
// @Tags players
// @Accept json
// @Produce json
// @Param playerId path string true "Player id"
// @Param body body players.UpdatePlayerModel true "Request body"
// @Success 200 {object} players.PlayerModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/players/{playerId} [put]
func (h *handler) putPlayer(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "update player")
}
