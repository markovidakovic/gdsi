package courts

import (
	"context"

	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/pagination"
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

func (s *service) processGetCourts(ctx context.Context, query *pagination.QueryParams) ([]CourtModel, int, error) {
	count, err := s.store.countCourts(ctx)
	if err != nil {
		return nil, 0, failure.New("unable to get courts", err)
	}

	limit, offset := query.CalcLimitAndOffset(count)

	result, err := s.store.findCourts(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return result, count, nil
}
