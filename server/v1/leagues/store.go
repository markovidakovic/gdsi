package leagues

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

func (s *store) insertLeague(ctx context.Context, tx pgx.Tx, title string, description *string, creatorId string, seasonId string) (LeagueModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	sql := `
		with inserted_league as (
			insert into league (title, description, season_id, creator_id)
			values ($1, $2, $3, $4)
			returning id, title, description, season_id, creator_id, created_at
		)
		select
			il.id,
			il.title,
			il.description,
			season.id as season_id,
			season.title as season_title,
			account.id as creator_id,
			account.name as creator_name,
			il.created_at
		from inserted_league il
		join season on il.season_id = season.id
		join account on il.creator_id = account.id
	`

	var dest LeagueModel
	row := q.QueryRow(ctx, sql, title, description, seasonId, creatorId)
	err := dest.ScanRow(row)
	if err != nil {
		return dest, failure.New("failed to insert league", err)
	}

	return dest, nil
}

func (s *store) findLeagues(ctx context.Context, seasonId string) ([]LeagueModel, error) {
	sql := `
		select 
			league.id,
			league.title,
			league.description,
			season.id as season_id,
			season.title as season_title,
			account.id as creator_id,
			account.name as creator_name,
			league.created_at
		from league
		join season on league.season_id = season.id
		join account on league.creator_id = account.id
		where league.season_id = $1
		order by league.created_at desc
	`

	dest := []LeagueModel{}

	rows, err := s.db.Query(ctx, sql, seasonId)
	if err != nil {
		return nil, failure.New("unable to find leagues", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	defer rows.Close()

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

func (s *store) findLeague(ctx context.Context, seasonId, leagueId string) (*LeagueModel, error) {
	sql := `
		select 
			league.id,
			league.title,
			league.description,
			season.id as season_id,
			season.title as season_title,
			account.id as creator_id,
			account.name as creator_name,
			league.created_at
		from league
		join season on league.season_id = season.id
		join account on league.creator_id = account.id
		where league.season_id = $1 and league.id = $2
	`

	var dest LeagueModel
	row := s.db.QueryRow(ctx, sql, seasonId, leagueId)
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

	sql := `
		with updated_league as (
			update league
			set title = $1, description = $2
			where id = $3 and season_id = $4
			returning id, title, description, season_id, creator_id, created_at
		)
		select
			ul.id,
			ul.title,
			ul.description,
			season.id as season_id,
			season.title as season_title,
			account.id as creator_id,
			account.name as creator_name,
			ul.created_at
		from updated_league ul
		join season on ul.season_id = season.id
		join account on ul.creator_id = account.id
	`

	var dest LeagueModel
	row := q.QueryRow(ctx, sql, title, description, leagueId, seasonId)
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

	sql := `
		delete from league where id = $1 and season_id = $2
	`

	ct, err := q.Exec(ctx, sql, leagueId, seasonId)
	if err != nil {
		return failure.New("unable to delete league", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	if ct.RowsAffected() == 0 {
		return failure.New("league not found", failure.ErrNotFound)
	}

	return nil
}
