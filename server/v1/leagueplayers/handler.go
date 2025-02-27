package leagueplayers

import (
	"fmt"
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
// @Param season_id path string true "Season id"
// @Param league_id path string true "League id"
// @Param page query int false "Page"
// @Param per_page query int false "Per page"
// @Param match_available query bool false "Match available"
// @Success 200 {array} players.PlayerModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id}/leagues/{league_id}/players [get]
func (h *handler) getLeaguePlayers(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	fmt.Printf("query: %v\n", query)

	result, err := h.service.processGetLeaguePlayers(r.Context(), chi.URLParam(r, "season_id"), chi.URLParam(r, "league_id"))
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
// @Param season_id path string true "Season id"
// @Param league_id path string true "League id"
// @Param player_id path string true "Player id"
// @Success 200 {object} players.PlayerModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id}/leagues/{league_id}/players/{player_id} [get]
func (h *handler) getLeaguePlayer(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.processGetLeaguePlayer(r.Context(), chi.URLParam(r, "season_id"), chi.URLParam(r, "league_id"), chi.URLParam(r, "player_id"))
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
// @Param season_id path string true "Season id"
// @Param league_id path string true "League id"
// @Param player_id path string true "Player id"
// @Success 200 {object} players.PlayerModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 404 {object} failure.Failure "Not found"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id}/leagues/{league_id}/players/{player_id}/assign [post]
func (h *handler) assignPlayerToLeague(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.processAssignPlayerToLeague(r.Context(), chi.URLParam(r, "season_id"), chi.URLParam(r, "league_id"), chi.URLParam(r, "player_id"))
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
// @Param season_id path string true "Season id"
// @Param league_id path string true "League id"
// @Param player_id path string true "Player id"
// @Success 200 {object} players.PlayerModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 404 {object} failure.Failure "Not found"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id}/leagues/{league_id}/players/{player_id}/assign [delete]
func (h *handler) unassignPlayerFromLeague(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.processUnassignPlayerFromLeague(r.Context(), chi.URLParam(r, "season_id"), chi.URLParam(r, "league_id"), chi.URLParam(r, "player_id"))
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
