package players

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/response"
)

type handler struct {
	service *service
	store   *store
}

func newHandler(cfg *config.Config, db *db.Conn) *handler {
	h := &handler{}
	h.store = newStore(db)
	h.service = newService(cfg, h.store)
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
	// call store
	result, err := h.store.findPlayers(r.Context())
	if err != nil {
		switch {
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
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
	// call store
	result, err := h.store.findPlayer(r.Context(), chi.URLParam(r, "playerId"))
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure("player not found"))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
}

// @Summary Update
// @Description Update an existing player
// @Tags players
// @Accept json
// @Produce json
// @Param playerId path string true "Player id"
// @Param body body players.UpdatePlayerRequestModel true "Request body"
// @Success 200 {object} players.PlayerModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/players/{playerId} [put]
func (h *handler) updatePlayer(w http.ResponseWriter, r *http.Request) {
	var input UpdatePlayerRequestModel
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	// validate input
	if valErr := input.Validate(); valErr != nil {
		response.WriteFailure(w, response.NewValidationFailure("validation failed", valErr))
		return
	}

	// call store
	result, err := h.store.updatePlayer(r.Context(), chi.URLParam(r, "playerId"), input)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure("player not found"))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
}
