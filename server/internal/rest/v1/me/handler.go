package me

import (
	"errors"
	"net/http"

	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
	custommiddleware "github.com/markovidakovic/gdsi/server/internal/middleware"
	"github.com/markovidakovic/gdsi/server/pkg/response"
)

type handler struct {
	service *service
}

// @Summary Get
// @Description Get my account and player profile data
// @Tags me
// @Produce json
// @Success 200 {object} me.MeModel "OK"
// @Failure 401 {object} response.BaseError "Unauthorized"
// @Failure 500 {object} response.BaseError "Internal server error"
// @Security BearerAuth
// @Router /v1/me [get]
func (h *handler) getMe(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(custommiddleware.AccountIdKey).(string)

	result, err := h.service.getMe(r.Context(), accountId)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			response.WriteError(w, response.NewNotFoundError("account not found"))
		default:
			response.WriteError(w, response.NewInternalError(err))
		}
		return
	}

	response.WriteSuccess(w, http.StatusOK, result)
}

// @Summary Update
// @Description Update my account and player profile data
// @Tags me
// @Accept json
// @Produce json
// @Param body body me.UpdateMeModel true "Request body"
// @Success 200 {object} me.MeModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.BaseError "Unauthorized"
// @Failure 404 {object} response.BaseError "Not found"
// @Failure 500 {object} response.BaseError "Internal server error"
// @Security BearerAuth
// @Router /v1/me [put]
func (h *handler) putMe(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "update me")
}

// @Summary Delete
// @Description Delete my account and player profile data
// @Tags me
// @Produce json
// @Success 204 "No content"
// @Failure 401 {object} response.BaseError "Unauthorized"
// @Failure 404 {object} response.BaseError "Not found"
// @Failure 500 {object} response.BaseError "Internal server error"
// @Security BearerAuth
// @Router /v1/me [delete]
func (h *handler) deleteMe(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusNoContent, "deleted me")
}

func newHandler(cfg *config.Config, db *db.Conn) *handler {
	h := &handler{}
	store := newStore(db)
	h.service = newService(cfg, store)
	return h
}
