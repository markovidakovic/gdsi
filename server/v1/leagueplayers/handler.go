package leagueplayers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
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
// @Description Get league players
// @Tags players
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Success 200 {array} players.PlayerModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/players [get]
func (h *handler) getLeaguePlayers(w http.ResponseWriter, r *http.Request) {
	// call the service
	result, err := h.service.processGetLeaguePlayers(r.Context(), chi.URLParam(r, "seasonId"), chi.URLParam(r, "leagueId"))
	if err != nil {
		switch {
		case errors.Is(err, response.ErrBadRequest):
			response.WriteFailure(w, response.NewBadRequestFailure(err.Error()))
			return
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure(err.Error()))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
}

// @Summary Get
// @Description Get league player
// @Tags players
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Param playerId path string true "Player id"
// @Success 200 {object} players.PlayerModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/players/{playerId} [get]
func (h *handler) getLeaguePlayer(w http.ResponseWriter, r *http.Request) {
	// call the service
	result, err := h.service.processGetLeaguePlayer(r.Context(), chi.URLParam(r, "seasonId"), chi.URLParam(r, "leagueId"), chi.URLParam(r, "playerId"))
	if err != nil {
		switch {
		case errors.Is(err, response.ErrBadRequest):
			response.WriteFailure(w, response.NewBadRequestFailure(err.Error()))
			return
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure(err.Error()))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
}

// @Summary Assign
// @Descriptions Assigns player to a league
// @Tags players
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Param playerId path string true "Player id"
// @Success 200 {object} players.PlayerModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/players/{playerId}/assign [post]
func (h *handler) assignPlayerToLeague(w http.ResponseWriter, r *http.Request) {
	// call the service
	result, err := h.service.processAssignPlayerToLeague(r.Context(), chi.URLParam(r, "seasonId"), chi.URLParam(r, "leagueId"), chi.URLParam(r, "playerId"))
	if err != nil {
		switch {
		case errors.Is(err, response.ErrBadRequest):
			response.WriteFailure(w, response.NewBadRequestFailure(err.Error()))
			return
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure(err.Error()))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
}

// @Summary Remove
// @Descriptions Removes player from a league
// @Tags players
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Param playerId path string true "Player id"
// @Success 200 {object} players.PlayerModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/players/{playerId}/assign [delete]
func (h *handler) unassignPlayerFromLeague(w http.ResponseWriter, r *http.Request) {
	// call the service
	result, err := h.service.processUnassignPlayerFromLeague(r.Context(), chi.URLParam(r, "seasonId"), chi.URLParam(r, "leagueId"), chi.URLParam(r, "playerId"))
	if err != nil {
		switch {
		case errors.Is(err, response.ErrBadRequest):
			response.WriteFailure(w, response.NewBadRequestFailure(err.Error()))
			return
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure(err.Error()))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
}
