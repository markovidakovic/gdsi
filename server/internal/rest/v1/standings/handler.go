package standings

import (
	"net/http"

	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
	"github.com/markovidakovic/gdsi/server/pkg/response"
)

type handler struct {
	service *service
}

// @Summary Get
// @Description Get standings
// @Tags standings
// @Produce json
// @Param seasonId path string true "Season id"
// @Param leagueId path string true "League id"
// @Success 200 {array} standings.StandingModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 500 {object} response.Error "Internal server error"
// @Security BearerAuth
// @Router /v1/seasons/{seasonId}/leagues/{leagueId}/standings [get]
func (h *handler) getStandings(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "get standings")
}

func newHandler(cfg *config.Config, db *db.Conn) *handler {
	h := &handler{}
	store := newStore(db)
	h.service = newService(cfg, store)
	return h
}
