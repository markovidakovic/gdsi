package seasons

import (
	"net/http"

	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
	"github.com/markovidakovic/gdsi/server/pkg/response"
)

type handler struct {
}

// @Summary Create
// @Description Create a new season
// @Tags seasons
// @Accept json
// @Produce json
// @Param body body seasons.CreateSeasonModel true "Request body"
// @Success 201 {object} seasons.SeasonModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/seasons [post]
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusCreated, "create season")
}

// @Summary Get
// @Description Get seasons
// @Tags seasons
// @Produce json
// @Success 200 {array} seasons.SeasonModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/seasons [get]
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "get seasons")
}

// @Summary Get by id
// @Description Get season by id
// @Tags seasons
// @Produce json
// @Param seasonId path string true "Season id"
// @Success 200 {object} seasons.SeasonModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 404 {object} response.Error "Not found"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/seasons/{seasonId} [get]
func (h *handler) GetById(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "get season by id")
}

// @Summary Update
// @Description Update an existing season
// @Tags seasons
// @Accept json
// @Produce json
// @Param seasonId path string true "Season id"
// @Param body body seasons.UpdateSeasonModel true "Request body"
// @Success 200 {object} seasons.SeasonModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 404 {object} response.Error "Not found"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/seasons/{seasonId} [put]
func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "update season")
}

// @Summary Delete
// @Description Delete an existing season
// @Tags seasons
// @Produce json
// @Param seasonId path string true "Season id"
// @Success 204 "No content"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 404 {object} response.Error "Not found"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/seasons/{seasonId} [delete]
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusNoContent, "deleted season")
}

func NewHandler(cfg *config.Config, db *db.Conn) *handler {
	h := &handler{}
	return h
}
