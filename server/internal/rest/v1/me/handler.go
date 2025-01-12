package me

import (
	"net/http"

	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
	"github.com/markovidakovic/gdsi/server/pkg/response"
)

type handler struct {
}

// @Summary Get
// @Description Get my account and player profile data
// @Tags me
// @Produce json
// @Success 200 {object} me.MeModel "OK"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/me [get]
func (h *handler) GetMe(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "me")
}

// @Summary Update
// @Description Update my account and player profile data
// @Tags me
// @Accept json
// @Produce json
// @Param body body me.UpdateMeModel true "Request body"
// @Success 200 {object} me.MeModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 404 {object} response.Error "Not found"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/me [put]
func (h *handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "update me")
}

// @Summary Delete
// @Description Delete my account and player profile data
// @Tags me
// @Produce json
// @Success 204 "No content"
// @Failure 401 {object} response.Error "Unauthorized"
// @Failure 404 {object} response.Error "Not found"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/me [delete]
func (h *handler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusNoContent, "deleted me")
}

func NewHandler(cfg *config.Config, db *db.Conn) *handler {
	h := &handler{}

	return h
}