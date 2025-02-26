package matches

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/failure"
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
// @Param season_id path string true "Season id"
// @Param league_id path string true "League id"
// @Param body body matches.CreateMatchRequestModel true "Request body"
// @Success 201 {object} matches.MatchModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id}/leagues/{league_id}/matches [post]
func (h *handler) createMatch(w http.ResponseWriter, r *http.Request) {
	// decode model
	var model CreateMatchRequestModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		response.WriteFailure(w, failure.New("invalid request body", fmt.Errorf("%w -> %v", failure.ErrBadRequest, err)))
		return
	}

	// validate model
	if valErr := model.Validate(); valErr != nil {
		response.WriteFailure(w, failure.NewValidation("validation failed", valErr))
		return
	}

	ctx := r.Context()

	// set additional values in model
	model.SeasonId = chi.URLParam(r, "season_id")
	model.LeagueId = chi.URLParam(r, "league_id")
	model.PlayerOneId = ctx.Value(middleware.PlayerIdCtxKey).(string)

	// call the service
	result, err := h.service.processCreateMatch(ctx, model)
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
// @Description Get matches
// @Tags matches
// @Produce json
// @Param season_id path string true "Season id"
// @Param league_id path string true "League id"
// @Success 200 {array} matches.MatchModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id}/leagues/{league_id}/matches [get]
func (h *handler) getMatches(w http.ResponseWriter, r *http.Request) {
	// call the service
	result, err := h.service.processGetMatches(r.Context(), chi.URLParam(r, "season_id"), chi.URLParam(r, "league_id"))
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
// @Description Get match by id
// @Tags matches
// @Produce json
// @Param season_id path string true "Season id"
// @Param league_id path string true "League id"
// @Param match_id path string true "Match id"
// @Success 200 {object} matches.MatchModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 404 {object} failure.Failure "Not found"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id}/leagues/{league_id}/matches/{match_id} [get]
func (h *handler) getMatch(w http.ResponseWriter, r *http.Request) {
	// call the service
	result, err := h.service.processGetMatch(r.Context(), chi.URLParam(r, "season_id"), chi.URLParam(r, "league_id"), chi.URLParam(r, "match_id"))
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
// @Description Update match
// @Tags matches
// @Accept json
// @Produce json
// @Param season_id path string true "Season id"
// @Param league_id path string true "League id"
// @Param match_id path string true "Match id"
// @Param body body matches.UpdateMatchRequestModel true "Request body"
// @Success 200 {object} matches.MatchModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 404 {object} failure.Failure "Not found"
// @Failure 409 {object} failure.Failure "Conflict"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id}/leagues/{league_id}/matches/{match_id} [put]
func (h *handler) updateMatch(w http.ResponseWriter, r *http.Request) {
	// decode model
	var model UpdateMatchRequestModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		response.WriteFailure(w, failure.New("invalid request body", fmt.Errorf("%w -> %v", failure.ErrBadRequest, err)))
		return
	}

	if valErr := model.Validate(); valErr != nil {
		response.WriteFailure(w, failure.NewValidation("validation failed", valErr))
		return
	}

	ctx := r.Context()

	model.SeasonId = chi.URLParam(r, "season_id")
	model.LeagueId = chi.URLParam(r, "league_id")
	model.MatchId = chi.URLParam(r, "match_id")
	model.PlayerOneId = ctx.Value(middleware.PlayerIdCtxKey).(string)

	// call the service
	result, err := h.service.processUpdateMatch(ctx, model)
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

// @Summary Score
// @Description Submit a match score
// @Tags matches
// @Accept json
// @Produce json
// @Param season_id path string true "Season id"
// @Param league_id path string true "League id"
// @Param match_id path string true "Match id"
// @Param body body matches.SubmitMatchScoreRequestModel true "Request body"
// @Success 200 {object} matches.MatchModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 404 {object} failure.Failure "Not found"
// @Failure 409 {object} failure.Failure "Conflict"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id}/leagues/{league_id}/matches/{match_id}/score [post]
func (h *handler) submitMatchScore(w http.ResponseWriter, r *http.Request) {
	// decode req body
	var model SubmitMatchScoreRequestModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		response.WriteFailure(w, failure.New("invalid request body", fmt.Errorf("%w -> %v", failure.ErrBadRequest, err)))
		return
	}

	// validate model
	if valErr := model.Validate(); valErr != nil {
		response.WriteFailure(w, failure.NewValidation("validation failed", valErr))
		return
	}

	// add params to model
	model.SeasonId = chi.URLParam(r, "season_id")
	model.LeagueId = chi.URLParam(r, "league_id")
	model.MatchId = chi.URLParam(r, "match_id")

	// call the service
	result, err := h.service.processSubmitMatchScore(r.Context(), model)
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
