package seasons

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
		return fmt.Errorf("failed to read insert.sql file -> %v", err)
	}
	listBytes, err := sqlFiles.ReadFile("queries/list.sql")
	if err != nil {
		return fmt.Errorf("failed to read list.sql file -> %v", err)
	}
	findByIdBytes, err := sqlFiles.ReadFile("queries/find_by_id.sql")
	if err != nil {
		return fmt.Errorf("failed to read find_by_id.sql file -> %v", err)
	}
	updateBytes, err := sqlFiles.ReadFile("queries/update.sql")
	if err != nil {
		return fmt.Errorf("failed to read update.sql file -> %v", err)
	}
	deleteBytes, err := sqlFiles.ReadFile("queries/delete.sql")
	if err != nil {
		return fmt.Errorf("failed to read delete.sql file -> %v", err)
	}

	s.queries.insert = string(insertBytes)
	s.queries.list = string(listBytes)
	s.queries.findById = string(findByIdBytes)
	s.queries.update = string(updateBytes)
	s.queries.delete = string(deleteBytes)

	return nil
}

var allowedSortFields = map[string]string{
	"created_at": "season.created_at",
}

func (s *store) insertSeason(ctx context.Context, tx pgx.Tx, model CreateSeasonRequestModel) (SeasonModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	var dest SeasonModel
	row := q.QueryRow(ctx, s.queries.insert, model.Title, model.Description, model.StartDate, model.EndDate, model.CreatorId)
	err := dest.ScanRow(row)
	if err != nil {
		return dest, failure.New("failed to insert season", err)
	}

	return dest, nil
}

func (s *store) findSeasons(ctx context.Context, limit, offset int, sort *params.OrderBy) ([]SeasonModel, error) {
	if sort != nil && sort.IsValid(allowedSortFields) {
		s.queries.list += fmt.Sprintf("order by %s %s\n", allowedSortFields[sort.Field], sort.Direction)
	} else {
		s.queries.list += fmt.Sprintln("order by season.created_at desc")
	}

	var err error
	var rows pgx.Rows
	if limit >= 0 {
		s.queries.list += `limit $1 offset $2`
		rows, err = s.db.Query(ctx, s.queries.list, limit, offset)
	} else {
		rows, err = s.db.Query(ctx, s.queries.list)
	}

	if err != nil {
		return nil, failure.New("unable to find seasons", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	defer rows.Close()

	dest := []SeasonModel{}
	for rows.Next() {
		var sm SeasonModel

		err := sm.ScanRows(rows)
		if err != nil {
			return nil, failure.New("unable to find seasons", err)
		}

		dest = append(dest, sm)
	}

	if err := rows.Err(); err != nil {
		return nil, failure.New("unable to find seasons", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return dest, nil
}

func (s *store) countSeasons(ctx context.Context) (int, error) {
	var count int
	sql := `select count(*) from season`
	err := s.db.QueryRow(ctx, sql).Scan(&count)
	if err != nil {
		return 0, failure.New("unable to count seasons", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	return count, nil
}

func (s *store) findSeason(ctx context.Context, seasonId string) (*SeasonModel, error) {
	var dest SeasonModel
	row := s.db.QueryRow(ctx, s.queries.findById, seasonId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("season not found", err)
		}
		return nil, failure.New("unable to find season", err)
	}

	return &dest, nil
}

func (s *store) updateSeason(ctx context.Context, tx pgx.Tx, seasonId string, model UpdateSeasonRequestModel) (*SeasonModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	var dest SeasonModel
	row := q.QueryRow(ctx, s.queries.update, model.Title, model.Description, model.StartDate, model.EndDate, seasonId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("season for update not found", err)
		}
		return nil, failure.New("unable to update season", err)
	}

	return &dest, nil
}

func (s *store) deleteSeason(ctx context.Context, tx pgx.Tx, seasonId string) error {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	ct, err := q.Exec(ctx, s.queries.delete, seasonId)
	if err != nil {
		return failure.New("unable to delete season", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	if ct.RowsAffected() == 0 {
		return failure.New("season for deletion not found", failure.ErrNotFound)
	}

	return nil
}
