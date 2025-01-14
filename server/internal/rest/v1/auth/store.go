package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/internal/db"
	"github.com/markovidakovic/gdsi/server/pkg/response"
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

func (s *store) readAccountByEmail(ctx context.Context, email string) (*Account, error) {
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
			return nil, response.ErrNotFound
		}
		return &account, response.ErrInternal
	}

	return &account, nil
}

func (s *store) insertRefreshToken(ctx context.Context, accountId string, token string, issuedAt, expiresAt time.Time) (inserted bool, err error) {
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
		err = response.ErrInternal
		return
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
		err = response.ErrInternal
		return
	}

	// insert the new refresh token
	_, err = tx.Exec(ctx, sql1, accountId, token, issuedAt, expiresAt)
	if err != nil {
		err = response.ErrInternal
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		err = response.ErrInternal
		return
	}

	inserted = true

	return
}

func newStore(db *db.Conn) *store {
	var r = &store{
		db,
	}
	return r
}
