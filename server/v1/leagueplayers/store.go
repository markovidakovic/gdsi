package leagueplayers

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/params"
	"github.com/markovidakovic/gdsi/server/v1/players"
)

//go:embed queries/*.sql
var sqlFiles embed.FS

type store struct {
	db      *db.Conn
	queries struct {
		list                   string
		findById               string
		update                 string
		incrementSeasonsPlayed string
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
	listBytes, err := sqlFiles.ReadFile("queries/list.sql")
	if err != nil {
		return fmt.Errorf("failed to read list.sql -> %v", err)
	}
	findByIdBytes, err := sqlFiles.ReadFile("queries/find_by_id.sql")
	if err != nil {
		return fmt.Errorf("failed to read find_by_id.sql -> %v", err)
	}
	updateBytes, err := sqlFiles.ReadFile("queries/update_current_league.sql")
	if err != nil {
		return fmt.Errorf("failed to read update.sql -> %v", err)
	}
	incrementSeasonsPlayedBytes, err := sqlFiles.ReadFile("queries/increment_seasons_played.sql")
	if err != nil {
		return fmt.Errorf("failed to read increment_seasons_played.sql -> %v", err)
	}

	s.queries.list = string(listBytes)
	s.queries.findById = string(findByIdBytes)
	s.queries.update = string(updateBytes)
	s.queries.incrementSeasonsPlayed = string(incrementSeasonsPlayedBytes)

	return nil
}

var allowedSortFeilds = map[string]string{
	"created_at": "player.created_at",
}

func (s *store) findLeaguePlayers(ctx context.Context, leagueId string, requestingPlayerId string, matchAvailable bool, limit, offset int, sort *params.OrderBy) ([]players.PlayerModel, error) {
	args := []interface{}{leagueId}
	argCounter := 2 // starting with $2 since $1 is already used for current_league_id

	if matchAvailable {
		// exclude the player requesting
		s.queries.list += fmt.Sprintf(" and player.id != $%d", argCounter)
		args = append(args, requestingPlayerId)
		argCounter++

		// exclude players who have already played a match with the requesting player
		s.queries.list += fmt.Sprintf(`
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
		s.queries.list += fmt.Sprintf("order by %s %s\n", allowedSortFeilds[sort.Field], sort.Direction)
	} else {
		s.queries.list += fmt.Sprintln("order by player.created_at desc")
	}

	if limit >= 0 {
		s.queries.list += fmt.Sprintf("limit $%d offset $%d", argCounter, argCounter+1)
		args = append(args, limit, offset)
	}

	rows, err := s.db.Query(ctx, s.queries.list, args...)
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
	row := s.db.QueryRow(ctx, s.queries.findById, playerId, leagueId)
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

	var dest players.PlayerModel
	row := q.QueryRow(ctx, s.queries.update, leagueId, playerId)
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

	var dest players.PlayerModel
	row := q.QueryRow(ctx, s.queries.incrementSeasonsPlayed, playerId, leagueId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return dest, failure.New("league player for incrementing seasons played not found", err)
		}
		return dest, failure.New("unable to increment league player seasons played", err)
	}

	return dest, nil
}
