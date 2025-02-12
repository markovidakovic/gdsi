package courts

import (
	"context"
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

func (s *store) insertCourt(ctx context.Context, tx pgx.Tx, name, creatorId string) (CourtModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	// cte - common table expression
	sql := `
		with inserted_court as (
			insert into court (name, creator_id)
			values ($1, $2)
			returning id, name, creator_id, created_at			
		)
		select court.id, court.name, court.created_at, account.id as creator_id, account.name as creator_name
		from inserted_court court
		join account on court.creator_id = account.id
	`

	var dest CourtModel
	row := q.QueryRow(ctx, sql, creatorId)
	err := dest.ScanRow(row)
	if err != nil {
		return dest, fmt.Errorf("inserting court: %v", err)
	}

	return dest, nil
}

func (s *store) findCourts(ctx context.Context) ([]CourtModel, error) {
	sql := `
		select 
			court.id as court_id,
			court.name as court_name,
			court.created_at as court_created_at,
			account.id as creator_id,
			account.name as creator_name
		from court
		join account on court.creator_id = account.id
		order by court.created_at desc		
	`

	rows, err := s.db.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("querying courts: %v", err)
	}
	defer rows.Close()

	var dest = []CourtModel{}
	for rows.Next() {
		var cm CourtModel
		err := cm.ScanRows(rows)
		if err != nil {
			return nil, err
		}
		dest = append(dest, cm)
	}

	if err = rows.Err(); err != nil {
		// todo: what is the appropriate error msg here?
		return nil, fmt.Errorf("scanning court rows: %v", err)
	}

	return dest, nil
}

func (s *store) findCourt(ctx context.Context, courtId string) (*CourtModel, error) {
	var dest CourtModel

	sql := `
		select 
			court.id as court_id,
			court.name as court_name,
			court.created_at as court_created_at,
			account.id as creator_id,
			account.name as creator_name
		from court
		join account on court.creator_id = account.id
		where court.id = $1
	`

	row := s.db.QueryRow(ctx, sql, courtId)
	err := dest.ScanRow(row)
	if err != nil {
		// todo: maybe wrap this err with an additional msg?
		return nil, err
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
			court.id as court_id,
			court.name as court_name,
			account.id as creator_id,
			account.name as creator_name,
			court.created_at as court_created_at
		from updated_court court
		join account on court.creator_id = account.id
	`

	row := q.QueryRow(ctx, sql, name, courtId)
	err := dest.ScanRow(row)
	if err != nil {
		return nil, fmt.Errorf("updating court: %w", err)
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
		return err
	}
	if ct.RowsAffected() == 0 {
		return response.ErrNotFound
	}

	return nil
}
