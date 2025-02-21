package matches

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/failure"
)

type store struct {
	db *db.Conn
}

func newStore(db *db.Conn) *store {
	return &store{
		db,
	}
}

func (s *store) insertMatch(ctx context.Context, tx pgx.Tx, courtId, scheduledAt, playerOneId, playerTwoId string, winnerId, score *string, seasonId, leagueId string) (MatchModel, error) {
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

	row := q.QueryRow(ctx, sql, courtId, scheduledAt, playerOneId, playerTwoId, winnerId, score, seasonId, leagueId, playerOneId)
	err := dest.ScanRow(row)
	if err != nil {
		return dest, failure.New("unable to insert match", err)
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
		return nil, failure.New("unable to find matches", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	defer rows.Close()

	for rows.Next() {
		var mm MatchModel
		err := mm.ScanRows(rows)
		if err != nil {
			return nil, failure.New("unable to find matches", err)
		}

		dest = append(dest, mm)
	}

	if err := rows.Err(); err != nil {
		return nil, failure.New("unable to find matches", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
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
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("match not found", err)
		}
		return nil, failure.New("unable to find match", err)
	}

	return &dest, nil
}

func (s *store) updateMatch(ctx context.Context, tx pgx.Tx, courtId, scheduledAt, playerTwoId, seasonId, leagueId, matchId string) (*MatchModel, error) {
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

	row := q.QueryRow(ctx, sql, courtId, scheduledAt, playerTwoId, matchId, seasonId, leagueId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("match for update not found", err)
		}
		return nil, failure.New("unable to update match", err)
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

	sql := `
		update player
		set 
			matches_played = matches_played + 1,
			matches_won = matches_won + case when id = $1 then 1 else 0 end
		where id in ($2, $3)
	`

	_, err := q.Exec(ctx, sql, winnerId, playerOneId, playerTwoId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return failure.New("players for updating statistics not found", fmt.Errorf("%w -> %v", failure.ErrNotFound, err))
		}
		return failure.New("unable to update player statistics", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return nil
}

func (s *store) incrementPlayerMatchesScheduled(ctx context.Context, tx pgx.Tx, playerId string) error {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	sql := `
		update player
		set
			matches_scheduled = matches_scheduled + 1
		where id = $1	
	`

	_, err := q.Exec(ctx, sql, playerId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return failure.New("player for updating matches scheduled not found", fmt.Errorf("%w -> %v", failure.ErrNotFound, err))
		}
		return failure.New("unable to increment player matches scheduled", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return nil
}

func (s *store) updateStanding(ctx context.Context, tx pgx.Tx, seasonId, leagueId, playerId string, plStats MatchStats) error {
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

	_, err := q.Exec(ctx, sql, plStats.Pts, plStats.WonMatches, plStats.SetsWon, plStats.SetsLost, plStats.GamesWon, plStats.GamesLost, seasonId, leagueId, playerId)
	if err != nil {
		return failure.New("unable to update standings", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
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
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("match for updating score not found", err)
		}
		return nil, failure.New("unable to update match score", err)
	}

	return &dest, nil
}

// helper - is the player part of a match
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
		return false, failure.New("unable to check if player is a match participant", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return exists, nil
}

// helper - did the player create the match
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
		return false, failure.New("unable to check match ownership", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	return exists, nil
}

// helper - check if the match score has been entered (not null)
func (s *store) checkMatchScore(ctx context.Context, matchId string) (bool, error) {
	var score sql.NullString

	err := s.db.QueryRow(ctx, `select score from match where id = $1`, matchId).Scan(&score)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, failure.New("checking match score - match not found", fmt.Errorf("%w -> %v", failure.ErrNotFound, err))
		}
		return false, failure.New("unable to find match score", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return score.Valid, nil
}
