package courts

import (
	"encoding/json"
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
// @Description Create a new court
// @Tags courts
// @Accept json
// @Produce json
// @Param body body courts.CreateCourtModel true "Request body"
// @Success 200 {object} courts.CourtModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/courts [post]
func (h *handler) postCourt(w http.ResponseWriter, r *http.Request) {
	var input CreateCourtModel
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	// validate
	valErr := validatePostCourt(input)
	if valErr != nil {
		response.WriteFailure(w, response.NewValidationFailure("validation failed", valErr))
		return
	}

	// service call
	result, err := h.service.createCourt(r.Context(), input)
	if err != nil {
		switch {
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusCreated, result)
}

// @Summary Get
// @Description Get courts
// @Tags courts
// @Produce json
// @Success 200 {array} courts.CourtModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/courts [get]
func (h *handler) getCourts(w http.ResponseWriter, r *http.Request) {
	// service call
	result, err := h.service.getCourts(r.Context())
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
// @Description Get court by id
// @Tags courts
// @Produce json
// @Param courtId path string true "Court id"
// @Success 200 {object} courts.CourtModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/courts/{courtId} [get]
func (h *handler) getCourt(w http.ResponseWriter, r *http.Request) {
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
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/courts/{courtId} [put]
func (h *handler) putCourt(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "update court")
}

// @Summary Delete
// @Description Delete an existing court
// @Tags courts
// @Produce json
// @Param courtId path string true "Court id"
// @Success 204 "No content"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/courts/{courtId} [delete]
func (h *handler) deleteCourt(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusNoContent, "deleted court")
}
