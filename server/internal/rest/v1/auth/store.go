package auth

import (
	"context"
	"errors"

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

func newStore(db *db.Conn) *store {
	var r = &store{
		db,
	}
	return r
}
