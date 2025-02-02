package leagues

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

// @Summary Create
// @Description Create a new league
// @Tags leagues
// @Accept json
// @Produce json
// @Param seasonId path string true "Season id"
// @Param body body leagues.CreateLeagueRequestModel true "Request body"
// @Success 201 {object} leagues.LeagueModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues [post]
func (h *handler) postLeague(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusCreated, "create league")
}

// @Summary Get
// @Description Get leagues
// @Tags leagues
// @Produce json
// @Param seasonId path string true "Season id"
// @Success 200 {array} leagues.LeagueModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues [get]
func (h *handler) getLeagues(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "get leagues")
}

// @Summary Get by id
// @Description Get league by id
// @Tags leagues
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Success 200 {object} leagues.LeagueModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId} [get]
func (h *handler) getLeague(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "get league by id")
}

// @Summary Update
// @Description Update an existing league
// @Tags leagues
// @Accept json
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Param body body leagues.UpdateLeagueRequestModel true "Request body"
// @Success 200 {object} leagues.LeagueModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId} [put]
func (h *handler) putLeague(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "update league")
}

// @Summary Delete
// @Description Delete an existing league
// @Tags leagues
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Success 204 "No content"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId} [delete]
func (h *handler) deleteLeague(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusNoContent, "deleted league")
}
