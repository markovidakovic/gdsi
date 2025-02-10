package courts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/middleware"
	"github.com/markovidakovic/gdsi/server/response"
)

type handler struct {
	store   *store
	service *service
}

func newHandler(cfg *config.Config, db *db.Conn) *handler {
	h := &handler{}
	h.store = newStore(db)
	h.service = newService(cfg, h.store)
	return h
}

// @Summary Create
// @Description Create a new court
// @Tags courts
// @Accept json
// @Produce json
// @Param body body courts.CreateCourtRequestModel true "Request body"
// @Success 200 {object} courts.CourtModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/courts [post]
func (h *handler) createCourt(w http.ResponseWriter, r *http.Request) {
	var input CreateCourtRequestModel
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	// validate
	if valErr := input.Validate(); valErr != nil {
		response.WriteFailure(w, response.NewValidationFailure("validation failed", valErr))
		return
	}

	// attach account id
	input.CreatorId = r.Context().Value(middleware.AccountIdCtxKey).(string)

	// store call
	result, err := h.store.insertCourt(r.Context(), input)
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
	// store call
	result, err := h.store.findCourts(r.Context())
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
	// store call
	result, err := h.store.findCourt(r.Context(), chi.URLParam(r, "courtId"))
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure("court not found"))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}
	response.WriteSuccess(w, http.StatusOK, result)
}

// @Summary Update
// @Description Update an existing court
// @Tags courts
// @Accept json
// @Produce json
// @Param courtId path string true "Court id"
// @Param body body courts.UpdateCourtRequestModel true "Request body"
// @Success 200 {object} courts.CourtModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/courts/{courtId} [put]
func (h *handler) updateCourt(w http.ResponseWriter, r *http.Request) {
	var input UpdateCourtRequestModel
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	// validate input
	if valErr := input.Validate(); valErr != nil {
		response.WriteFailure(w, response.NewValidationFailure("validation failed", valErr))
		return
	}

	result, err := h.store.updateCourt(r.Context(), chi.URLParam(r, "courtId"), input)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure("court not found"))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
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
	// call store
	err := h.store.deleteCourt(r.Context(), chi.URLParam(r, "courtId"))
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure("court not found"))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusNoContent, nil)
}
