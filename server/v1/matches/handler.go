package matches

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/middleware"
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
// @Description Create a new match
// @Tags matches
// @Accept json
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Param body body matches.CreateMatchRequestModel true "Request body"
// @Success 201 {object} matches.MatchModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/matches [post]
func (h *handler) createMatch(w http.ResponseWriter, r *http.Request) {
	// decode model
	var model CreateMatchRequestModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	// validate model
	if valErr := model.Validate(); valErr != nil {
		response.WriteFailure(w, response.NewValidationFailure("validation failed", valErr))
		return
	}

	ctx := r.Context()

	// set additional values in model
	model.SeasonId = chi.URLParam(r, "seasonId")
	model.LeagueId = chi.URLParam(r, "leagueId")
	model.PlayerOneId = ctx.Value(middleware.PlayerIdCtxKey).(string)

	// call the service
	result, err := h.service.processCreateMatch(ctx, model)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrBadRequest):
			response.WriteFailure(w, response.NewBadRequestFailure(err.Error()))
			return
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
// @Description Get matches
// @Tags matches
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Success 200 {array} matches.MatchModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/matches [get]
func (h *handler) getMatches(w http.ResponseWriter, r *http.Request) {
	// call the service
	result, err := h.service.processGetMatches(r.Context(), chi.URLParam(r, "seasonId"), chi.URLParam(r, "leagueId"))
	if err != nil {
		switch {
		case errors.Is(err, response.ErrBadRequest):
			response.WriteFailure(w, response.NewBadRequestFailure(err.Error()))
			return
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

// @Summary Get by id
// @Description Get match by id
// @Tags matches
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Param matchId path string true "Match id"
// @Success 200 {object} matches.MatchModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/matches/{matchId} [get]
func (h *handler) getMatch(w http.ResponseWriter, r *http.Request) {
	// call the service
	result, err := h.service.processGetMatch(r.Context(), chi.URLParam(r, "seasonId"), chi.URLParam(r, "leagueId"), chi.URLParam(r, "matchId"))
	if err != nil {
		switch {
		case errors.Is(err, response.ErrBadRequest):
			response.WriteFailure(w, response.NewBadRequestFailure(err.Error()))
			return
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
// @Description Update match
// @Tags matches
// @Accept json
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Param matchId path string true "Match id"
// @Param body body matches.UpdateMatchRequestModel true "Request body"
// @Success 200 {object} matches.MatchModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 409 {object} response.Failure "Conflict"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/matches/{matchId} [put]
func (h *handler) updateMatch(w http.ResponseWriter, r *http.Request) {
	// decode model
	var model UpdateMatchRequestModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	if valErr := model.Validate(); valErr != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("validation failed"))
		return
	}

	ctx := r.Context()

	model.SeasonId = chi.URLParam(r, "seasonId")
	model.LeagueId = chi.URLParam(r, "leagueId")
	model.MatchId = chi.URLParam(r, "matchId")
	model.PlayerOneId = ctx.Value(middleware.PlayerIdCtxKey).(string)

	// call the service
	result, err := h.service.processUpdateMatch(ctx, model)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrBadRequest):
			response.WriteFailure(w, response.NewBadRequestFailure(err.Error()))
			return
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

// @Summary Score
// @Description Submit a match score
// @Tags matches
// @Accept json
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Param matchId path string true "Match id"
// @Param body body matches.SubmitMatchScoreRequestModel true "Request body"
// @Success 200 {object} matches.MatchModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 409 {object} response.Failure "Conflict"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/matches/{matchId}/score [post]
func (h *handler) submitMatchScore(w http.ResponseWriter, r *http.Request) {
	// decode req body
	var model SubmitMatchScoreRequestModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	// validation
	if valErr := model.Validate(); valErr != nil {
		response.WriteFailure(w, response.NewValidationFailure("validation failed", valErr))
		return
	}

	// add params to model
	model.SeasonId = chi.URLParam(r, "seasonId")
	model.LeagueId = chi.URLParam(r, "leagueId")
	model.MatchId = chi.URLParam(r, "matchId")

	// call the service
	result, err := h.service.processSubmitMatchScore(r.Context(), model)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrBadRequest):
			response.WriteFailure(w, response.NewBadRequestFailure(err.Error()))
			return
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure(err.Error()))
			return
		case errors.Is(err, response.ErrConflict):
			response.WriteFailure(w, response.NewConflictFailure(err.Error()))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
}
