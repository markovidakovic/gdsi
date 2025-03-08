package matches

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/params"
)

//go:embed queries/*.sql
var sqlFiles embed.FS

type store struct {
	db      *db.Conn
	queries struct {
		insert                          string
		list                            string
		findById                        string
		update                          string
		updatePlayerStatistics          string
		incrementPlayerMatchesScheduled string
		updateStanding                  string
		updateScore                     string
	}
}

func newStore(db *db.Conn) (*store, error) {
	s := &store{
		db: db,
	}
	if err := s.loadQueries(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *store) loadQueries() error {
	insertBytes, err := sqlFiles.ReadFile("queries/insert.sql")
	if err != nil {
		return fmt.Errorf("failed to read insert.sql -> %v", err)
	}
	listBytes, err := sqlFiles.ReadFile("queries/list.sql")
	if err != nil {
		return fmt.Errorf("failed to read list.sql -> %v", err)
	}
	findByIdBytes, err := sqlFiles.ReadFile("queries/find_by_id.sql")
	if err != nil {
		return fmt.Errorf("failed to read find_by_id.sql -> %v", err)
	}
	updateBytes, err := sqlFiles.ReadFile("queries/update.sql")
	if err != nil {
		return fmt.Errorf("failed to read update.sql -> %v", err)
	}
	updatePlayerStatisticsBytes, err := sqlFiles.ReadFile("queries/update_player_statistics.sql")
	if err != nil {
		return fmt.Errorf("failed to read update_player_statistics.sql -> %v", err)
	}
	incrementPlayerMatchesScheduledBytes, err := sqlFiles.ReadFile("queries/increment_player_matches_scheduled.sql")
	if err != nil {
		return fmt.Errorf("failed to read increment_player_matches_scheduled.sql -> %v", err)
	}
	updateStandingBytes, err := sqlFiles.ReadFile("queries/update_standing.sql")
	if err != nil {
		return fmt.Errorf("failed to read update_standing.sql -> %v", err)
	}
	updateScoreBytes, err := sqlFiles.ReadFile("queries/update_score.sql")
	if err != nil {
		return fmt.Errorf("failed to read update_score.sql -> %v", err)
	}

	s.queries.insert = string(insertBytes)
	s.queries.list = string(listBytes)
	s.queries.findById = string(findByIdBytes)
	s.queries.update = string(updateBytes)
	s.queries.updatePlayerStatistics = string(updatePlayerStatisticsBytes)
	s.queries.incrementPlayerMatchesScheduled = string(incrementPlayerMatchesScheduledBytes)
	s.queries.updateStanding = string(updateStandingBytes)
	s.queries.updateScore = string(updateScoreBytes)

	return nil
}

var allowedSortFields = map[string]string{
	"created_at": "match.created_at",
}

func (s *store) insertMatch(ctx context.Context, tx pgx.Tx, courtId, scheduledAt, playerOneId, playerTwoId string, winnerId, score *string, seasonId, leagueId string) (MatchModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	var dest MatchModel

	row := q.QueryRow(ctx, s.queries.insert, courtId, scheduledAt, playerOneId, playerTwoId, winnerId, score, seasonId, leagueId, playerOneId)
	err := dest.ScanRow(row)
	if err != nil {
		return dest, failure.New("unable to insert match", err)
	}

	return dest, nil
}

func (s *store) findMatches(ctx context.Context, seasonId, leagueId string, limit, offset int, sort *params.OrderBy) ([]MatchModel, error) {
	if sort != nil && sort.IsValid(allowedSortFields) {
		s.queries.list += fmt.Sprintf("order by %s %s\n", allowedSortFields[sort.Field], sort.Direction)
	} else {
		s.queries.list += fmt.Sprintln("order by match.created_at desc")
	}

	var err error
	var rows pgx.Rows
	if limit >= 0 {
		s.queries.list += `limit $3 offset $4`
		rows, err = s.db.Query(ctx, s.queries.list, seasonId, leagueId, limit, offset)
	} else {
		rows, err = s.db.Query(ctx, s.queries.list, seasonId, leagueId)
	}

	if err != nil {
		return nil, failure.New("unable to find matches", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	defer rows.Close()

	dest := []MatchModel{}
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

func (s *store) countMatches(ctx context.Context, seasonId, leagueId string) (int, error) {
	var count int
	sql := `select count(*) from match where season_id = $1 and league_id = $2`
	err := s.db.QueryRow(ctx, sql, seasonId, leagueId).Scan(&count)
	if err != nil {
		return 0, failure.New("unable to count matches", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	return count, nil
}

func (s *store) findMatch(ctx context.Context, seasonId, leagueId, matchId string) (*MatchModel, error) {
	var dest MatchModel
	row := s.db.QueryRow(ctx, s.queries.findById, matchId, seasonId, leagueId)
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

	var dest MatchModel
	row := q.QueryRow(ctx, s.queries.update, courtId, scheduledAt, playerTwoId, matchId, seasonId, leagueId)
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

	_, err := q.Exec(ctx, s.queries.updatePlayerStatistics, winnerId, playerOneId, playerTwoId)
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

	_, err := q.Exec(ctx, s.queries.incrementPlayerMatchesScheduled, playerId)
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

	_, err := q.Exec(ctx, s.queries.updateStanding, plStats.Pts, plStats.WonMatches, plStats.SetsWon, plStats.SetsLost, plStats.GamesWon, plStats.GamesLost, seasonId, leagueId, playerId)
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

	var dest MatchModel

	row := q.QueryRow(ctx, s.queries.updateScore, score, winnerId, matchId, seasonId, leagueId)
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
