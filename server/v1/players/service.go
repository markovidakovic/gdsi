package players

import (
	"context"

	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/failure"
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

func (s *service) processGetPlayers(ctx context.Context, query *params.Query) ([]PlayerModel, int, error) {
	count, err := s.store.countPlayers(ctx)
	if err != nil {
		return nil, 0, failure.New("unable to get players", err)
	}

	limit, offset := query.CalcLimitAndOffset(count)

	result, err := s.store.findPlayers(ctx, limit, offset, query.OrderBy)
	if err != nil {
		return nil, 0, err
	}

	return result, count, nil
}
