package standings

import "github.com/markovidakovic/gdsi/server/internal/config"

type service struct {
	cfg   *config.Config
	store *store
}

func newService(cfg *config.Config, store *store) *service {
	return &service{
		cfg,
		store,
	}
}
