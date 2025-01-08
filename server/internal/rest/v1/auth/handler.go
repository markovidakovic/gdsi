package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
	"github.com/markovidakovic/gdsi/server/pkg/response"
)

type handler struct {
	service *service
}

// @Summary Signup
// @Description Signup a new account
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Param body body SignupRequestModel true "Request body"
// @Success 200 {object} auth.AccessTokenResponseModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/auth/signup [post]
func (h *handler) Signup(w http.ResponseWriter, r *http.Request) {
	var model SignupRequestModel

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		response.WriteError(w, response.Error{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
		})
		return
	}

	// Validate input
	validationErr := validateSignup(model)
	if validationErr != nil {
		response.WriteError(w, response.ValidationError{
			Error: response.Error{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
			},
			InvalidFields: validationErr,
		})
		return
	}

	// Call the service
	result, err := h.service.signupNewAccount(r.Context(), model)
	if err != nil {
		if errors.Is(err, response.ErrDuplicateRecord) {
			response.WriteError(w, response.Error{
				Status:  http.StatusBadRequest,
				Message: "Account with email already exists",
			})
			return
		}

		response.WriteError(w, response.Error{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	resp := AccessTokenResponseModel{
		AccessToken: result,
	}

	response.WriteSuccess(w, http.StatusCreated, resp)
}

// @Summary Login
// @Description Login and get a new access token
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Param body body LoginRequestModel true "Request body"
// @Success 200 {object} auth.AccessTokenResponseModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/auth/tokens/access [post]
func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var model LoginRequestModel

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		response.WriteError(w, response.Error{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
		})
		return
	}

	// Validate input
	validationErr := validateLogin(model)
	if validationErr != nil {
		response.WriteError(w, response.ValidationError{
			Error: response.Error{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
			},
			InvalidFields: validationErr,
		})
		return
	}

	// Call the service
	result, err := h.service.getAccessToken(r.Context(), model)
	if err != nil {
		if errors.Is(err, response.ErrNotFound) {
			response.WriteError(w, response.Error{
				Status:  http.StatusBadRequest,
				Message: "Invalid email or password",
			})
			return
		}

		response.WriteError(w, response.Error{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	resp := AccessTokenResponseModel{
		AccessToken: result,
	}

	response.WriteSuccess(w, http.StatusOK, resp)
}

// @Summary Refresh token
// @Description Get a refreshed access token
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} auth.AccessTokenResponseModel "OK"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/auth/tokens/refresh [get]
func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("refresh token endpoint\n"))
}

func NewHandler(cfg *config.Config, db *db.Conn) *handler {
	h := &handler{}

	store := newStore(db)
	service := newService(cfg, store)

	h.service = service

	return h
}
