package auth

import (
	"context"

	"github.com/markovidakovic/gdsi/server/internal/config"
)

type service struct {
	cfg   *config.Config
	store *store
}

func (s *service) signupNewAccount(ctx context.Context) (string, error) {
	result, err := s.store.createNewAccount(ctx)
	if err != nil {
		return "", err
	}
	return result, err
}

func newService(cfg *config.Config, store *store) *service {
	var s = &service{
		cfg,
		store,
	}
	return s
}
