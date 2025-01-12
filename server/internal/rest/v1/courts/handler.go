package courts

import (
	"net/http"

	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
	"github.com/markovidakovic/gdsi/server/pkg/response"
)

type handler struct{}

// @Summary Create
// @Description Create a new court
// @Tags courts
// @Accept json
// @Produce json
// @Param body body courts.CreateCourtModel true "Request body"
// @Success 200 {object} courts.CourtModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/courts [post]
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusCreated, "create court")
}

// @Summary Get
// @Description Get courts
// @Tags courts
// @Produce json
// @Success 200 {array} courts.CourtModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/courts [get]
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "get courts")
}

// @Summary Get by id
// @Description Get court by id
// @Tags courts
// @Produce json
// @Param courtId path string true "Court id"
// @Success 200 {object} courts.CourtModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 404 {object} response.Error "Not found"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/courts/{courtId} [get]
func (h *handler) GetById(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "get court by id")
}

// @Summary Update
// @Description Update an existing court
// @Tags courts
// @Accept json
// @Produce json
// @Param courtId path string true "Court id"
// @Param body body courts.UpdateCourtModel true "Request body"
// @Success 200 {object} courts.CourtModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 404 {object} response.Error "Not found"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/courts/{courtId} [put]
func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "update court")
}

// @Summary Delete
// @Description Delete an existing court
// @Tags courts
// @Produce json
// @Param courtId path string true "Court id"
// @Success 204 "No content"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 404 {object} response.Error "Not found"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/courts/{courtId} [delete]
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusNoContent, "deleted court")
}

func NewHandler(cfg *config.Config, db *db.Conn) *handler {
	h := &handler{}

	return h
}
