package leagueplayers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/response"
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
	sql1 := `
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

	rows, err := s.db.Query(ctx, sql1, leagueId)
	if err != nil {
		return nil, fmt.Errorf("quering league players: %v", err)
	}

	for rows.Next() {
		var pm players.PlayerModel
		var currLeagueId, currLeagueTitle sql.NullString
		err := rows.Scan(&pm.Id, &pm.Height, &pm.Weight, &pm.Handedness, &pm.Racket, &pm.MatchesExpected, &pm.MatchesPlayed, &pm.MatchesWon, &pm.MatchesScheduled, &pm.SeasonsPlayed, &pm.Account.Id, &pm.Account.Name, &currLeagueId, &currLeagueTitle, &pm.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scanning league player row: %v", err)
		}

		if !currLeagueId.Valid {
			pm.CurrentLeague = nil
		} else {
			pm.CurrentLeague = &players.CurrentLeagueModel{
				Id:    currLeagueId.String,
				Title: currLeagueTitle.String,
			}
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
	var currLeagueId, currLeagueTitle sql.NullString

	sql1 := `
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

	err := s.db.QueryRow(ctx, sql1, playerId, leagueId).Scan(&dest.Id, &dest.Height, &dest.Weight, &dest.Handedness, &dest.Racket, &dest.MatchesExpected, &dest.MatchesPlayed, &dest.MatchesWon, &dest.MatchesScheduled, &dest.SeasonsPlayed, &dest.Account.Id, &dest.Account.Name, &currLeagueId, &currLeagueTitle, &dest.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dest, fmt.Errorf("finding league player: %w", response.ErrNotFound)
		}
		return dest, err
	}

	if !currLeagueId.Valid {
		dest.CurrentLeague = nil
	} else {
		dest.CurrentLeague = &players.CurrentLeagueModel{
			Id:    currLeagueId.String,
			Title: currLeagueTitle.String,
		}
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

	sql1 := `
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
	var currLeagueId, currLeagueTitle sql.NullString

	err := q.QueryRow(ctx, sql1, leagueId, playerId).Scan(&dest.Id, &dest.Height, &dest.Weight, &dest.Handedness, &dest.Racket, &dest.MatchesExpected, &dest.MatchesPlayed, &dest.MatchesWon, &dest.MatchesScheduled, &dest.SeasonsPlayed, &dest.Account.Id, &dest.Account.Name, &currLeagueId, &currLeagueTitle, &dest.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dest, fmt.Errorf("finding league player: %w", response.ErrNotFound)
		}
		return dest, err
	}

	if !currLeagueId.Valid {
		dest.CurrentLeague = nil
	} else {
		dest.CurrentLeague = &players.CurrentLeagueModel{
			Id:    currLeagueId.String,
			Title: currLeagueTitle.String,
		}
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
