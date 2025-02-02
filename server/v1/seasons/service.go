package seasons

import (
	"context"

	"github.com/markovidakovic/gdsi/server/config"
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

func (s *service) processCreateSeason(ctx context.Context, input CreateSeasonRequestModel) (SeasonModel, error) {
	return SeasonModel{}, nil
}
