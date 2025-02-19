package standings

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

func (s *service) processGetStandings(ctx context.Context, seasonId, leagueId string) ([]StandingModel, error) {
	// call the store
	standings, err := s.store.findStandings(ctx, seasonId, leagueId)
	if err != nil {
		return nil, err
	}

	return standings, nil
}
