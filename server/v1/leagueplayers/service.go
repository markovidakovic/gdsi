package leagueplayers

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/v1/players"
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

func (s *service) processGetLeaguePlayers(ctx context.Context, seasonId, leagueId string) ([]players.PlayerModel, error) {
	// validation
	err := s.validator.NewValidation(ctx).
		SeasonExists(seasonId, "path").
		LeagueExists(leagueId, "path").
		LeagueInSeason(seasonId, leagueId, "path").
		Result()
	if err != nil {
		return nil, err
	}

	// find league players
	lps, err := s.store.findLeaguePlayers(ctx, leagueId)
	if err != nil {
		return nil, err
	}

	return lps, nil
}

func (s *service) processGetLeaguePlayer(ctx context.Context, seasonId, leagueId, playerId string) (*players.PlayerModel, error) {
	// validation
	err := s.validator.NewValidation(ctx).
		SeasonExists(seasonId, "path").
		LeagueExists(leagueId, "path").
		LeagueInSeason(seasonId, leagueId, "path").
		PlayerExists(playerId, "path").
		Result()
	if err != nil {
		return nil, err
	}

	// find league player
	lp, err := s.store.findLeaguePlayer(ctx, leagueId, playerId)
	if err != nil {
		return nil, err
	}

	return &lp, nil
}

func (s *service) processAssignPlayerToLeague(ctx context.Context, seasonId, leagueId, playerId string) (*players.PlayerModel, error) {
	// validation
	err := s.validator.NewValidation(ctx).
		SeasonExists(seasonId, "path").
		LeagueExists(leagueId, "path").
		LeagueInSeason(seasonId, leagueId, "path").
		PlayerExists(playerId, "path").
		Result()
	if err != nil {
		return nil, err
	}

	// begin tx
	tx, err := s.store.db.Begin(ctx)
	if err != nil {
		return nil, failure.New("unable to assign player to league", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err != pgx.ErrTxClosed {
			log.Printf("failed to rollback assign player to league tx: %v", err)
		}
	}()

	// update player current league
	player, err := s.store.updatePlayerCurrentLeague(ctx, tx, &leagueId, playerId)
	if err != nil {
		// another option to do here? could just return nil, err
		// the updatePlayerCurrentLeague will never return (i think so) "player not found"
		// because we check the player existance in the validation above.
		// so checking for errors.Is(err, failure.ErrNotFound) is not necessary
		// so we could just return a generic message from the service layer and not the specific from the store
		return nil, failure.New("unable to assign player to league", err)
	}

	// increment player seasons played
	player, err = s.store.incrementPlayerSeasonsPlayed(ctx, tx, leagueId, playerId)
	if err != nil {
		// same as above, this will most likely be a db error
		return nil, failure.New("unable to assign player to league", err)
	}

	// commit tx
	err = tx.Commit(ctx)
	if err != nil {
		return nil, failure.New("unable to assign player to league", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return &player, nil
}

func (s *service) processUnassignPlayerFromLeague(ctx context.Context, seasonId, leagueId, playerId string) (*players.PlayerModel, error) {
	// validation
	// todo: maybe do a validation in validation.go for playerInLeague
	err := s.validator.NewValidation(ctx).
		SeasonExists(seasonId, "path").
		LeagueExists(leagueId, "path").
		LeagueInSeason(seasonId, leagueId, "path").
		PlayerExists(playerId, "path").
		Result()
	if err != nil {
		return nil, err
	}

	// update player current league
	lp, err := s.store.updatePlayerCurrentLeague(ctx, nil, nil, playerId)
	if err != nil {
		// same as in other methods. this will most likely be a db error
		// and not a player not found because the player existance is confirmed in the validation above
		return nil, failure.New("unable to unassign player from league", err)
	}

	return &lp, nil
}
