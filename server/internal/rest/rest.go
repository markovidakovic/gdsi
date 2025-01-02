package rest

import "github.com/markovidakovic/gdsi/server/internal/config"

type server struct {
	cfg *config.Config
}

func New() (*server, error) {
	var srvr = &server{}

	return srvr, nil
}
