package players

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

func (s *service) processGetLeaguePlayers(ctx context.Context, seasonId, leagueId string) ([]PlayerModel, error) {
	// validate params
	seasonExists, leagueExists, leagueInSeason, err := s.store.validateFindLeaguePlayers(ctx, seasonId, leagueId)
	if err != nil {
		return nil, err
	}

	if !seasonExists {
		return nil, fmt.Errorf("%w: finding season", response.ErrNotFound)
	}
	if !leagueExists {
		return nil, fmt.Errorf("%w: finding league", response.ErrNotFound)
	}
	if !leagueInSeason {
		return nil, fmt.Errorf("%w: league not in season", response.ErrBadRequest)
	}

	// find league players
	lps, err := s.store.findLeaguePlayers(ctx, leagueId)
	if err != nil {
		return nil, err
	}

	return lps, nil
}

func (s *service) processGetLeaguePlayer(ctx context.Context, seasonId, leagueId, playerId string) (*PlayerModel, error) {
	// validate params
	seasonExists, leagueExists, leagueInSeason, playerExists, err := s.store.validateFindLeaguePlayer(ctx, seasonId, leagueId, playerId)
	if err != nil {
		return nil, err
	}

	if !seasonExists {
		return nil, fmt.Errorf("%w: finding season", response.ErrNotFound)
	}
	if !leagueExists {
		return nil, fmt.Errorf("%w: finding league", response.ErrNotFound)
	}
	if !leagueInSeason {
		return nil, fmt.Errorf("%w: league not in season", response.ErrBadRequest)
	}
	if !playerExists {
		return nil, fmt.Errorf("%w: finding player", response.ErrNotFound)
	}

	// find league player
	lp, err := s.store.findLeaguePlayer(ctx, leagueId, playerId)
	if err != nil {
		return nil, err
	}

	return &lp, nil
}
