package standings

import (
	"context"

	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/validation"
)

type service struct {
	cfg       *config.Config
	store     *store
	validator *validation.Validator
}

func newService(cfg *config.Config, store *store, validator *validation.Validator) *service {
	return &service{
		cfg,
		store,
		validator,
	}
}

func (s *service) processGetStandings(ctx context.Context, seasonId, leagueId string) ([]StandingModel, error) {
	// validation
	err := s.validator.NewValidation(ctx).
		SeasonExists(seasonId).
		LeagueExists(leagueId).
		LeagueInSeason(seasonId, leagueId).
		Result()
	if err != nil {
		return nil, err
	}

	// call the store
	standings, err := s.store.findStandings(ctx, seasonId, leagueId)
	if err != nil {
		return nil, err
	}

	return standings, nil
}
