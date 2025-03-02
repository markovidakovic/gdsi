package seasons

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/params"
)

type store struct {
	db *db.Conn
}

func newStore(db *db.Conn) *store {
	return &store{
		db,
	}
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
		return dest, failure.New("failed to insert season", err)
	}

	return dest, nil
}

func (s *store) findSeasons(ctx context.Context, limit, offset int, sort *params.OrderBy) ([]SeasonModel, error) {
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
	`

	if sort != nil && sort.IsValid(allowedSortFields) {
		sql += fmt.Sprintf("order by %s %s\n", allowedSortFields[sort.Field], sort.Direction)
	} else {
		sql += fmt.Sprintln("order by season.created_at desc")
	}

	var err error
	var rows pgx.Rows
	if limit >= 0 {
		sql += `limit $1 offset $2`
		rows, err = s.db.Query(ctx, sql, limit, offset)
	} else {
		rows, err = s.db.Query(ctx, sql)
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

	sql := `
		delete from season where id = $1
	`

	ct, err := q.Exec(ctx, sql, seasonId)
	if err != nil {
		return failure.New("unable to delete season", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	if ct.RowsAffected() == 0 {
		return failure.New("season for deletion not found", failure.ErrNotFound)
	}

	return nil
}
