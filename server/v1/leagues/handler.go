package leagues

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/middleware"
	"github.com/markovidakovic/gdsi/server/pagination"
	"github.com/markovidakovic/gdsi/server/params"
	"github.com/markovidakovic/gdsi/server/response"
	"github.com/markovidakovic/gdsi/server/validation"
)

type handler struct {
	service *service
	store   *store
}

func newHandler(cfg *config.Config, db *db.Conn, validator *validation.Validator) *handler {
	h := &handler{}
	h.store = newStore(db)
	h.service = newService(cfg, h.store, validator)
	return h
}

// @Summary Create
// @Description Create a new league
// @Tags leagues
// @Accept json
// @Produce json
// @Param season_id path string true "season id"
// @Param body body leagues.CreateLeagueRequestModel true "Request body"
// @Success 201 {object} leagues.LeagueModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id}/leagues [post]
func (h *handler) createLeague(w http.ResponseWriter, r *http.Request) {
	var model CreateLeagueRequestModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		response.WriteFailure(w, failure.New("invalid request body", fmt.Errorf("%w -> %v", failure.ErrBadRequest, err)))
		return
	}

	if valErr := model.Validate(); valErr != nil {
		response.WriteFailure(w, failure.NewValidation("validation failed", valErr))
		return
	}

	model.SeasonId = chi.URLParam(r, "season_id")
	model.CreatorId = r.Context().Value(middleware.AccountIdCtxKey).(string)

	result, err := h.service.processCreateLeague(r.Context(), model)
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
// @Description Get leagues
// @Tags leagues
// @Produce json
// @Param season_id path string true "season id"
// @Param page query int false "page"
// @Param per_page query int false "per page"
// @Param order_by query string false "order by"
// @Success 200 {array} leagues.LeagueModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id}/leagues [get]
func (h *handler) getLeagues(w http.ResponseWriter, r *http.Request) {
	query := params.NewQuery(r.URL.Query())

	leagues, count, err := h.service.processFindLeagues(r.Context(), chi.URLParam(r, "season_id"), query)
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

	result := pagination.NewPaginated(query.Page, query.PerPage, count, leagues)

	response.WriteSuccess(w, http.StatusOK, result)
}

// @Summary Get by id
// @Description Get league by id
// @Tags leagues
// @Produce json
// @Param season_id path string true "season id"
// @Param league_id path string true "league id"
// @Success 200 {object} leagues.LeagueModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 404 {object} failure.Failure "Not found"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id}/leagues/{league_id} [get]
func (h *handler) getLeague(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.processFindLeague(r.Context(), chi.URLParam(r, "season_id"), chi.URLParam(r, "league_id"))
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
// @Description Update an existing league
// @Tags leagues
// @Accept json
// @Produce json
// @Param season_id path string true "season id"
// @Param league_id path string true "league id"
// @Param body body leagues.UpdateLeagueRequestModel true "Request body"
// @Success 200 {object} leagues.LeagueModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 404 {object} failure.Failure "Not found"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id}/leagues/{league_id} [put]
func (h *handler) updateLeague(w http.ResponseWriter, r *http.Request) {
	var model UpdateLeagueRequestModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		response.WriteFailure(w, failure.New("invalid request body", fmt.Errorf("%w -> %v", failure.ErrBadRequest, err)))
		return
	}

	if valErr := model.Validate(); valErr != nil {
		response.WriteFailure(w, failure.NewValidation("validation vailed", valErr))
		return
	}

	model.SeasonId = chi.URLParam(r, "season_id")
	model.LeagueId = chi.URLParam(r, "league_id")

	result, err := h.service.processUpdateLeague(r.Context(), model)
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
// @Description Delete an existing league
// @Tags leagues
// @Produce json
// @Param season_id path string true "season id"
// @Param league_id path string true "league id"
// @Success 204 "No content"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 404 {object} failure.Failure "Not found"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id}/leagues/{league_id} [delete]
func (h *handler) deleteLeague(w http.ResponseWriter, r *http.Request) {
	err := h.service.processDeleteLeague(r.Context(), chi.URLParam(r, "season_id"), chi.URLParam(r, "league_id"))
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
