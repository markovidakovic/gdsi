package leagues

import (
	"context"
	"fmt"

	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/response"
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

func (s *service) processCreateLeague(ctx context.Context, input CreateLeagueRequestModel) (*LeagueModel, error) {
	// check if season exists
	seasonExists, err := s.store.checkSeasonExistance(ctx, input.SeasonId)
	if err != nil {
		return nil, err
	}

	if !seasonExists {
		return nil, fmt.Errorf("finding season: %w", response.ErrNotFound)
	}

	// create league
	lm, err := s.store.insertLeague(ctx, input)
	if err != nil {
		return nil, err
	}

	return &lm, nil
}

func (s *service) processFindLeagues(ctx context.Context, seasonId string) ([]LeagueModel, error) {
	// check if season exists
	seasonExists, err := s.store.checkSeasonExistance(ctx, seasonId)
	if err != nil {
		return nil, err
	}

	if !seasonExists {
		return nil, fmt.Errorf("finding season: %w", response.ErrNotFound)
	}

	// find leagues
	lms, err := s.store.findLeagues(ctx, seasonId)
	if err != nil {
		return nil, err
	}

	return lms, err
}

func (s *service) processFindLeague(ctx context.Context, seasonId, leagueId string) (*LeagueModel, error) {
	// check if season exists
	seasonExists, err := s.store.checkSeasonExistance(ctx, seasonId)
	if err != nil {
		return nil, err
	}

	if !seasonExists {
		return nil, fmt.Errorf("finding season: %w", response.ErrNotFound)
	}

	// find league
	lm, err := s.store.findLeague(ctx, seasonId, leagueId)
	if err != nil {
		return nil, err
	}

	return lm, nil
}

func (s *service) processUpdateLeague(ctx context.Context, input UpdateLeagueRequestModel) (*LeagueModel, error) {
	// check if season exists
	seasonExists, err := s.store.checkSeasonExistance(ctx, input.SeasonId)
	if err != nil {
		return nil, err
	}

	if !seasonExists {
		return nil, fmt.Errorf("finding season: %w", response.ErrNotFound)
	}

	// update league
	lm, err := s.store.updateLeague(ctx, input)
	if err != nil {
		return nil, err
	}

	return &lm, nil
}

func (s *service) processDeleteLeague(ctx context.Context, seasonId, leagueId string) error {
	// check if season exists
	seasonExists, err := s.store.checkSeasonExistance(ctx, seasonId)
	if err != nil {
		return err
	}

	if !seasonExists {
		return fmt.Errorf("finding season: %w", response.ErrNotFound)
	}

	// delete league
	err = s.store.deleteLeague(ctx, seasonId, leagueId)
	if err != nil {
		return err
	}

	return nil
}
