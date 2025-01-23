package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/response"
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
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 500 {object} response.Failure "Internal server error"
// @Router /v1/auth/signup [post]
func (h *handler) postSignup(w http.ResponseWriter, r *http.Request) {
	var input SignupRequestModel

	// decode request body
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	// validation
	valErr := validateSignup(input)
	if valErr != nil {
		response.WriteFailure(w, response.NewValidationFailure("validation failed", valErr))
		return
	}

	// call the service
	accessToken, refreshToken, err := h.service.signup(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrDuplicateRecord):
			response.WriteFailure(w, response.NewBadRequestFailure("account with email already exists"))
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
		}
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
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 500 {object} response.Failure "Internal server error"
// @Router /v1/auth/tokens/access [post]
func (h *handler) postLogin(w http.ResponseWriter, r *http.Request) {
	var input LoginRequestModel

	// decode request body
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	// validate input
	valErr := validateLogin(input)
	if valErr != nil {
		response.WriteFailure(w, response.NewValidationFailure("validation failed", valErr))
		return
	}

	// for now not used, just to see in the future what to maybe store
	ra := r.RemoteAddr
	host, port, err := net.SplitHostPort(ra)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("ip: %v\n", ra)
	fmt.Printf("host: %v\n", host)
	fmt.Printf("port: %v\n", port)
	fmt.Printf("r.UserAgent(): %v\n", r.UserAgent())

	// Call the service
	accessToken, refreshToken, err := h.service.login(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrDuplicateRecord):
			response.WriteFailure(w, response.NewBadRequestFailure("invalid email or password"))
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
		}
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
// @Param body body RefreshTokenRequestModel true "Request body"
// @Success 200 {object} auth.TokensResponseModel "OK"
// @Failure 500 {object} response.Failure "Internal server error"
// @Router /v1/auth/tokens/refresh [post]
func (h *handler) postRefreshToken(w http.ResponseWriter, r *http.Request) {
	var input RefreshTokenRequestModel

	// decode request body
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	// validate input
	valErr := validateRefreshToken(input)
	if valErr != nil {
		response.WriteFailure(w, response.NewValidationFailure("validation failed", valErr))
		return
	}

	access, refresh, err := h.service.refreshTokens(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			response.WriteFailure(w, response.NewUnauthorizedFailure("invalid refresh token"))
			return
		case errors.Is(err, response.ErrUnauthorized):
			response.WriteFailure(w, response.NewUnauthorizedFailure(err.Error()))
			return
		default:
			response.WriteFailure(w, response.NewInternalFailure(err))
		}
		return
	}

	resp := TokensResponseModel{
		AccessToken:  access,
		RefreshToken: refresh,
	}

	response.WriteSuccess(w, http.StatusOK, resp)
}

// @Summary Forgotten password
// @Description Get an email with a password reset link
// @Tags auth
// @Accept json
// @Produce json
// @Param body body auth.ForgottenPasswordRequestModel true "Request body"
// @Success 200 {object} auth.ForgottenPasswordResponseModel "OK"
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 500 {object} response.Failure "Internal server error"
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
// @Failure 400 {object} response.ValidationFailure "Bad request"
// @Failure 500 {object} response.Failure "Internal server error"
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
