package matches

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

func (s *service) processCreateMatch(ctx context.Context, input CreateMatchRequestModel) (*MatchModel, error) {
	// validate input in the db layer
	seasonExists, leagueExists, leagueInSeason, playerOneExists, playerTwoExists, playersInLeague, err := s.store.validateInsertMatch(ctx, input.SeasonId, input.LeagueId, input.PlayerOneId, input.PlayerTwoId)
	if err != nil {
		return nil, err
	}

	if !seasonExists {
		return nil, fmt.Errorf("finding season: %w", response.ErrNotFound)
	}
	if !leagueExists {
		return nil, fmt.Errorf("finding league: %w", response.ErrNotFound)
	}
	if !leagueInSeason {
		return nil, fmt.Errorf("league not in season: %w", response.ErrBadRequest)
	}
	if !playerOneExists {
		return nil, fmt.Errorf("finding player one: %w", response.ErrNotFound)
	}
	if !playerTwoExists {
		return nil, fmt.Errorf("finding player two: %w", response.ErrNotFound)
	}
	if !playersInLeague {
		return nil, fmt.Errorf("players not in league: %w", response.ErrBadRequest)
	}

	// insert match
	cm, err := s.store.insertMatch(ctx, input)
	if err != nil {
		return nil, err
	}

	return &cm, nil
}
