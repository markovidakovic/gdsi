package me

import (
	"context"

	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/middleware"
)

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

func (s *service) getMe(ctx context.Context) (*MeModel, error) {
	accountId := ctx.Value(middleware.AccountIdCtxKey).(string)
	me, err := s.store.queryMe(ctx, accountId)
	if err != nil {
		return nil, err
	}
	return me, nil
}
