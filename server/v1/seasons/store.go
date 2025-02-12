package seasons

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

func (s *store) insertSeason(ctx context.Context, tx pgx.Tx, model CreateSeasonRequestModel) (SeasonModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	sql := `
		with inserted_season as (
			insert into season (title, description, start_date, end_date, creator_id)
			values ($1, $2, $3, $4, $5)
			returning id, title, description, start_date, end_date, creator_id, created_at
		)
		select s.id, s.title, s.description, s.start_date, s.end_date, account.id as creator_id, account.name as creator_name, s.created_at
		from inserted_season s
		join account on s.creator_id = account.id
	`

	var dest SeasonModel

	row := q.QueryRow(ctx, sql, model.Title, model.Description, model.StartDate, model.EndDate, model.CreatorId)
	err := dest.ScanRow(row)
	if err != nil {
		return dest, fmt.Errorf("inserting season: %w", err)
	}

	return dest, nil
}

func (s *store) findSeasons(ctx context.Context) ([]SeasonModel, error) {
	sql := `
		select
			season.id,
			season.title,
			season.description,
			season.start_date,
			season.end_date,
			account.id as creator_id,
			account.name as creator_name,
			season.created_at
		from season
		join account on season.creator_id = account.id
		order by season.created_at desc
	`

	rows, err := s.db.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("quering seasons: %v", err)
	}
	defer rows.Close()

	dest := []SeasonModel{}
	for rows.Next() {
		var sm SeasonModel

		err := sm.ScanRows(rows)
		if err != nil {
			return nil, err
		}

		dest = append(dest, sm)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scanning season rows: %v", err)
	}

	return dest, nil
}

func (s *store) findSeason(ctx context.Context, seasonId string) (*SeasonModel, error) {
	sql := `
		select
			season.id,
			season.title,
			season.description,
			season.start_date,
			season.end_date,
			account.id as creator_id,
			account.name as creator_name,
			season.created_at
		from season
		join account on season.creator_id = account.id
		where season.id = $1
	`

	var dest SeasonModel

	row := s.db.QueryRow(ctx, sql, seasonId)
	err := dest.ScanRow(row)
	if err != nil {
		return nil, fmt.Errorf("finding season: %w", err)
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

	sql := `
		with updated_season as (
			update season 
			set title = $1, description = $2, start_date = $3, end_date = $4
			where id = $5
			returning id, title, description, start_date, end_date, creator_id, created_at
		)
		select 
			us.id as season_id,
			us.title as season_title,
			us.description as season_description,
			us.start_date as season_start_date,
			us.end_date as season_end_date,
			account.id as creator_id,
			account.name as creator_name,
			us.created_at as season_created_at
		from updated_season us
		join account on us.creator_id = account.id
	`

	var dest SeasonModel

	row := q.QueryRow(ctx, sql, model.Title, model.Description, model.StartDate, model.EndDate, seasonId)
	err := dest.ScanRow(row)
	if err != nil {
		return nil, fmt.Errorf("updating season: %w", err)
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

	sql := `
		delete from season where id = $1
	`

	ct, err := q.Exec(ctx, sql, seasonId)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return response.ErrNotFound
	}

	return nil
}
