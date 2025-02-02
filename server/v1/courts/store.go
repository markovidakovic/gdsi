package courts

import (
	"context"
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

func (s *store) insertCourt(ctx context.Context, input CreateCourtRequestModel) (CourtModel, error) {
	// cte - common table expression
	sql1 := `
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

	err := s.db.QueryRow(ctx, sql1, input.Name, input.CreatorId).Scan(&dest.Id, &dest.Name, &dest.CreatedAt, &dest.Creator.Id, &dest.Creator.Name)
	if err != nil {
		return dest, fmt.Errorf("inserting court: %v", err)
	}

	return dest, nil
}

func (s *store) findCourts(ctx context.Context) ([]CourtModel, error) {
	sql1 := `
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

	rows, err := s.db.Query(ctx, sql1)
	if err != nil {
		return nil, fmt.Errorf("querying courts: %v", err)
	}
	defer rows.Close()

	var dest = []CourtModel{}
	for rows.Next() {
		var cm CourtModel
		err := rows.Scan(&cm.Id, &cm.Name, &cm.CreatedAt, &cm.Creator.Id, &cm.Creator.Name)
		if err != nil {
			return nil, fmt.Errorf("scanning court row: %v", err)
		}
		dest = append(dest, cm)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("scanning court rows: %v", err)
	}

	return dest, nil
}

func (s *store) findCourt(ctx context.Context, courtId string) (*CourtModel, error) {
	var dest CourtModel

	sql1 := `
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

	err := s.db.QueryRow(ctx, sql1, courtId).Scan(&dest.Id, &dest.Name, &dest.CreatedAt, &dest.Creator.Id, &dest.Creator.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, response.ErrNotFound
		}
		return nil, err
	}

	return &dest, nil
}

func (s *store) updateCourt(ctx context.Context, courtId string, input UpdateCourtRequestModel) (*CourtModel, error) {
	var dest CourtModel

	sql1 := `
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

	err := s.db.QueryRow(ctx, sql1, input.Name, courtId).Scan(&dest.Id, &dest.Name, &dest.Creator.Id, &dest.Creator.Name, &dest.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, response.ErrNotFound
		}
		return nil, err
	}

	return &dest, nil
}

func (s *store) deleteCourt(ctx context.Context, courtId string) error {
	sql1 := `
		delete from court where id = $1
	`

	ct, err := s.db.Exec(ctx, sql1, courtId)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return response.ErrNotFound
	}

	return nil
}
