package matches

import (
	"net/http"

	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/response"
)

type handler struct {
	service *service
}

// @Summary Create
// @Description Create a new match
// @Tags matches
// @Accept json
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Param body body matches.CreateMatchModel true "Request body"
// @Success 201 {object} matches.MatchModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.BaseError "Unauthorized"
// @Failure 500 {object} response.BaseError "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/matches [post]
func (h *handler) postMatch(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusCreated, "create match")
}

// @Summary Get
// @Description Get matches
// @Tags matches
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Success 200 {array} matches.MatchModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.BaseError "Unauthorized"
// @Failure 500 {object} response.BaseError "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/matches [get]
func (h *handler) getMatches(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "get matches")
}

// @Summary Get by id
// @Description Get match by id
// @Tags matches
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Param matchId path string true "Match id"
// @Success 200 {object} matches.MatchModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.BaseError "Unauthorized"
// @Failure 404 {object} response.BaseError "Not found"
// @Failure 500 {object} response.BaseError "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/matches/{matchId} [get]
func (h *handler) getMatch(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "get match by id")
}

// @Summary Update
// @Description Update an existing match
// @Tags matches
// @Accept json
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Param matchId path string true "Match id"
// @Param body body matches.UpdateMatchModel true "Request body"
// @Success 200 {object} matches.MatchModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.BaseError "Unauthorized"
// @Failure 404 {object} response.BaseError "Not found"
// @Failure 500 {object} response.BaseError "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/matches/{matchId} [put]
func (h *handler) putMatch(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "update match")
}

// @Summary Delete
// @Description Delete an existing league
// @Tags matches
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Param matchId path string true "Match id"
// @Success 204 "No content"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.BaseError "Unauthorized"
// @Failure 404 {object} response.BaseError "Not found"
// @Failure 500 {object} response.BaseError "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/matches/{matchId} [delete]
func (h *handler) deleteMatch(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusNoContent, "deleted match")
}

func newHandler(cfg *config.Config, db *db.Conn) *handler {
	h := &handler{}
	store := newStore(db)
	h.service = newService(cfg, store)
	return h
}
