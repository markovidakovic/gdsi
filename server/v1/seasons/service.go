package seasons

import (
	"context"

	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/middleware"
	"github.com/markovidakovic/gdsi/server/params"
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

	sm, err := s.store.insertSeason(ctx, nil, model)
	if err != nil {
		return sm, err
	}

	return sm, nil
}

func (s *service) processGetSeasons(ctx context.Context, query *params.Query) ([]SeasonModel, int, error) {
	count, err := s.store.countSeasons(ctx)
	if err != nil {
		return nil, 0, failure.New("unable to get seasons", err)
	}

	limit, offset := query.CalcLimitAndOffset(count)

	result, err := s.store.findSeasons(ctx, limit, offset, query.OrderBy)
	if err != nil {
		return nil, 0, err
	}

	return result, count, nil
}
