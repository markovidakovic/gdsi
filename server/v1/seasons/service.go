package seasons

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

func (s *service) processCreateSeason(ctx context.Context, model CreateSeasonRequestModel) (SeasonModel, error) {
	model.CreatorId = ctx.Value(middleware.AccountIdCtxKey).(string)

	// call the store
	sm, err := s.store.insertSeason(ctx, nil, model)
	if err != nil {
		return sm, err
	}

	return sm, nil
}
