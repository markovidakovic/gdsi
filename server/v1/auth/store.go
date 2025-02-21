package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/failure"
)

type store struct {
	db *db.Conn
}

func newStore(db *db.Conn) *store {
	var r = &store{
		db,
	}
	return r
}

func (s *store) insertAccount(ctx context.Context, tx pgx.Tx, model SignupRequestModel) (AccountModel, error) {
	sql := `
		insert into account (name, email, dob, gender, phone_number, password)
		values ($1, $2, $3, $4, $5, $6)
		returning id, name, email, dob, gender, phone_number, password, role, NULL as player_id, created_at
	`

	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	var dest AccountModel
	row := q.QueryRow(ctx, sql, model.Name, model.Email, model.Dob, model.Gender, model.PhoneNumber, model.Password)
	err := dest.ScanRow(row)
	if err != nil {
		return dest, failure.New("failed to insert account", err)
	}

	return dest, nil
}

// todo: instead of returning string return the full model here
func (s *store) insertPlayer(ctx context.Context, tx pgx.Tx, accountId string) (string, error) {
	sql := `
		insert into player (account_id)
		values ($1)
		returning id
	`

	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	var playerId string

	err := q.QueryRow(ctx, sql, accountId).Scan(&playerId)
	if err != nil {
		return "", failure.New("failed to insert player", fmt.Errorf("%w: %v", failure.ErrInternal, err))
	}

	return playerId, nil
}

func (s *store) findAccountByEmail(ctx context.Context, tx pgx.Tx, email string) (*AccountModel, error) {
	var dest AccountModel

	sql := `
		select 
			account.id as account_id, 
			account.name as account_name, 
			account.email as account_email, 
			account.dob as account_dob, 
			account.gender as account_gender, 
			account.phone_number as account_phone_number, 
			account.password as account_password, 
			account.role as account_role, 
			player.id as player_id, 
			account.created_at as account_created_at
		from account
		left join player on player.account_id = account.id
		where account.email = $1
	`

	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	row := q.QueryRow(ctx, sql, email)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New(fmt.Sprintf("account with email %s not found", email), err)
		}
		return nil, failure.New("unable to retreive account", err)
	}

	return &dest, nil
}

func (s *store) insertRefreshToken(ctx context.Context, tx pgx.Tx, accountId string, token string, issuedAt, expiresAt time.Time) error {
	sql := `
		insert into refresh_token (account_id, token_hash, issued_at, expires_at)
		values ($1, $2, $3, $4)
		returning id, account_id, token_hash, device_id, ip_address, user_agent, issued_at, expires_at, last_used_at, is_revoked
	`

	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	_, err := q.Exec(ctx, sql, accountId, token, issuedAt, expiresAt)
	if err != nil {
		return failure.New("failed to insert refresh token", fmt.Errorf("%w: %v", failure.ErrInternal, err))
	}

	return nil
}

func (s *store) revokeAccountRefreshTokens(ctx context.Context, tx pgx.Tx, accountId string) error {
	sql := `
		update refresh_token
		set is_revoked = true
		where account_id = $1 and is_revoked = false
	`

	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	_, err := q.Exec(ctx, sql, accountId)
	if err != nil {
		return failure.New("failed to revoke account refresh tokens", fmt.Errorf("%w: %v", failure.ErrInternal, err))
	}

	return nil
}

func (s *store) findRefreshTokenByHash(ctx context.Context, tx pgx.Tx, rt string) (*RefreshTokenModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	sql := `
		select
			refresh_token.id,
			account.id as account_id,
			account.role as account_role,
			refresh_token.token_hash,
			refresh_token.device_id,
			refresh_token.ip_address,
			refresh_token.user_agent,
			refresh_token.issued_at,
			refresh_token.expires_at,
			refresh_token.last_used_at,
			refresh_token.is_revoked,
			player.id as player_id
		from refresh_token
		join account on refresh_token.account_id = account.id
		join player on account.id = player.account_id
		where refresh_token.token_hash = $1
	`

	var dest RefreshTokenModel

	row := q.QueryRow(ctx, sql, rt)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("refresh token not found", err)
		}
		return nil, failure.New("unable to retreive refresh token", err)
	}

	return &dest, nil
}

// todo: refactor this to return a full refresh token model and not just the error
func (s *store) updateRefreshToken(ctx context.Context, tx pgx.Tx, rtId string) error {
	sql := `
		update refresh_token
		set last_used_at = $1
		where id = $2
	`

	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	// last updated
	lua := time.Now()

	_, err := q.Exec(ctx, sql, lua, rtId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return failure.New("refresh token not found", fmt.Errorf("%w: %v", failure.ErrNotFound, err))
		}
		return failure.New("unable to update refresh token", fmt.Errorf("%w: %v", failure.ErrInternal, err))
	}

	return nil
}

func (s *store) revokeRefreshToken(ctx context.Context, tx pgx.Tx, rtId string) error {
	sql := `
		update refresh_token
		set is_revoked = true
		where id = $1	
	`

	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	_, err := q.Exec(ctx, sql, rtId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return failure.New("refresh token not found", fmt.Errorf("%w: %v", failure.ErrNotFound, err))
		}
		return failure.New("unable to revoke refresh token", fmt.Errorf("%w: %v", failure.ErrInternal, err))
	}

	return nil
}
