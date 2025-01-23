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

func (s *store) insertAccount(ctx context.Context, model SignupRequestModel) (Account, error) {
	sql1 := `
		insert into account (name, email, dob, gender, phone_number, password)
		values ($1, $2, $3, $4, $5, $6)
		returning id, name, email, dob, gender, phone_number, password
	`

	sql2 := `
		insert into player (account_id)
		values ($1)
	`

	var account Account

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return account, response.ErrInternal
	}

	// tx rollback
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	err = tx.QueryRow(ctx, sql1, model.Name, model.Email, model.Dob, model.Gender, model.PhoneNumber, model.Password).Scan(&account.Id, &account.Name, &account.Email, &account.Dob, &account.Gender, &account.PhoneNumber, &account.Password)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return account, response.ErrInternal
	}

	_, err = tx.Exec(ctx, sql2, account.Id)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return account, response.ErrInternal
	}

	err = tx.Commit(ctx)
	if err != nil {
		return account, response.ErrInternal
	}

	return account, nil
}

func (s *store) findAccountByEmail(ctx context.Context, email string) (*Account, error) {
	var account Account

	sql := `
		select id, name, email, dob, gender, phone_number, password
		from account
		where email = $1
	`

	err := s.db.QueryRow(ctx, sql, email).Scan(
		&account.Id,
		&account.Name,
		&account.Email,
		&account.Dob,
		&account.Gender,
		&account.PhoneNumber,
		&account.Password,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: account not found", response.ErrNotFound)
		}
		return nil, err
	}

	return &account, nil
}

func (s *store) insertRefreshToken(ctx context.Context, accountId string, token string, issuedAt, expiresAt time.Time) error {
	sql1 := `
		insert into refresh_token (account_id, token_hash, issued_at, expires_at)
		values ($1, $2, $3, $4)
		returning id, account_id, token_hash, device_id, ip_address, user_agent, issued_at, expires_at, last_used_at, is_revoked
	`

	sql2 := `
		update refresh_token
		set is_revoked = true
		where account_id = $1 and is_revoked = false
	`

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	// tx rollback
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// revoke the previous refresh token
	_, err = tx.Exec(ctx, sql2, accountId)
	if err != nil {
		return err
	}

	// insert the new refresh token
	_, err = tx.Exec(ctx, sql1, accountId, token, issuedAt, expiresAt)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *store) findRefreshToken(ctx context.Context, rt string) (*RefreshToken, error) {
	// maybe do just an update query where we do returning
	sql1 := `
		select id, account_id, token_hash, device_id, ip_address, user_agent, issued_at, expires_at, last_used_at, is_revoked
		from refresh_token
		where token_hash = $1
	`

	sql2 := `
		update refresh_token
		set last_used_at = $1
		where token_hash = $2
	`

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	var dest RefreshToken

	err = tx.QueryRow(ctx, sql1, rt).Scan(&dest.Id, &dest.AccountId, &dest.TokenHash, &dest.DeviceId, &dest.IpAddress, &dest.UserAgent, &dest.IssuedAt, &dest.ExpiresAt, &dest.LastUsedAt, &dest.IsRevoked)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: refresh token not found", response.ErrNotFound)
		}
		return nil, err
	}

	lua := time.Now()

	_, err = tx.Exec(ctx, sql2, lua, rt)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &dest, nil
}

func (s *store) revokeAllAccountRefreshTokens(ctx context.Context, accountId string) error {
	sql1 := `
		update refresh_token
		set is_revoked = true
		where account_id = $1 and is_revoked = false
	`

	_, err := s.db.Exec(ctx, sql1, accountId)
	if err != nil {
		return err
	}

	return nil
}

func (s *store) revokeRefreshToken(ctx context.Context, rtId string) error {
	sql1 := `
		update refresh_token
		set is_revoked = true
		where id = $1
	`

	_, err := s.db.Exec(ctx, sql1, rtId)
	if err != nil {
		return err
	}

	return nil
}

func newStore(db *db.Conn) *store {
	var r = &store{
		db,
	}
	return r
}
