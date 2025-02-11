package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/response"
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
	sql1 := `
		insert into account (name, email, dob, gender, phone_number, password)
		values ($1, $2, $3, $4, $5, $6)
		returning id, name, email, dob, gender, phone_number, password, role, created_at
	`

	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	var dest AccountModel
	err := q.QueryRow(ctx, sql1, model.Name, model.Email, model.Dob, model.Gender, model.PhoneNumber, model.Password).Scan(&dest.Id, &dest.Name, &dest.Email, &dest.Dob, &dest.Gender, &dest.PhoneNumber, &dest.Password, &dest.Role, &dest.CreatedAt)
	if err != nil {
		return dest, response.ErrInternal
	}

	return dest, nil
}

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
		return "", fmt.Errorf("inserting player: %v", err)
	}

	return playerId, nil
}

func (s *store) findAccountByEmail(ctx context.Context, tx pgx.Tx, email string) (*AccountModel, error) {
	var dest AccountModel

	sql := `
		select account.id, account.name, account.email, account.dob, account.gender, account.phone_number, account.password, account.role, player.id as player_id, account.created_at
		from account
		join player on player.account_id = account.id
		where account.email = $1
	`

	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	err := q.QueryRow(ctx, sql, email).Scan(
		&dest.Id,
		&dest.Name,
		&dest.Email,
		&dest.Dob,
		&dest.Gender,
		&dest.PhoneNumber,
		&dest.Password,
		&dest.Role,
		&dest.PlayerId,
		&dest.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("finding account: %w", response.ErrNotFound)
		}
		return nil, err
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
		return err
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
		return fmt.Errorf("revoking refresh tokens: %v", err)
	}

	return nil
}

func (s *store) findRefreshTokenByHash(ctx context.Context, tx pgx.Tx, rt string) (*RefreshTokenModel, error) {
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

	err := tx.QueryRow(ctx, sql, rt).Scan(&dest.Id, &dest.AccountId, &dest.AccountRole, &dest.TokenHash, &dest.DeviceId, &dest.IpAddress, &dest.UserAgent, &dest.IssuedAt, &dest.ExpiresAt, &dest.LastUsedAt, &dest.IsRevoked, &dest.PlayerId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("finding refresh token: %w", response.ErrNotFound)
		}
		return nil, err
	}

	return &dest, nil
}

func (s *store) updateRefreshToken(ctx context.Context, tx pgx.Tx, rtHash string) error {
	sql := `
		update refresh_token
		set last_used_at = $1
		where token_hash = $2
	`

	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	// last updated
	lua := time.Now()

	_, err := q.Exec(ctx, sql, lua, rtHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("finding refresh token for updated: %w", response.ErrNotFound)
		}
		return fmt.Errorf("updating refresh token: %v", err)
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
		return fmt.Errorf("revoking refresh token: %v", err)
	}

	return nil
}
