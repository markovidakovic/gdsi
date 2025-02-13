package leagueplayers

import (
	"context"
	"fmt"

	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/response"
	"github.com/markovidakovic/gdsi/server/v1/players"
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

func (s *service) processGetLeaguePlayers(ctx context.Context, seasonId, leagueId string) ([]players.PlayerModel, error) {
	// validate params
	seasonExists, leagueExists, err := s.store.validateFindLeaguePlayers(ctx, seasonId, leagueId)
	if err != nil {
		return nil, err
	}

	if !seasonExists {
		return nil, fmt.Errorf("finding season: %w", response.ErrNotFound)
	}
	if !leagueExists {
		return nil, fmt.Errorf("finding league: %w", response.ErrNotFound)
	}

	// find league players
	lps, err := s.store.findLeaguePlayers(ctx, leagueId)
	if err != nil {
		return nil, err
	}

	return lps, nil
}

func (s *service) processGetLeaguePlayer(ctx context.Context, seasonId, leagueId, playerId string) (*players.PlayerModel, error) {
	// validate params
	seasonExists, leagueExists, playerExists, err := s.store.validateFindLeaguePlayer(ctx, seasonId, leagueId, playerId)
	if err != nil {
		return nil, err
	}

	if !seasonExists {
		return nil, fmt.Errorf("finding season: %w", response.ErrNotFound)
	}
	if !leagueExists {
		return nil, fmt.Errorf("finding league: %w", response.ErrNotFound)
	}
	if !playerExists {
		return nil, fmt.Errorf("finding player: %w", response.ErrNotFound)
	}

	// find league player
	lp, err := s.store.findLeaguePlayer(ctx, leagueId, playerId)
	if err != nil {
		return nil, err
	}

	return &lp, nil
}

func (s *service) processAssignPlayerToLeague(ctx context.Context, seasonId, leagueId, playerId string) (*players.PlayerModel, error) {
	// validate params
	seasonExists, leagueExists, playerExists, err := s.store.validateUpdatePlayerCurrentLeague(ctx, seasonId, leagueId, playerId)
	if err != nil {
		return nil, err
	}

	if !seasonExists {
		return nil, fmt.Errorf("finding season: %w", response.ErrNotFound)
	}
	if !leagueExists {
		return nil, fmt.Errorf("finding league: %w", response.ErrNotFound)
	}
	if !playerExists {
		return nil, fmt.Errorf("finding player: %w", response.ErrNotFound)
	}

	// update player current league
	lp, err := s.store.updatePlayerCurrentLeague(ctx, nil, &leagueId, playerId)
	if err != nil {
		return nil, err
	}

	return &lp, nil
}

func (s *service) processUnassignPlayerFromLeague(ctx context.Context, seasonId, leagueId, playerId string) (*players.PlayerModel, error) {
	// validate params
	seasonExists, leagueExists, playerExists, err := s.store.validateUpdatePlayerCurrentLeague(ctx, seasonId, leagueId, playerId)
	if err != nil {
		return nil, err
	}

	if !seasonExists {
		return nil, fmt.Errorf("finding season: %w", response.ErrNotFound)
	}
	if !leagueExists {
		return nil, fmt.Errorf("finding league: %w", response.ErrNotFound)
	}
	if !playerExists {
		return nil, fmt.Errorf("finding player: %w", response.ErrNotFound)
	}

	// update player current league
	lp, err := s.store.updatePlayerCurrentLeague(ctx, nil, nil, playerId)
	if err != nil {
		return nil, err
	}

	return &lp, nil
}
