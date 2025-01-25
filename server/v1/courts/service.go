package courts

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

func (s *service) createCourt(ctx context.Context, input CreateCourtModel) (CourtModel, error) {
	input.CreatorId = ctx.Value(middleware.AccountIdCtxKey).(string)

	result, err := s.store.insertCourt(ctx, input)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *service) getCourts(ctx context.Context) ([]CourtModel, error) {
	result, err := s.store.findCourts(ctx)
	if err != nil {
		return result, err
	}
	return result, nil
}
