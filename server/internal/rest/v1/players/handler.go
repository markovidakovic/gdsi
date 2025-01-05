package players

import (
	"log"
	"net/http"

	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
)

type handler struct {
	service *service
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.getAllPlayers(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch players", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(result))
}

func NewHandler(cfg *config.Config, db db.Conn) *handler {
	h := &handler{}

	store := newStore(db)
	service := newService(store)

	h.service = service

	log.Println("players endpoints mounted")
	return h
}
