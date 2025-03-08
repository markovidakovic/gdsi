package leagues

import (
	"context"
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
		insert   string
		list     string
		findById string
		update   string
		delete   string
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
	deleteBytes, err := sqlFiles.ReadFile("queries/delete.sql")
	if err != nil {
		return fmt.Errorf("failed to read delete.sql -> %v", err)
	}

	s.queries.insert = string(insertBytes)
	s.queries.list = string(listBytes)
	s.queries.findById = string(findByIdBytes)
	s.queries.update = string(updateBytes)
	s.queries.delete = string(deleteBytes)

	return nil
}

var allowedSortFields = map[string]string{
	"title":      "league.title",
	"start_date": "league.start_date",
	"end_date":   "league.end_date",
	"created_at": "league.created_at",
}

func (s *store) insertLeague(ctx context.Context, tx pgx.Tx, title string, description *string, creatorId string, seasonId string) (LeagueModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	var dest LeagueModel
	row := q.QueryRow(ctx, s.queries.insert, title, description, seasonId, creatorId)
	err := dest.ScanRow(row)
	if err != nil {
		return dest, failure.New("failed to insert league", err)
	}

	return dest, nil
}

func (s *store) findLeagues(ctx context.Context, seasonId string, limit, offset int, sort *params.OrderBy) ([]LeagueModel, error) {
	if sort != nil && sort.IsValid(allowedSortFields) {
		s.queries.list += fmt.Sprintf("order by %s %s\n", allowedSortFields[sort.Field], sort.Direction)
	} else {
		s.queries.list += fmt.Sprintln("order by league.created_at desc")
	}

	var err error
	var rows pgx.Rows
	if limit >= 0 {
		s.queries.list += `limit $2 offset $3`
		rows, err = s.db.Query(ctx, s.queries.list, seasonId, limit, offset)
	} else {
		rows, err = s.db.Query(ctx, s.queries.list, seasonId)
	}

	if err != nil {
		return nil, failure.New("unable to find leagues", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	defer rows.Close()

	dest := []LeagueModel{}
	for rows.Next() {
		var lm LeagueModel
		err := lm.ScanRows(rows)
		if err != nil {
			return nil, failure.New("unable to find leagues", err)
		}

		dest = append(dest, lm)
	}

	if err = rows.Err(); err != nil {
		return nil, failure.New("unable to find leagues", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return dest, nil
}

func (s *store) countLeagues(ctx context.Context, seasonId string) (int, error) {
	var count int
	sql := `select count(*) from league where league.season_id = $1`
	err := s.db.QueryRow(ctx, sql, seasonId).Scan(&count)
	if err != nil {
		return 0, failure.New("unable to count leagues", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	return count, nil
}

func (s *store) findLeague(ctx context.Context, seasonId, leagueId string) (*LeagueModel, error) {
	var dest LeagueModel
	row := s.db.QueryRow(ctx, s.queries.findById, seasonId, leagueId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("league not found", err)
		}
		return nil, failure.New("unable to find league", err)
	}

	return &dest, nil
}

func (s *store) updateLeague(ctx context.Context, tx pgx.Tx, title string, description *string, seasonId, leagueId string) (LeagueModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	var dest LeagueModel
	row := q.QueryRow(ctx, s.queries.update, title, description, leagueId, seasonId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return dest, failure.New("league not found", err)
		}
		return dest, failure.New("unable to update league", err)
	}

	return dest, nil
}

func (s *store) deleteLeague(ctx context.Context, tx pgx.Tx, seasonId, leagueId string) error {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	ct, err := q.Exec(ctx, s.queries.delete, leagueId, seasonId)
	if err != nil {
		return failure.New("unable to delete league", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	if ct.RowsAffected() == 0 {
		return failure.New("league not found", failure.ErrNotFound)
	}

	return nil
}
