package seasons

import (
	"net/http"

	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
	"github.com/markovidakovic/gdsi/server/pkg/response"
)

type handler struct {
}

// @Summary Create
// @Description Create a new season
// @Tags seasons
// @Accept application/json
// @Produce application/json
// @Param body body string true "Request body"
// @Success 200 {object} string "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/seasons [post]
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {

	response.WriteSuccess(w, http.StatusOK, "result")
}

func NewHandler(cfg *config.Config, db *db.Conn) *handler {
	h := &handler{}
	return h
}
