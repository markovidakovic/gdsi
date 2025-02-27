package standings

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/response"
	"github.com/markovidakovic/gdsi/server/validation"
)

type handler struct {
	service *service
}

func newHandler(cfg *config.Config, db *db.Conn, validator *validation.Validator) *handler {
	h := &handler{}
	store := newStore(db)
	h.service = newService(cfg, store, validator)
	return h
}

// @Summary Get
// @Description Get standings
// @Tags standings
// @Produce json
// @Param season_id path string true "season id"
// @Param league_id path string true "league id"
// @Param page query int false "page"
// @Param per_page query int false "per page"
// @Param order_by query string false "order by"
// @Success 200 {array} standings.StandingModel "OK"
// @Failure 400 {object} failure.ValidationFailure "Bad request"
// @Failure 401 {object} failure.Failure "Unauthorized"
// @Failure 500 {object} failure.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{season_id}/leagues/{league_id}/standings [get]
func (h *handler) getStandings(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.processGetStandings(r.Context(), chi.URLParam(r, "season_id"), chi.URLParam(r, "league_id"))

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
