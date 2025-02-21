package me

import (
	"context"
	"errors"

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

func (s *store) findMe(ctx context.Context, accountId string) (*MeModel, error) {
	sql := `
		select 
			account.id as account_id,
			account.name as account_name,
			account.email as account_email,
			account.dob as account_dob,
			account.gender as account_gender,
			account.phone_number as account_phone_number,
			account.role as account_role,
			player.id as player_id,
			player.height as player_height,
			player.weight as player_weight,
			player.handedness as player_handedness,
			player.racket as player_racket,
			player.matches_expected as player_matches_expected,
			player.matches_played as player_matches_played,
			player.matches_won as player_matches_won,
			player.matches_scheduled as player_matches_scheduled,
			player.seasons_played as player_seasons_played,
			league.id as league_id,
			league.title as league_title,
			player.created_at as player_created_at,
			account.created_at as account_created_at
		from account
		left join player on account.id = player.account_id
		left join league on player.current_league_id = league.id
		where account.id = $1
	`

	var dest MeModel

	row := s.db.QueryRow(ctx, sql, accountId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("account not found", err)
		}
		return nil, failure.New("unable to find account", err)
	}

	return &dest, nil
}

func (s *store) updateMe(ctx context.Context, tx pgx.Tx, accountId string, model UpdateMeRequestModel) (*MeModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	var dest MeModel

	sql := `
		with updated_account as (
			update account 
			set name = $1
			where id = $2
			returning id, name, email, dob, gender, phone_number, role, created_at
		)
		select 
			ua.id as account_id,
			ua.name as account_name,
			ua.email as account_email,
			ua.dob as account_dob,
			ua.gender as account_gender,
			ua.phone_number as account_phone_number,
			ua.role as account_role,
			player.id as player_id,
			player.height as player_height,
			player.weight as player_weight,
			player.handedness as player_handedness,
			player.racket as player_racket,
			player.matches_expected as player_matches_expected,
			player.matches_played as player_matches_played,
			player.matches_won as player_matches_won,
			player.matches_scheduled as player_matches_scheduled,
			player.seasons_played as player_seasons_played,
			league.id as league_id,
			league.title as league_title,
			player.created_at as player_created_at,
			ua.created_at as account_created_at
		from updated_account ua
		left join player on ua.id = player.account_id
		left join league on player.current_league_id = league.id
	`

	row := q.QueryRow(ctx, sql, model.Name, accountId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("account for updating not found", err)
		}
		return nil, failure.New("unable to update account", err)
	}

	return &dest, nil
}
