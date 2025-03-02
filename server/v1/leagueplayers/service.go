package leagueplayers

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/params"
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

func (s *service) processGetLeaguePlayers(ctx context.Context, seasonId, leagueId string, query *params.Query) ([]players.PlayerModel, int, error) {
	err := s.validator.NewValidation(ctx).
		SeasonExists(seasonId, "path").
		LeagueExists(leagueId, "path").
		LeagueInSeason(seasonId, leagueId, "path").
		Result()
	if err != nil {
		return nil, 0, err
	}

	count, err := s.store.countLeaguePlayers(ctx, leagueId)
	if err != nil {
		return nil, 0, failure.New("unable to get league players", err)
	}

	limit, offset := query.CalcLimitAndOffset(count)

	lps, err := s.store.findLeaguePlayers(ctx, leagueId, limit, offset, query.OrderBy)
	if err != nil {
		return nil, 0, err
	}

	return lps, count, nil
}

func (s *service) processGetLeaguePlayer(ctx context.Context, seasonId, leagueId, playerId string) (*players.PlayerModel, error) {
	err := s.validator.NewValidation(ctx).
		SeasonExists(seasonId, "path").
		LeagueExists(leagueId, "path").
		LeagueInSeason(seasonId, leagueId, "path").
		PlayerExists(playerId, "path").
		Result()
	if err != nil {
		return nil, err
	}

	lp, err := s.store.findLeaguePlayer(ctx, leagueId, playerId)
	if err != nil {
		return nil, err
	}

	return &lp, nil
}

func (s *service) processAssignPlayerToLeague(ctx context.Context, seasonId, leagueId, playerId string) (*players.PlayerModel, error) {
	err := s.validator.NewValidation(ctx).
		SeasonExists(seasonId, "path").
		LeagueExists(leagueId, "path").
		LeagueInSeason(seasonId, leagueId, "path").
		PlayerExists(playerId, "path").
		Result()
	if err != nil {
		return nil, err
	}

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

	player, err := s.store.updatePlayerCurrentLeague(ctx, tx, &leagueId, playerId)
	if err != nil {
		// another option to do here? could just return nil, err
		// the updatePlayerCurrentLeague will never return (i think so) "player not found"
		// because we check the player existance in the validation above.
		// so checking for errors.Is(err, failure.ErrNotFound) is not necessary
		// so we could just return a generic message from the service layer and not the specific from the store
		return nil, failure.New("unable to assign player to league", err)
	}

	player, err = s.store.incrementPlayerSeasonsPlayed(ctx, tx, leagueId, playerId)
	if err != nil {
		// same as above, this will most likely be a db error
		return nil, failure.New("unable to assign player to league", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, failure.New("unable to assign player to league", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return &player, nil
}

func (s *service) processUnassignPlayerFromLeague(ctx context.Context, seasonId, leagueId, playerId string) (*players.PlayerModel, error) {
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

	lp, err := s.store.updatePlayerCurrentLeague(ctx, nil, nil, playerId)
	if err != nil {
		// same as in other methods. this will most likely be a db error
		// and not a player not found because the player existance is confirmed in the validation above
		return nil, failure.New("unable to unassign player from league", err)
	}

	return &lp, nil
}
