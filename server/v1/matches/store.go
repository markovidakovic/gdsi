package matches

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/response"
)

type store struct {
	db *db.Conn
}

func newStore(db *db.Conn) *store {
	return &store{
		db,
	}
}

func (s *store) insertMatch(ctx context.Context, input CreateMatchRequestModel) (MatchModel, error) {
	sql1 := `
		with inserted_match as (
			insert into match (court_id, scheduled_at, player_one_id, player_two_id, winner_id, score, season_id, league_id)
			values ($1, $2, $3, $4, $5, $6, $7, $8)
			returning id, court_id, scheduled_at, player_one_id, player_two_id, winner_id, score, season_id, league_id, created_at
		)
		select
			im.id,
			court.id as court_id,
			court.name as court_name,
			im.scheduled_at,
			player1.id as player_one_id,
			account1.name as player_one_name,
			player2.id as player_two_id,
			account2.name as player_two_name,
			winner.id as winner_id,
			account3.name as winner_name,
			im.score,
			season.id as season_id,
			season.title as season_title,
			league.id as league_id,
			league.title as league_title,
			im.created_at
		from inserted_match im
		join court on im.court_id = court.id
		join player player1 on im.player_one_id = player1.id
		join account account1 on player1.account_id = account1.id
		join player player2 on im.player_two_id = player2.id
		join account account2 on player2.account_id = account2.id
		left join player winner on im.winner_id = winner.id
		left join account account3 on winner.account_id = account3.id
		join season on im.season_id = season.id
		join league on im.league_id = league.id
	`

	var dest MatchModel
	var winnerId, winnerName sql.NullString

	err := s.db.QueryRow(ctx, sql1, input.CourtId, input.ScheduledAt, input.PlayerOneId, input.PlayerTwoId, input.WinnerId, input.Score, input.SeasonId, input.LeagueId).Scan(
		&dest.Id,
		&dest.Court.Id,
		&dest.Court.Name,
		&dest.ScheduledAt,
		&dest.PlayerOne.Id,
		&dest.PlayerOne.Name,
		&dest.PlayerTwo.Id,
		&dest.PlayerTwo.Name,
		&winnerId,
		&winnerName,
		&dest.Score,
		&dest.Season.Id,
		&dest.Season.Title,
		&dest.League.Id,
		&dest.League.Title,
		&dest.CreatedAt,
	)
	if err != nil {
		return dest, err
	}

	if !winnerId.Valid {
		dest.Winner = nil
	} else {
		dest.Winner = &PlayerModel{
			Id:   winnerId.String,
			Name: winnerName.String,
		}
	}

	return dest, nil
}

func (s *store) findMatches(ctx context.Context, seasonId, leagueId string) ([]MatchModel, error) {
	sql1 := `
		select
			match.id,
			court.id as court_id,
			court.name as court_name,
			match.scheduled_at,
			player1.id as player_one_id,
			account1.name as player_one_name,
			player2.id as player_two_id,
			account2.name as player_two_name,
			winner.id as winner_id,
			account3.name as winner_name,
			match.score,
			season.id as season_id,
			season.title as season_title,
			league.id as league_id,
			league.title as league_title,
			match.created_at
		from match
		join court on match.court_id = court.id
		join player player1 on match.player_one_id = player1.id
		join account account1 on player1.account_id = account1.id
		join player player2 on match.player_two_id = player2.id
		join account account2 on player2.account_id = account2.id
		left join player winner on match.winner_id = winner.id
		left join account account3 on winner.account_id = account3.id
		join season on match.season_id = season.id
		join league on match.league_id = league.id
		where match.season_id = $1 and match.league_id = $2
		order by match.created_at desc
	`

	dest := []MatchModel{}

	rows, err := s.db.Query(ctx, sql1, seasonId, leagueId)
	if err != nil {
		return nil, fmt.Errorf("quering match rows: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mm MatchModel
		var winnerId, winnerName sql.NullString
		err := rows.Scan(&mm.Id, &mm.Court.Id, &mm.Court.Name, &mm.ScheduledAt, &mm.PlayerOne.Id, &mm.PlayerOne.Name, &mm.PlayerTwo.Id, &mm.PlayerTwo.Name, &winnerId, &winnerName, &mm.Score, &mm.Season.Id, &mm.Season.Title, &mm.League.Id, &mm.League.Title, &mm.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scanning match row: %v", err)
		}

		if !winnerId.Valid {
			mm.Winner = nil
		} else {
			mm.Winner = &PlayerModel{
				Id:   winnerId.String,
				Name: winnerName.String,
			}
		}

		dest = append(dest, mm)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating match rows: %v", err)
	}

	return dest, nil
}

func (s *store) findMatch(ctx context.Context, seasonId, leagueId, matchId string) (*MatchModel, error) {
	sql1 := `
		select
			match.id,
			court.id as court_id,
			court.name as court_name,
			match.scheduled_at,
			player1.id as player_one_id,
			account1.name as player_one_name,
			player2.id as player_two_id,
			account2.name as player_two_name,
			winner.id as winner_id,
			account3.name as winner_name,
			match.score,
			season.id as season_id,
			season.title as season_title,
			league.id as league_id,
			league.title as league_title,
			match.created_at
		from match
		join court on match.court_id = court.id
		join player player1 on match.player_one_id = player1.id
		join account account1 on player1.account_id = account1.id
		join player player2 on match.player_two_id = player2.id
		join account account2 on player2.account_id = account2.id
		left join player winner on match.winner_id = winner.id
		left join account account3 on winner.account_id = account3.id
		join season on match.season_id = season.id
		join league on match.league_id = league.id
		where match.id = $1 and match.season_id = $2 and match.league_id = $3
	`

	var dest MatchModel
	var winnerId, winnerName sql.NullString

	err := s.db.QueryRow(ctx, sql1, matchId, seasonId, leagueId).Scan(&dest.Id, &dest.Court.Id, &dest.Court.Name, &dest.ScheduledAt, &dest.PlayerOne.Id, &dest.PlayerOne.Name, &dest.PlayerTwo.Id, &dest.PlayerTwo.Name, &winnerId, &winnerName, &dest.Score, &dest.Season.Id, &dest.Season.Title, &dest.League.Id, &dest.League.Title, &dest.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("finding match: %w", response.ErrNotFound)
		}
		return nil, err
	}

	if !winnerId.Valid {
		dest.Winner = nil
	} else {
		dest.Winner = &PlayerModel{
			Id:   winnerId.String,
			Name: winnerName.String,
		}
	}

	return &dest, nil
}

