package seasons

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
// @Description Create a new season
// @Tags seasons
// @Accept json
// @Produce json
// @Param body body seasons.CreateSeasonRequestModel true "Request body"
// @Success 201 {object} seasons.SeasonModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons [post]
func (h *handler) postSeason(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusCreated, "create season")
}

// @Summary Get
// @Description Get seasons
// @Tags seasons
// @Produce json
// @Success 200 {array} seasons.SeasonModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons [get]
func (h *handler) getSeasons(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "get seasons")
}

// @Summary Get by id
// @Description Get season by id
// @Tags seasons
// @Produce json
// @Param seasonId path string true "Season id"
// @Success 200 {object} seasons.SeasonModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId} [get]
func (h *handler) getSeason(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "get season by id")
}

// @Summary Update
// @Description Update an existing season
// @Tags seasons
// @Accept json
// @Produce json
// @Param seasonId path string true "Season id"
// @Param body body seasons.UpdateSeasonRequestModel true "Request body"
// @Success 200 {object} seasons.SeasonModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId} [put]
func (h *handler) putSeason(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "update season")
}

// @Summary Delete
// @Description Delete an existing season
// @Tags seasons
// @Produce json
// @Param seasonId path string true "Season id"
// @Success 204 "No content"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId} [delete]
func (h *handler) deleteSeason(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusNoContent, "deleted season")
}
