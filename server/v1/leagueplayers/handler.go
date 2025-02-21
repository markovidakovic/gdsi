package leagueplayers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/response"
	"github.com/markovidakovic/gdsi/server/validation"
)

type handler struct {
	service *service
}

func newHandler(cfg *config.Config, db *db.Conn, validator *validation.Validator) *handler {
	h := &handler{}
	store := newStore(db)
	h.service = newService(cfg, store, validator)
	return h
}

// @Summary Get
// @Description Get league players
// @Tags players
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Success 200 {array} players.PlayerModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/players [get]
func (h *handler) getLeaguePlayers(w http.ResponseWriter, r *http.Request) {
	// call the service
	result, err := h.service.processGetLeaguePlayers(r.Context(), chi.URLParam(r, "seasonId"), chi.URLParam(r, "leagueId"))
	if err != nil {
		switch f := err.(type) {
		case *failure.ValidationFailure:
			response.WriteFailure(w, f)
			return
		case *failure.Failure:
			response.WriteFailure(w, f)
			return
		default:
			response.WriteFailure(w, failure.New("internal server error", err))
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
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/players/{playerId} [get]
func (h *handler) getLeaguePlayer(w http.ResponseWriter, r *http.Request) {
	// call the service
	result, err := h.service.processGetLeaguePlayer(r.Context(), chi.URLParam(r, "seasonId"), chi.URLParam(r, "leagueId"), chi.URLParam(r, "playerId"))
	if err != nil {
		switch f := err.(type) {
		case *failure.ValidationFailure:
			response.WriteFailure(w, f)
			return
		case *failure.Failure:
			response.WriteFailure(w, f)
			return
		default:
			response.WriteFailure(w, failure.New("internal server error", err))
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
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 404 {object} failure.Failure "Not found"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/players/{playerId}/assign [post]
func (h *handler) assignPlayerToLeague(w http.ResponseWriter, r *http.Request) {
	// call the service
	result, err := h.service.processAssignPlayerToLeague(r.Context(), chi.URLParam(r, "seasonId"), chi.URLParam(r, "leagueId"), chi.URLParam(r, "playerId"))
	if err != nil {
		switch f := err.(type) {
		case *failure.ValidationFailure:
			response.WriteFailure(w, f)
			return
		case *failure.Failure:
			response.WriteFailure(w, f)
			return
		default:
			response.WriteFailure(w, failure.New("internal server error", err))
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
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 404 {object} failure.Failure "Not found"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/players/{playerId}/assign [delete]
func (h *handler) unassignPlayerFromLeague(w http.ResponseWriter, r *http.Request) {
	// call the service
	result, err := h.service.processUnassignPlayerFromLeague(r.Context(), chi.URLParam(r, "seasonId"), chi.URLParam(r, "leagueId"), chi.URLParam(r, "playerId"))
	if err != nil {
		switch f := err.(type) {
		case *failure.ValidationFailure:
			response.WriteFailure(w, f)
			return
		case *failure.Failure:
			response.WriteFailure(w, f)
			return
		default:
			response.WriteFailure(w, failure.New("internal server error", err))
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
}
