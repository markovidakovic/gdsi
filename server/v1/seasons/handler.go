package seasons

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/response"
)

type handler struct {
	service *service
	store   *store
}

func newHandler(cfg *config.Config, db *db.Conn) *handler {
	h := &handler{}
	h.store = newStore(db)
	h.service = newService(cfg, h.store)
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
	var input CreateSeasonRequestModel
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	// validate input
	valErr := validatePostSeason(input)
	if valErr != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("validation failed"))
		return
	}

	// call the service
	result, err := h.service.processCreateSeason(r.Context(), input)
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
	// call the store
	result, err := h.store.findSeasons(r.Context())
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
	// call the store
	result, err := h.store.findSeason(r.Context(), chi.URLParam(r, "seasonId"))
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure("season not found"))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
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
	// decode req body
	var input UpdateSeasonRequestModel
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	// validate input
	if valErr := validatePutSeason(input); valErr != nil {
		response.WriteFailure(w, response.NewValidationFailure("validation failed", valErr))
		return
	}

	// call the store
	result, err := h.store.updateSeason(r.Context(), chi.URLParam(r, "seasonId"), input)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure("season not found"))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
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
	// call the store
	err := h.store.deleteSeason(r.Context(), chi.URLParam(r, "seasonId"))
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure("season not found"))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}
	response.WriteSuccess(w, http.StatusNoContent, "deleted")
}
