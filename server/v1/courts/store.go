package courts

import (
	"context"
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

func (s *store) insertCourt(ctx context.Context, tx pgx.Tx, name, creatorId string) (CourtModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	sql := `
		with inserted_court as (
			insert into court (name, creator_id)
			values ($1, $2)
			returning id, name, creator_id, created_at			
		)
		select 
			ic.id as court_id, 
			ic.name as court_name, 
			account.id as creator_id, 
			account.name as creator_name, 
			ic.created_at as court_created_at
		from inserted_court ic
		join account on ic.creator_id = account.id
	`

	var dest CourtModel
	row := q.QueryRow(ctx, sql, name, creatorId)
	err := dest.ScanRow(row)
	if err != nil {
		return dest, failure.New("failed to insert court", err)
	}

	return dest, nil
}

func (s *store) findCourts(ctx context.Context, limit, offset int) ([]CourtModel, error) {
	sql := `
		select 
			court.id as court_id,
			court.name as court_name,
			account.id as creator_id,
			account.name as creator_name,
			court.created_at as court_created_at
		from court
		join account on court.creator_id = account.id
		order by court.created_at desc		
	`

	var err error
	var rows pgx.Rows
	if limit >= 0 {
		sql += `limit $1 offset $2`
		rows, err = s.db.Query(ctx, sql, limit, offset)
		if err != nil {
			return nil, failure.New("unable to find courts", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
		}
	} else {
		rows, err = s.db.Query(ctx, sql)
		if err != nil {
			return nil, failure.New("unable to find courts", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
		}
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

	sql := `
		select 
			court.id as court_id,
			court.name as court_name,
			account.id as creator_id,
			account.name as creator_name,
			court.created_at as court_created_at
		from court
		join account on court.creator_id = account.id
		where court.id = $1
	`

	row := s.db.QueryRow(ctx, sql, courtId)
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

	sql := `
		with updated_court as (
			update court
			set name = $1
			where id = $2
			returning id, name, creator_id, created_at
		)
		select
			uc.id as court_id,
			uc.name as court_name,
			account.id as creator_id,
			account.name as creator_name,
			uc.created_at as court_created_at
		from updated_court uc
		join account on uc.creator_id = account.id
	`

	row := q.QueryRow(ctx, sql, name, courtId)
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

	sql := `
		delete from court where id = $1
	`

	ct, err := q.Exec(ctx, sql, courtId)
	if err != nil {
		return failure.New("unable to delete court", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	if ct.RowsAffected() == 0 {
		return failure.New("court not found", failure.ErrNotFound)
	}

	return nil
}
