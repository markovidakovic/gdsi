package players

import (
	"net/http"

	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
	"github.com/markovidakovic/gdsi/server/pkg/response"
)

type handler struct {
}

// @Summary Get
// @Description Get players
// @Tags players
// @Produce json
// @Success 200 {array} players.PlayerModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/players [get]
func (h *handler) GetPlayers(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "get players")
}

// @Summary Get by id
// @Description Get player by id
// @Tags players
// @Produce json
// @Param playerId path string true "Player id"
// @Success 200 {object} matches.MatchModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 404 {object} response.Error "Not found"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/players/{playerId} [get]
func (h *handler) GetPlayerById(w http.ResponseWriter, r *http.Request) {
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
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 404 {object} response.Error "Not found"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/players/{playerId} [put]
func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "update player")
}

func NewHandler(cfg *config.Config, db *db.Conn) *handler {
	h := &handler{}
	return h
}
