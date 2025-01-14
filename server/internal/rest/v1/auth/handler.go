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
// @Accept json
// @Produce json
// @Param body body SignupRequestModel true "Request body"
// @Success 200 {object} auth.TokensResponseModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/auth/signup [post]
func (h *handler) postSignup(w http.ResponseWriter, r *http.Request) {
	var model SignupRequestModel

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		response.WriteError(w, response.Error{
			Status:  http.StatusBadRequest,
			Message: response.ErrBadRequest.Error(),
		})
		return
	}

	// Validate input
	validationErr := validateSignup(model)
	if validationErr != nil {
		response.WriteError(w, response.ValidationError{
			Error: response.Error{
				Status:  http.StatusBadRequest,
				Message: response.ErrBadRequest.Error(),
			},
			InvalidFields: validationErr,
		})
		return
	}

	// Call the service
	accessToken, refreshToken, err := h.service.signup(r.Context(), model)
	if err != nil {
		if errors.Is(err, response.ErrDuplicateRecord) {
			response.WriteError(w, response.Error{
				Status:  http.StatusBadRequest,
				Message: "account with email already exists",
			})
			return
		}

		response.WriteError(w, response.Error{
			Status:  http.StatusInternalServerError,
			Message: response.ErrInternal.Error(),
		})
		return
	}

	resp := TokensResponseModel{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.WriteSuccess(w, http.StatusCreated, resp)
}

// @Summary Login
// @Description Login and get a new access token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginRequestModel true "Request body"
// @Success 200 {object} auth.TokensResponseModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/auth/tokens/access [post]
func (h *handler) postLogin(w http.ResponseWriter, r *http.Request) {
	var model LoginRequestModel

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		response.WriteError(w, response.Error{
			Status:  http.StatusBadRequest,
			Message: response.ErrBadRequest.Error(),
		})
		return
	}

	// Validate input
	validationErr := validateLogin(model)
	if validationErr != nil {
		response.WriteError(w, response.ValidationError{
			Error: response.Error{
				Status:  http.StatusBadRequest,
				Message: response.ErrBadRequest.Error(),
			},
			InvalidFields: validationErr,
		})
		return
	}

	// Call the service
	accessToken, refreshToken, err := h.service.login(r.Context(), model)
	if err != nil {
		if errors.Is(err, response.ErrNotFound) {
			response.WriteError(w, response.Error{
				Status:  http.StatusBadRequest,
				Message: "invalid email or password",
			})
			return
		}

		response.WriteError(w, response.Error{
			Status:  http.StatusInternalServerError,
			Message: response.ErrInternal.Error(),
		})
		return
	}

	resp := TokensResponseModel{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.WriteSuccess(w, http.StatusOK, resp)
}

// @Summary Refresh token
// @Description Get a refreshed access token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} auth.TokensResponseModel "OK"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/auth/tokens/refresh [get]
func (h *handler) getRefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("refresh token endpoint\n"))
}

// @SUmmary Forgotten password
// @Description Get an email with a password reset link
// @Tags auth
// @Accept json
// @Produce json
// @Param body body auth.ForgottenPasswordRequestModel true "Request body"
// @Success 200 {object} auth.ForgottenPasswordResponseModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/auth/passwords/forgotten [post]
func (h *handler) postForgottenPassword(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "email sent")
}

// @SUmmary Forgotten password
// @Description Reset forgotten password
// @Tags auth
// @Accept json
// @Produce json
// @Param body body auth.ChangeForgottenPasswordRequestModel true "Request body"
// @Success 200 {object} auth.ChangeForgottenPasswordResponseModel "OK"
// @Failure 400 {object} response.ValidationError "Bad request"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /v1/auth/passwords/forgotten [put]
func (h *handler) putForgottenPassword(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "email sent")
}

func newHandler(cfg *config.Config, db *db.Conn) *handler {
	h := &handler{}
	store := newStore(db)
	h.service = newService(cfg, store)
	return h
}
