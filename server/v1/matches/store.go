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

func (s *store) insertMatch(ctx context.Context, tx pgx.Tx, model CreateMatchRequestModel) (MatchModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	sql := `
		with inserted_match as (
			insert into match (court_id, scheduled_at, player_one_id, player_two_id, winner_id, score, season_id, league_id, creator_id)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
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

	row := q.QueryRow(ctx, sql, model.CourtId, model.ScheduledAt, model.PlayerOneId, model.PlayerTwoId, model.WinnerId, model.Score, model.SeasonId, model.LeagueId, model.PlayerOneId)
	err := dest.ScanRow(row)
	if err != nil {
		return dest, err
	}

	return dest, nil
}

func (s *store) findMatches(ctx context.Context, seasonId, leagueId string) ([]MatchModel, error) {
	sql := `
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

	rows, err := s.db.Query(ctx, sql, seasonId, leagueId)
	if err != nil {
		return nil, fmt.Errorf("quering match rows: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mm MatchModel
		err := mm.ScanRows(rows)
		if err != nil {
			return nil, err
		}

		dest = append(dest, mm)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating match rows: %v", err)
	}

	return dest, nil
}

func (s *store) findMatch(ctx context.Context, seasonId, leagueId, matchId string) (*MatchModel, error) {
	sql := `
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

	row := s.db.QueryRow(ctx, sql, matchId, seasonId, leagueId)
	err := dest.ScanRow(row)
	if err != nil {
		return nil, err
	}

	return &dest, nil
}

func (s *store) updateMatch(ctx context.Context, tx pgx.Tx, model UpdateMatchRequestModel) (*MatchModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	sql := `
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

	row := q.QueryRow(ctx, sql, model.CourtId, model.ScheduledAt, model.PlayerTwoId, model.MatchId, model.SeasonId, model.LeagueId)
	err := dest.ScanRow(row)
	if err != nil {
		return nil, fmt.Errorf("updating match: %w", err)
	}

	return &dest, nil
}

func (s *store) updatePlayerStatistics(ctx context.Context, tx pgx.Tx, winnerId, playerOneId, playerTwoId string) error {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	// sql := `
	// 	update player
	// 	set
	// 		matches_played = matches_played + 1,
	// 		matches_won = matches_won + case when id = $1 then 1 else 0 end,
	// 		winning_ratio = case
	// 			when (matches_played + 1) > 0
	// 			then (matches_won + case when id = $1 then 1 else 0 end)::float / (matches_played + 1)
	// 			else 0
	// 		end
	// 	where id in ($2, $3)
	// `

	sql := `
	`

	_, err := q.Exec(ctx, sql, winnerId, playerOneId, playerTwoId)
	if err != nil {
		return fmt.Errorf("updating player stats: %v", err)
	}

	return nil
}

// todo: refactor be an atomic operation for one user, then call from service for each player
func (s *store) updateStandings(ctx context.Context, tx pgx.Tx, seasonId, leagueId, playerOneId, playerTwoId string, pl1Stats MatchStats, pl2Stats MatchStats) error {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	sql := `
		insert into standing (points, matches_played, matches_won, sets_won, sets_lost, games_won, games_lost, season_id, league_id, player_id)
		values ($1, 1, $2, $3, $4, $5, $6, $7, $8, $9)
		on conflict (season_id, league_id, player_id) do update
		set
			points = standing.points + $1,
			matches_played = standing.matches_played + 1,
			matches_won = standing.matches_won + $2,
			sets_won = standing.sets_won + $3,
			sets_lost = standing.sets_lost + $4,
			games_won = standing.games_won + $5,
			games_lost = standing.games_lost + $6
	`

	_, err := q.Exec(ctx, sql, pl1Stats.Pts, pl1Stats.WonMatches, pl1Stats.SetsWon, pl1Stats.SetsLost, pl1Stats.GamesWon, pl1Stats.GamesLost, seasonId, leagueId, playerOneId)
	if err != nil {
		return fmt.Errorf("updating player one standing: %v", err)
	}

	_, err = q.Exec(ctx, sql, pl2Stats.Pts, pl2Stats.WonMatches, pl2Stats.SetsWon, pl2Stats.SetsLost, pl2Stats.GamesWon, pl2Stats.GamesLost, seasonId, leagueId, playerTwoId)
	if err != nil {
		return fmt.Errorf("updating player two standing: %v", err)
	}

	return nil
}

func (s *store) updateMatchScore(ctx context.Context, tx pgx.Tx, seasonId, leagueId, matchId, score, winnerId string) (*MatchModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	sql := `
		with updated_match as (
			update match 
			set score = $1, winner_id = $2
			where id = $3 and season_id = $4 and league_id = $5
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

	row := q.QueryRow(ctx, sql, score, winnerId, matchId, seasonId, leagueId)
	err := dest.ScanRow(row)
	if err != nil {
		return nil, fmt.Errorf("updating match score: %w", err)
	}

	return &dest, nil
}

// helper
func (s *store) validateInsertUpdateMatch(ctx context.Context, courtId, seasonId, leagueId, player1Id, player2Id string) (courtExists bool, seasonExists bool, leagueExists bool, playerOneExists bool, playerTwoExists bool, playersInLeague bool, err error) {
	sql1 := `
		select
			exists (
				select 1 from court where id = $1
			) as court_exists,
			exists ( 
				select 1 from season where id = $2
			) as season_exists,
			exists (
				select 1 from league where id = $3 and season_id = $2
			) as league_exists,
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

	err = s.db.QueryRow(ctx, sql1, courtId, seasonId, leagueId, player1Id, player2Id).Scan(&courtExists, &seasonExists, &leagueExists, &playerOneExists, &playerTwoExists, &playersInLeague)
	if err != nil {
		return
	}

	return
}

// helper
func (s *store) validateFindMatches(ctx context.Context, seasonId, leagueId string) (seasonExists bool, leagueExists bool, err error) {
	sql1 := `
		select
			exists (
				select 1 from season where id = $1
			) as season_exists,
			exists (
				select 1 from league where id = $2 and season_id = $1
			) as league_exists
	`

	err = s.db.QueryRow(ctx, sql1, seasonId, leagueId).Scan(&seasonExists, &leagueExists)
	if err != nil {
		return
	}
	return
}

// helper
func (s *store) validateSubmitMatchScore(ctx context.Context, seasonId, leagueId, matchId string) (seasonExists bool, leagueExists bool, matchExists bool, err error) {
	sql1 := `
		select
			exists (
				select 1 from season where id = $1
			) as season_exists,
			exists (
				select 1 from league where id = $2 and season_id = $1
			) as league_exists,
			exists (
				select 1 from match where id = $3 and season_id = $1 and league_id = $2
			) as match_exists
	`

	err = s.db.QueryRow(ctx, sql1, seasonId, leagueId, matchId).Scan(&seasonExists, &leagueExists, &matchExists)
	if err != nil {
		return
	}
	return
}

// helper
func (s *store) checkMatchParticipation(ctx context.Context, matchId, playerId string) (bool, error) {
	sql := `
		select exists (
			select 1 from match
			where id = $1 and (player_one_id = $2 or player_two_id = $2)
		)
	`

	var exists bool
	err := s.db.QueryRow(ctx, sql, matchId, playerId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// helper
func (s *store) checkMatchOwnership(ctx context.Context, matchId, playerId string) (bool, error) {
	sql := `
		select exists (
			select 1 from match
			where id = $1 and creator_id = $2
		)
	`

	var exists bool
	err := s.db.QueryRow(ctx, sql, matchId, playerId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// helper
func (s *store) checkMatchScore(ctx context.Context, matchId string) (bool, error) {
	var score sql.NullString

	err := s.db.QueryRow(ctx, `
		select score
		from match
		where id = $1
	`, matchId).Scan(&score)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("finding match: %w", response.ErrNotFound)
		}
		return false, err
	}

	return score.Valid, nil
}
