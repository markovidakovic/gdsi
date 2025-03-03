package leagueplayers

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/params"
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

var allowedSortFeilds = map[string]string{
	"created_at": "player.created_at",
}

func (s *store) findLeaguePlayers(ctx context.Context, leagueId string, requestingPlayerId string, matchAvailable bool, limit, offset int, sort *params.OrderBy) ([]players.PlayerModel, error) {
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
	`

	args := []interface{}{leagueId}
	argCounter := 2 // starting with $2 since $1 is already used for current_league_id

	if matchAvailable {
		// exclude the player requesting
		sql += fmt.Sprintf(" and player.id != $%d", argCounter)
		args = append(args, requestingPlayerId)
		argCounter++

		// exclude players who have already played a match with the requesting player
		sql += fmt.Sprintf(`
			and not exists (
				select 1
				from match
				where (
					(match.player_one_id = player.id and match.player_two_id = $%d)
					or
					(match.player_one_id = $%d and match.player_two_id = player.id)
				)
				and match.league_id = $1
			)
		`, argCounter, argCounter)
		args = append(args, requestingPlayerId)
		argCounter++
	}

	if sort != nil && sort.IsValid(allowedSortFeilds) {
		sql += fmt.Sprintf("order by %s %s\n", allowedSortFeilds[sort.Field], sort.Direction)
	} else {
		sql += fmt.Sprintln("order by player.created_at desc")
	}

	if limit >= 0 {
		sql += fmt.Sprintf("limit $%d offset $%d", argCounter, argCounter+1)
		args = append(args, limit, offset)
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, failure.New("unable to find league players", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	defer rows.Close()

	dest := []players.PlayerModel{}
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

func (s *store) countLeaguePlayers(ctx context.Context, leagueId string) (int, error) {
	var count int
	sql := `select count(*) from player where current_league_id = $1`
	err := s.db.QueryRow(ctx, sql, leagueId).Scan(&count)
	if err != nil {
		return 0, failure.New("unable to count league players", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	return count, nil
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
