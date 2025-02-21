package leagueplayers

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/v1/players"
)

type store struct {
	db *db.Conn
}

func newStore(db *db.Conn) *store {
	return &store{
		db,
	}
}

func (s *store) findLeaguePlayers(ctx context.Context, leagueId string) ([]players.PlayerModel, error) {
	sql := `
		select
			player.id,
			player.height,
			player.weight,
			player.handedness,
			player.racket,
			player.matches_expected,
			player.matches_played,
			player.matches_won,
			player.matches_scheduled,
			player.seasons_played,
			account.id as account_id,
			account.name as account_name,
			league.id as current_league_id,
			league.title as current_league_title,
			player.created_at
		from player
		join account on player.account_id = account.id
		left join league on player.current_league_id = league.id
		where player.current_league_id = $1
		order by player.created_at desc
	`

	dest := []players.PlayerModel{}

	rows, err := s.db.Query(ctx, sql, leagueId)
	if err != nil {
		return nil, failure.New("unable to find league players", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	for rows.Next() {
		var pm players.PlayerModel
		err := pm.ScanRows(rows)
		if err != nil {
			return nil, failure.New("unable to find league players", err)
		}

		dest = append(dest, pm)
	}

	if err := rows.Err(); err != nil {
		return nil, failure.New("unable to find league players", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return dest, nil
}

func (s *store) findLeaguePlayer(ctx context.Context, leagueId, playerId string) (players.PlayerModel, error) {
	var dest players.PlayerModel

	sql := `
		select 
			player.id,
			player.height,
			player.weight,
			player.handedness,
			player.racket,
			player.matches_expected,
			player.matches_played,
			player.matches_won,
			player.matches_scheduled,
			player.seasons_played,
			account.id as account_id,
			account.name as account_name,
			league.id as current_league_id,
			league.title as current_league_title,
			player.created_at
		from player
		join account on player.account_id = account.id
		left join league on player.current_league_id = league.id
		where player.id = $1 and player.current_league_id = $2
	`

	row := s.db.QueryRow(ctx, sql, playerId, leagueId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return dest, failure.New("league player not found", err)
		}
		return dest, failure.New("unable to find league player", err)
	}

	return dest, nil
}

func (s *store) updatePlayerCurrentLeague(ctx context.Context, tx pgx.Tx, leagueId *string, playerId string) (players.PlayerModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	sql := `
		with updated_player as (
			update player
			set 
				current_league_id = $1
			where id = $2
			returning id, height, weight, handedness, racket, matches_expected, matches_played, matches_won, matches_scheduled, seasons_played, account_id, current_league_id, created_at
		)
		select
			up.id as player_id,
			up.height as player_height,
			up.weight as player_weight,
			up.handedness as player_handedness,
			up.racket as player_racket,
			up.matches_expected as player_matches_expected,
			up.matches_played as player_matches_played,
			up.matches_won as player_matches_won,
			up.matches_scheduled as player_matches_scheduled,
			up.seasons_played as player_seasons_played,
			account.id as player_account_id,
			account.name as player_account_name,
			league.id as player_current_league_id,
			league.title as player_current_league_title,
			up.created_at
		from updated_player up
		join account on up.account_id = account.id
		left join league on up.current_league_id = league.id
	`

	var dest players.PlayerModel

	row := q.QueryRow(ctx, sql, leagueId, playerId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return dest, failure.New("league player for updating current league not found", err)
		}
		return dest, failure.New("unable to update league player current league", err)
	}

	return dest, nil
}

func (s *store) incrementPlayerSeasonsPlayed(ctx context.Context, tx pgx.Tx, leagueId, playerId string) (players.PlayerModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	sql := `
		with updated_player as (
			update player
			set
				seasons_played = seasons_played + 1
			where id = $1 and current_league_id = $2
			returning id, height, weight, handedness, racket, matches_expected, matches_played, matches_won, matches_scheduled, seasons_played, account_id, current_league_id, created_at
		)
		select
			up.id as player_id,
			up.height as player_height,
			up.weight as player_weight,
			up.handedness as player_handedness,
			up.racket as player_racket,
			up.matches_expected as player_matches_expected,
			up.matches_played as player_matches_played,
			up.matches_won as player_matches_won,
			up.matches_scheduled as player_matches_scheduled,
			up.seasons_played as player_seasons_played,
			account.id as player_account_id,
			account.name as player_account_name,
			league.id as player_current_league_id,
			league.title as player_current_league_title,
			up.created_at as player_created_at
		from updated_player up
		join account on up.account_id = account.id
		left join league on up.current_league_id = league.id
	`

	var dest players.PlayerModel
	row := q.QueryRow(ctx, sql, playerId, leagueId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return dest, failure.New("league player for incrementing seasons played not found", err)
		}
		return dest, failure.New("unable to increment league player seasons played", err)
	}

	return dest, nil
}
