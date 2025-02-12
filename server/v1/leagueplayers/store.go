package leagueplayers

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/db"
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
		return nil, fmt.Errorf("quering league players: %v", err)
	}

	for rows.Next() {
		var pm players.PlayerModel
		err := pm.ScanRows(rows)
		if err != nil {
			return nil, err
		}

		dest = append(dest, pm)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating league player rows: %v", err)
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
		return dest, err
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
			set current_league_id = $1
			where id = $2
			returning id, height, weight, handedness, racket, matches_expected, matches_played, matches_won, matches_scheduled, seasons_played, account_id, current_league_id, created_at
		)
		select
			up.id,
			up.height,
			up.weight,
			up.handedness,
			up.racket,
			up.matches_expected,
			up.matches_played,
			up.matches_won,
			up.matches_scheduled,
			up.seasons_played,
			account.id as account_id,
			account.name as account_name,
			league.id as current_league_id,
			league.title as current_league_title,
			up.created_at
		from updated_player up
		join account on up.account_id = account.id
		left join league on up.current_league_id = league.id
	`

	var dest players.PlayerModel

	row := q.QueryRow(ctx, sql, leagueId, playerId)
	err := dest.ScanRow(row)
	if err != nil {
		// todo: better err msg
		return dest, err
	}

	return dest, nil
}

// helper
func (s *store) validateFindLeaguePlayers(ctx context.Context, seasonId, leagueId string) (seasonExists bool, leagueExists bool, leagueInSeason bool, err error) {
	sql1 := `
		select
			exists (
				select 1 from season where id = $1
			) as season_exists,
			exists (
				select 1 from league where id = $2
			) as league_exists,
			exists (
				select 1 from league where id = $2 and season_id = $1
			) as league_in_season
	`

	err = s.db.QueryRow(ctx, sql1, seasonId, leagueId).Scan(&seasonExists, &leagueExists, &leagueInSeason)
	if err != nil {
		return
	}

	return
}

// helper
func (s *store) validateFindLeaguePlayer(ctx context.Context, seasonId, leagueId, playerId string) (seasonExists bool, leagueExists bool, leagueInSeason bool, playerExists bool, err error) {
	sql1 := `
		select
			exists (
				select 1 from season where id = $1
			) as season_exists,
			exists (
				select 1 from league where id = $2
			) as league_exists,
			exists (
				select 1 from league where id = $2 and season_id = $1
			) as league_in_season,
			exists (
				select 1 from player where id = $3
			) as player_exists
	`

	err = s.db.QueryRow(ctx, sql1, seasonId, leagueId, playerId).Scan(&seasonExists, &leagueExists, &leagueInSeason, &playerExists)
	if err != nil {
		return
	}

	return
}

// helper
func (s *store) validateUpdatePlayerCurrentLeague(ctx context.Context, seasonId, leagueId, playerId string) (seasonExists bool, leagueExists bool, leagueInSeason bool, playerExists bool, err error) {
	sql1 := `
	select
		exists (
			select 1 from season where id = $1
		) as season_exists,
		exists (
			select 1 from league where id = $2
		) as league_exists,
		exists (
			select 1 from league where id = $2 and season_id = $1
		) as league_in_season,
		exists (
			select 1 from player where id = $3
		) as player_exists
`

	err = s.db.QueryRow(ctx, sql1, seasonId, leagueId, playerId).Scan(&seasonExists, &leagueExists, &leagueInSeason, &playerExists)
	if err != nil {
		return
	}

	return
}
