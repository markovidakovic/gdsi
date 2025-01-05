package auth

import (
	"net/http"

	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
)

type handler struct {
	service *service
}

// @Summary Login
// @Description Login and get a new access token
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Success 200
// @Failure 500
// @Router /v1/auth/tokens/access [post]
func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login endpoint\n"))
}

// @Summary Signup
// @Description Signup a new account
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Success 200
// @Failure 500
// @Router /v1/auth/signup [post]
func (h *handler) Signup(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.signupNewAccount(r.Context())
	if err != nil {
		http.Error(w, "failed to signup new account", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(result))
}

func NewHandler(cfg *config.Config, db db.Conn) *handler {
	h := &handler{}

	store := newStore(db)
	service := newService(cfg, store)

	h.service = service

	return h
}
