package players

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/pagination"
	"github.com/markovidakovic/gdsi/server/params"
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

// @Summary Get
// @Description Get players
// @Tags players
// @Produce json
// @Param page query int false "page"
// @Param per_page query int false "per page"
// @Param order_by query string false "order by"
// @Success 200 {array} players.PlayerModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/players [get]
func (h *handler) getPlayers(w http.ResponseWriter, r *http.Request) {
	query := params.NewQuery(r.URL.Query())

	players, count, err := h.service.processGetPlayers(r.Context(), query)
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

	result := pagination.NewPaginated(query.Page, query.PerPage, count, players)

	response.WriteSuccess(w, http.StatusOK, result)
}

// @Summary Get by id
// @Description Get player by id
// @Tags players
// @Produce json
// @Param player_id path string true "player id"
// @Success 200 {object} matches.MatchModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 404 {object} failure.Failure "Not found"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/players/{player_id} [get]
func (h *handler) getPlayer(w http.ResponseWriter, r *http.Request) {
	result, err := h.store.findPlayer(r.Context(), chi.URLParam(r, "player_id"))
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
// @Description Update an existing player
// @Tags players
// @Accept json
// @Produce json
// @Param player_id path string true "player id"
// @Param body body players.UpdatePlayerRequestModel true "Request body"
// @Success 200 {object} players.PlayerModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 404 {object} failure.Failure "Not found"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/players/{player_id} [put]
func (h *handler) updatePlayer(w http.ResponseWriter, r *http.Request) {
	var model UpdatePlayerRequestModel
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		response.WriteFailure(w, failure.New("invalid request body", fmt.Errorf("%w -> %v", failure.ErrBadRequest, err)))
		return
	}

	if valErr := model.Validate(); valErr != nil {
		response.WriteFailure(w, failure.NewValidation("validation failed", valErr))
		return
	}

	result, err := h.store.updatePlayer(r.Context(), nil, chi.URLParam(r, "player_id"), model)
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
