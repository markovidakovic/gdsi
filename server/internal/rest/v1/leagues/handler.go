package leagues

import (
	"net/http"

	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
	"github.com/markovidakovic/gdsi/server/pkg/response"
)

type handler struct {
}

// @Summary Create
// @Description Create a new league
// @Tags leagues
// @Accept json
// @Produce json
// @Param seasonId path string true "Season id"
// @Param body body leagues.CreateLeagueModel true "Request body"
// @Success 201 {object} leagues.LeagueModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/seasons/{seasonId}/leagues [post]
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusCreated, "create league")
}

// @Summary Get
// @Description Get leagues
// @Tags leagues
// @Produce json
// @Param seasonId path string true "Season id"
// @Success 200 {array} leagues.LeagueModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/seasons/{seasonId}/leagues [get]
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "get leagues")
}

// @Summary Get by id
// @Description Get league by id
// @Tags leagues
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Success 200 {object} leagues.LeagueModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 404 {object} response.Error "Not found"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/seasons/{seasonId}/leagues/{leagueId} [get]
func (h *handler) GetById(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "get league by id")
}

// @Summary Update
// @Description Update an existing league
// @Tags leagues
// @Accept json
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Param body body leagues.UpdateLeagueModel true "Request body"
// @Success 200 {object} leagues.LeagueModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 404 {object} response.Error "Not found"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/seasons/{seasonId}/leagues/{leagueId} [put]
func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "update league")
}

// @Summary Delete
// @Description Delete an existing league
// @Tags leagues
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Success 204 "No content"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 404 {object} response.Error "Not found"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/seasons/{seasonId}/leagues/{leagueId} [delete]
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusNoContent, "deleted league")
}

func NewHandler(cfg *config.Config, db *db.Conn) *handler {
	h := &handler{}
	return h
}
