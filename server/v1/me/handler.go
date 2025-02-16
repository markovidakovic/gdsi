package me

import (
	"encoding/json"
	"errors"
	"net/http"

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

// @Summary Get
// @Description Get my account and player profile data
// @Tags me
// @Produce json
// @Success 200 {object} me.MeModel "OK"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/me [get]
func (h *handler) getMe(w http.ResponseWriter, r *http.Request) {
	// get account id
	accountId := r.Context().Value(middleware.AccountIdCtxKey).(string)
	result, err := h.store.findMe(r.Context(), accountId)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure("account not found"))
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
}

// @Summary Update
// @Description Update my account and player profile data
// @Tags me
// @Accept json
// @Produce json
// @Param body body me.UpdateMeRequestModel true "Request body"
// @Success 200 {object} me.MeModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "Not found"
// @Failure 500 {object} response.Failure "Internal server error"
// @Security BearerAuth
// @Router /v1/me [put]
func (h *handler) updateMe(w http.ResponseWriter, r *http.Request) {
	var model UpdateMeRequestModel
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	// validation
	if valErr := model.Validate(); valErr != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("validation failed"))
		return
	}

	// get account id
	accountId := r.Context().Value(middleware.AccountIdCtxKey).(string)

	result, err := h.store.updateMe(r.Context(), nil, accountId, model)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewNotFoundFailure("account not found"))
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, result)
}
