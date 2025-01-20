package me

import (
	"context"

	"github.com/markovidakovic/gdsi/server/config"
)

type service struct {
	cfg   *config.Config
	store *store
}

func (s *service) getMe(ctx context.Context, accountId string) (*MeModel, error) {
	me, err := s.store.queryMe(ctx, accountId)
	if err != nil {
		return nil, err
	}
	return me, nil
}

func newService(cfg *config.Config, store *store) *service {
	return &service{
		cfg,
		store,
	}
}
