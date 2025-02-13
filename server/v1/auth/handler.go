package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/response"
)

type handler struct {
	service *service
}

func newHandler(cfg *config.Config, db *db.Conn) *handler {
	h := &handler{}
	store := newStore(db)
	h.service = newService(cfg, store)
	return h
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
func (h *handler) signup(w http.ResponseWriter, r *http.Request) {
	var model SignupRequestModel

	// decode request body
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	// validation
	if valErr := model.Validate(); valErr != nil {
		response.WriteFailure(w, response.NewValidationFailure("validation failed", valErr))
		return
	}

	// call the service
	accessToken, refreshToken, err := h.service.processSignup(r.Context(), model)
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
func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var model LoginRequestModel

	// decode request body
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	// validate model
	if valErr := model.Validate(); valErr != nil {
		response.WriteFailure(w, response.NewValidationFailure("validation failed", valErr))
		return
	}

	// for now not used, just to see in the future what to maybe store
	// ra := r.RemoteAddr
	// host, port, err := net.SplitHostPort(ra)
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// }
	// fmt.Printf("ip: %v\n", ra)
	// fmt.Printf("host: %v\n", host)
	// fmt.Printf("port: %v\n", port)
	// fmt.Printf("r.UserAgent(): %v\n", r.UserAgent())

	// Call the service
	accessToken, refreshToken, err := h.service.processLogin(r.Context(), model)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
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
func (h *handler) refreshToken(w http.ResponseWriter, r *http.Request) {
	var model RefreshTokenRequestModel

	// decode request body
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		response.WriteFailure(w, response.NewBadRequestFailure("invalid request body"))
		return
	}

	// validate model
	if valErr := model.Validate(); valErr != nil {
		response.WriteFailure(w, response.NewValidationFailure("validation failed", valErr))
		return
	}

	access, refresh, err := h.service.processRefreshTokens(r.Context(), model)
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
