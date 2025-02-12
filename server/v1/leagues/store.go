package leagues

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
		return dest, fmt.Errorf("inserting league: %w", err)
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
		return nil, fmt.Errorf("quering leagues: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var lm LeagueModel
		err := lm.ScanRows(rows)
		if err != nil {
			return nil, err
		}

		dest = append(dest, lm)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("scanning league rows: %v", err)
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
		return nil, err
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
		return dest, fmt.Errorf("updating league: %w", err)
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
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("deleting league: %w", response.ErrNotFound)
	}

	return nil
}

func (s *store) checkSeasonExistance(ctx context.Context, seasonId string) (bool, error) {
	sql := `
		select exists (
			select 1 from season where id = $1
		)
	`

	var exists bool
	err := s.db.QueryRow(ctx, sql, seasonId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
