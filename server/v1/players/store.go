package players

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

func (s *store) findPlayers(ctx context.Context) ([]PlayerModel, error) {
	sql := `
		select
			player.id,
			player.height,
			player.weight,
			player.handedness,
			player.racket,
			player.matches_expected,
			player.matches_played,
			player.matches_won,
			player.matches_scheduled,
			player.seasons_played,
			account.id as account_id,
			account.name as account_name,
			league.id as league_id,
			league.title as league_title,
			player.created_at
		from player
		join account on player.account_id = account.id
		left join league on player.current_league_id = league.id
		order by player.created_at desc
	`

	var dest = []PlayerModel{}

	rows, err := s.db.Query(ctx, sql)
	if err != nil {
		return nil, failure.New("unable to find players", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	defer rows.Close()

	for rows.Next() {
		var pm PlayerModel
		err := pm.ScanRows(rows)
		if err != nil {
			return nil, failure.New("unable to find players", err)
		}

		dest = append(dest, pm)
	}

	if err = rows.Err(); err != nil {
		return nil, failure.New("unable to find players", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return dest, nil
}

func (s *store) findPlayer(ctx context.Context, playerId string) (*PlayerModel, error) {
	sql := `
		select
			player.id,
			player.height,
			player.weight,
			player.handedness,
			player.racket,
			player.matches_expected,
			player.matches_played,
			player.matches_won,
			player.matches_scheduled,
			player.seasons_played,
			account.id as account_id,
			account.name as account_name,
			league.id as league_id,
			league.title as league_title,
			player.created_at
		from player
		join account on player.account_id = account.id
		left join league on player.current_league_id = league.id
		where player.id = $1
	`

	var dest PlayerModel

	row := s.db.QueryRow(ctx, sql, playerId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("player not found", err)
		}
		return nil, failure.New("unable to find player", err)
	}

	return &dest, nil
}

func (s *store) updatePlayer(ctx context.Context, tx pgx.Tx, playerId string, model UpdatePlayerRequestModel) (*PlayerModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	sql := `
		with updated_player as (
			update player 
			set height = $1, weight = $2, handedness = $3, racket = $4
			where id = $5
			returning id, height, weight, handedness, racket, matches_expected, matches_played, matches_won, matches_scheduled, seasons_played, account_id, current_league_id, created_at
		)
		select 
			up.id,
			up.height,
			up.weight,
			up.handedness,
			up.racket,
			up.matches_expected,
			up.matches_played,
			up.matches_won,
			up.matches_scheduled,
			up.seasons_played,
			account.id as account_id,
			account.name as account_name,
			league.id as league_id,
			league.title as league_title,
			up.created_at
		from updated_player up
		join account on up.account_id = account.id
		left join league on up.current_league_id = league.id
	`

	var dest PlayerModel

	row := q.QueryRow(ctx, sql, model.Height, model.Weight, model.Handedness, model.Racket, playerId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("player for update not found", err)
		}
		return nil, failure.New("unable to update player", err)
	}

	return &dest, nil
}

// helper
func (s *store) checkPlayerOwnership(ctx context.Context, playerId, accountId string) (bool, error) {
	sql1 := `
		select exists (
			select 1 from player where id = $1 and account_id = $2
		)
	`

	var exists bool
	err := s.db.QueryRow(ctx, sql1, playerId, accountId).Scan(&exists)
	if err != nil {
		return false, failure.New("unable to check player ownership", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return exists, nil
}