func (s *store) updateMatch(ctx context.Context, input UpdateMatchRequestModel) (*MatchModel, error) {
	sql1 := `
		with updated_match as (
			update match 
			set court_id = $1, scheduled_at = $2, player_two_id = $3
			where id = $4 and season_id = $5 and league_id = $6
			returning id, court_id, scheduled_at, player_one_id, player_two_id, winner_id, score, season_id, league_id, created_at
		)
		select
			um.id,
			court.id as court_id,
			court.name as court_name,
			um.scheduled_at,
			player1.id as player_one_id,
			account1.name as player_one_name,
			player2.id as player_two_id,
			account2.name as player_two_name,
			winner.id as winner_id,
			account3.name as winner_name,
			um.score,
			season.id as season_id,
			season.title as season_title,
			league.id as league_id,
			league.title as league_title,
			um.created_at
		from updated_match um
		join court on um.court_id = court.id
		join player player1 on um.player_one_id = player1.id
		join account account1 on player1.account_id = account1.id
		join player player2 on um.player_two_id = player2.id
		join account account2 on player2.account_id = account2.id
		left join player winner on um.winner_id = winner.id
		left join account account3 on winner.account_id = account3.id
		join season on um.season_id = season.id
		join league on um.league_id = league.id
	`

	var dest MatchModel
	var winnerId, winnerName sql.NullString

	err := s.db.QueryRow(ctx, sql1, input.CourtId, input.ScheduledAt, input.PlayerTwoId, input.MatchId, input.SeasonId, input.LeagueId).Scan(&dest.Id, &dest.Court.Id, &dest.Court.Name, &dest.ScheduledAt, &dest.PlayerOne.Id, &dest.PlayerOne.Name, &dest.PlayerTwo.Id, &dest.PlayerTwo.Name, &winnerId, &winnerName, &dest.Score, &dest.Season.Id, &dest.Season.Title, &dest.League.Id, &dest.League.Title, &dest.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("finding match: %w", response.ErrNotFound)
		}
		return nil, err
	}

	if !winnerId.Valid {
		dest.Winner = nil
	} else {
		dest.Winner = &PlayerModel{
			Id:   winnerId.String,
			Name: winnerName.String,
		}
	}

	return &dest, nil
}

func (s *store) submitMatchScore(ctx context.Context) {

}

// helper
func (s *store) validateInsertUpdateMatch(ctx context.Context, courtId, seasonId, leagueId, player1Id, player2Id string) (courtExists bool, seasonExists bool, leagueExists bool, leagueInSeason bool, playerOneExists bool, playerTwoExists bool, playersInLeague bool, err error) {
	sql1 := `
		select
			exists (
				select 1 from court where id = $1
			) as court_exists,
			exists ( 
				select 1 from season where id = $2
			) as season_exists,
			exists (
				select 1 from league where id = $3
			) as league_exists,
			exists (
				select 1 from league where id = $3 and season_id = $2
				) as league_in_season,
			exists (
				select 1 from player where id = $4
			) as player_one_exists,
			exists (
				select 1 from player where id = $5
			) as player_two_exists,
			exists (
				select 1 from player
				where id in ($4, $5)
				and current_league_id = $3
				having count(*) = 2
			) as players_in_league
	`

	err = s.db.QueryRow(ctx, sql1, courtId, seasonId, leagueId, player1Id, player2Id).Scan(&courtExists, &seasonExists, &leagueExists, &leagueInSeason, &playerOneExists, &playerTwoExists, &playersInLeague)
	if err != nil {
		return
	}

	return
}

// helper
func (s *store) validateFindMatches(ctx context.Context, seasonId, leagueId string) (seasonExists bool, leagueExists bool, leagueInSeason bool, err error) {
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
func (s *store) checkMatchParticipation(ctx context.Context, matchId, accountId string) (bool, error) {
	sql := `
		select exists (
			select 1 from match
			where id = $1 and (player_one_id = $2 or player_two_id = $2)
		)
	`

	var exists bool
	err := s.db.QueryRow(ctx, sql, matchId, accountId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
