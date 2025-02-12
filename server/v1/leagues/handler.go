package leagues

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
func (h *handler) createLeague(w http.ResponseWriter, r *http.Request) {
	// decode model
	var model CreateLeagueRequestModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	// validate model
	if valErr := model.Validate(); valErr != nil {
		response.WriteFailure(w, response.NewValidationFailure("validation failed", valErr))
		return
	}

	model.SeasonId = chi.URLParam(r, "seasonId")
	model.CreatorId = r.Context().Value(middleware.AccountIdCtxKey).(string)

	// call the service
	result, err := h.service.processCreateLeague(r.Context(), model)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure(err.Error()))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusCreated, result)
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
	// call the service
	result, err := h.service.processFindLeagues(r.Context(), chi.URLParam(r, "seasonId"))
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
	// call the service
	result, err := h.service.processFindLeague(r.Context(), chi.URLParam(r, "seasonId"), chi.URLParam(r, "leagueId"))
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure(err.Error()))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
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
func (h *handler) updateLeague(w http.ResponseWriter, r *http.Request) {
	// decode model
	var model UpdateLeagueRequestModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	// validate model
	if valErr := model.Validate(); valErr != nil {
		response.WriteFailure(w, response.NewValidationFailure("validation failed", valErr))
		return
	}

	model.SeasonId = chi.URLParam(r, "seasonId")
	model.LeagueId = chi.URLParam(r, "leagueId")

	// call the service
	result, err := h.service.processUpdateLeague(r.Context(), model)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure(err.Error()))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
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
	// call the service
	err := h.service.processDeleteLeague(r.Context(), chi.URLParam(r, "seasonId"), chi.URLParam(r, "leagueId"))
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure(err.Error()))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusNoContent, nil)
}
