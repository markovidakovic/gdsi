package courts

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

var sortingFields = map[string]string{
	"name":       "court.name",
	"created_at": "court.created_at",
}

func (s *store) insertCourt(ctx context.Context, tx pgx.Tx, name, creatorId string) (CourtModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	var dest CourtModel
	row := q.QueryRow(ctx, s.queries.insert, name, creatorId)
	err := dest.ScanRow(row)
	if err != nil {
		return dest, failure.New("failed to insert court", err)
	}

	return dest, nil
}

func (s *store) findCourts(ctx context.Context, limit, offset int, orderBy *params.OrderBy) ([]CourtModel, error) {
	if orderBy != nil && orderBy.IsValid(sortingFields) {
		s.queries.list += fmt.Sprintf("order by %s %s\n", sortingFields[orderBy.Field], orderBy.Direction)
	} else {
		s.queries.list += fmt.Sprintln("order by court.created_at desc")
	}

	var err error
	var rows pgx.Rows
	if limit >= 0 {
		s.queries.list += "limit $1 offset $2"
		rows, err = s.db.Query(ctx, s.queries.list, limit, offset)
	} else {
		rows, err = s.db.Query(ctx, s.queries.list)
	}

	if err != nil {
		return nil, failure.New("unable to find courts", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	defer rows.Close()

	var dest = []CourtModel{}
	for rows.Next() {
		var cm CourtModel
		err := cm.ScanRows(rows)
		if err != nil {
			return nil, failure.New("unable to find courts", err)
		}
		dest = append(dest, cm)
	}

	if err = rows.Err(); err != nil {
		return nil, failure.New("unable to find courts", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return dest, nil
}

func (s *store) countCourts(ctx context.Context) (int, error) {
	var count int
	sql := `select count(*) from court`
	err := s.db.QueryRow(ctx, sql).Scan(&count)
	if err != nil {
		return 0, failure.New("unable to count courts", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	return count, nil
}

func (s *store) findCourt(ctx context.Context, courtId string) (*CourtModel, error) {
	var dest CourtModel
	row := s.db.QueryRow(ctx, s.queries.findById, courtId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("court not found", err)
		}
		return nil, failure.New("unable to find court", err)
	}

	return &dest, nil
}

func (s *store) updateCourt(ctx context.Context, tx pgx.Tx, courtId string, name string) (*CourtModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	var dest CourtModel
	row := q.QueryRow(ctx, s.queries.update, name, courtId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("court not found", err)
		}
		return nil, failure.New("unable to update court", err)
	}

	return &dest, nil
}

func (s *store) deleteCourt(ctx context.Context, tx pgx.Tx, courtId string) error {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	ct, err := q.Exec(ctx, s.queries.delete, courtId)
	if err != nil {
		return failure.New("unable to delete court", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	if ct.RowsAffected() == 0 {
		return failure.New("court not found", failure.ErrNotFound)
	}

	return nil
}
