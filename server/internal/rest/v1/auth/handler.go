package auth

import (
	"log"
	"net/http"

	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
)

type handler struct {
	service *service
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login endpoint\n"))
}

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

	log.Println("auth endpoints mounted")
	return h
}
