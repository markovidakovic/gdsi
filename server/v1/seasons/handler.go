package seasons

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/failure"
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
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons [post]
func (h *handler) createSeason(w http.ResponseWriter, r *http.Request) {
	var model CreateSeasonRequestModel
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		response.WriteFailure(w, failure.New("invalid request body", fmt.Errorf("%w -> %v", failure.ErrBadRequest, err)))
		return
	}

	if valErr := model.Validate(); valErr != nil {
		response.WriteFailure(w, failure.NewValidation("validation failed", valErr))
		return
	}

	result, err := h.service.processCreateSeason(r.Context(), model)
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
			return
		}
	}

	response.WriteSuccess(w, http.StatusCreated, result)
}

// @Summary Get
// @Description Get seasons
// @Tags seasons
// @Produce json
// @Param page query int false "page"
// @Param per_page query int false "per page"
// @Param order_by query string false "order by"
// @Success 200 {array} seasons.SeasonModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons [get]
func (h *handler) getSeasons(w http.ResponseWriter, r *http.Request) {
	result, err := h.store.findSeasons(r.Context())
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
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
}

// @Summary Get by id
// @Description Get season by id
// @Tags seasons
// @Produce json
// @Param season_id path string true "season id"
// @Success 200 {object} seasons.SeasonModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 404 {object} failure.Failure "Not found"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id} [get]
func (h *handler) getSeason(w http.ResponseWriter, r *http.Request) {
	result, err := h.store.findSeason(r.Context(), chi.URLParam(r, "season_id"))
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
// @Param season_id path string true "season id"
// @Param body body seasons.UpdateSeasonRequestModel true "Request body"
// @Success 200 {object} seasons.SeasonModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 404 {object} failure.Failure "Not found"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id} [put]
func (h *handler) updateSeason(w http.ResponseWriter, r *http.Request) {
	var model UpdateSeasonRequestModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		response.WriteFailure(w, failure.New("invalid request body", fmt.Errorf("%w -> %v", failure.ErrBadRequest, err)))
		return
	}

	if valErr := model.Validate(); valErr != nil {
		response.WriteFailure(w, failure.NewValidation("validation failed", valErr))
		return
	}

	result, err := h.store.updateSeason(r.Context(), nil, chi.URLParam(r, "season_id"), model)
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
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
}

// @Summary Delete
// @Description Delete an existing season
// @Tags seasons
// @Produce json
// @Param season_id path string true "season id"
// @Success 204 "No content"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 404 {object} failure.Failure "Not found"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id} [delete]
func (h *handler) deleteSeason(w http.ResponseWriter, r *http.Request) {
	err := h.store.deleteSeason(r.Context(), nil, chi.URLParam(r, "season_id"))
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
			return
		}
	}
	response.WriteSuccess(w, http.StatusNoContent, nil)
}
