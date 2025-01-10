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
	var result Account

	query := `
		INSERT INTO account (first_name, last_name, email, dob, gender, phone_number, password)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, first_name, last_name, email, dob, gender, phone_number, password
	`

	err := s.db.QueryRow(ctx, query,
		model.FirstName,
		model.LastName,
		model.Email,
		model.Dob,
		model.Gender,
		model.PhoneNumber,
		model.Password,
	).Scan(
		&result.Id,
		&result.FirstName,
		&result.LastName,
		&result.Email,
		&result.Dob,
		&result.Gender,
		&result.PhoneNumber,
		&result.Password,
	)
	if err != nil {
		return result, response.ErrInternal
	}

	return result, nil
}

func (s *store) selectAccountByEmail(ctx context.Context, email string) (*Account, error) {
	var result Account

	query := `
		SELECT id, first_name, last_name, email, dob, gender, phone_number, password
		FROM account
		WHERE email = $1
	`

	err := s.db.QueryRow(ctx, query, email).Scan(
		&result.Id,
		&result.FirstName,
		&result.LastName,
		&result.Email,
		&result.Dob,
		&result.Gender,
		&result.PhoneNumber,
		&result.Password,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &result, response.ErrNotFound
		}
		return &result, response.ErrInternal
	}

	return &result, nil
}

func (s *store) insertRefreshToken(ctx context.Context, accountId string, token string, issuedAt, expiresAt time.Time) (*RefreshToken, error) {
	var result RefreshToken

	sql := `
		INSERT INTO refresh_token (account_id, token_hash, issued_at, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, account_id, token_hash, device_id, ip_address, user_agent, issued_at, expires_at, last_used_at, is_revoked
	`

	err := s.db.QueryRow(ctx, sql, accountId, token, issuedAt, expiresAt).Scan(
		&result.Id,
		&result.AccountId,
		&result.TokenHash,
		&result.DeviceId,
		&result.IpAddress,
		&result.UserAgent,
		&result.IssuedAt,
		&result.ExpiresAt,
		&result.LastUsedAt,
		&result.IsRevoked,
	)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, response.ErrInternal
	}

	return &result, nil
}

func newStore(db *db.Conn) *store {
	var r = &store{
		db,
	}
	return r
}
