package auth

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/failure"
)

//go:embed queries/*.sql
var sqlFiles embed.FS

type store struct {
	db      *db.Conn
	queries struct {
		insertAccount              string
		insertPlayer               string
		findAccountByEmail         string
		insertRefreshToken         string
		revokeAccountRefreshTokens string
		findRefreshTokenByHash     string
		updateRefreshToken         string
		revokeRefreshToken         string
	}
}

func newStore(db *db.Conn) (*store, error) {
	var s = &store{
		db: db,
	}
	if err := s.loadQueries(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *store) loadQueries() error {
	insertAccountBytes, err := sqlFiles.ReadFile("queries/insert_account.sql")
	if err != nil {
		return fmt.Errorf("failed to read insert_account.sql file -> %v", err)
	}
	insertPlayerBytes, err := sqlFiles.ReadFile("queries/insert_player.sql")
	if err != nil {
		return fmt.Errorf("failed to read insert_player.sql file -> %v", err)
	}
	findAccountByEmailBytes, err := sqlFiles.ReadFile("queries/find_account_by_email.sql")
	if err != nil {
		return fmt.Errorf("failed to read find_account_by_email.sql file -> %v", err)
	}
	insertRefreshTokenBytes, err := sqlFiles.ReadFile("queries/insert_refresh_token.sql")
	if err != nil {
		return fmt.Errorf("failed to read insert_refresh_token.sql file -> %v", err)
	}
	revokeAccountRefreshTokensBytes, err := sqlFiles.ReadFile("queries/revoke_account_refresh_tokens.sql")
	if err != nil {
		return fmt.Errorf("failed to read revoke_account_refresh_tokens.sql file -> %v", err)
	}
	findRefreshTokenByHashBytes, err := sqlFiles.ReadFile("queries/find_refresh_token_by_hash.sql")
	if err != nil {
		return fmt.Errorf("failed to read find_refresh_token_by_hash.sql file -> %v", err)
	}
	updateRefreshTokenBytes, err := sqlFiles.ReadFile("queries/update_refresh_token.sql")
	if err != nil {
		return fmt.Errorf("failed to read update_refresh_token.sql file -> %v", err)
	}
	revokeRefreshTokenBytes, err := sqlFiles.ReadFile("queries/revoke_refresh_token.sql")
	if err != nil {
		return fmt.Errorf("failed to read revoke_refresh_token.sql file -> %v", err)
	}

	s.queries.insertAccount = string(insertAccountBytes)
	s.queries.insertPlayer = string(insertPlayerBytes)
	s.queries.findAccountByEmail = string(findAccountByEmailBytes)
	s.queries.insertRefreshToken = string(insertRefreshTokenBytes)
	s.queries.revokeAccountRefreshTokens = string(revokeAccountRefreshTokensBytes)
	s.queries.findRefreshTokenByHash = string(findRefreshTokenByHashBytes)
	s.queries.updateRefreshToken = string(updateRefreshTokenBytes)
	s.queries.revokeRefreshToken = string(revokeRefreshTokenBytes)

	return nil
}

func (s *store) insertAccount(ctx context.Context, tx pgx.Tx, model SignupRequestModel) (AccountModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	var dest AccountModel
	row := q.QueryRow(ctx, s.queries.insertAccount, model.Name, model.Email, model.Dob, model.Gender, model.PhoneNumber, model.Password)
	err := dest.ScanRow(row)
	if err != nil {
		return dest, failure.New("failed to insert account", err)
	}

	return dest, nil
}

// todo: instead of returning string return the full model here
func (s *store) insertPlayer(ctx context.Context, tx pgx.Tx, accountId string) (string, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	var playerId string

	err := q.QueryRow(ctx, s.queries.insertPlayer, accountId).Scan(&playerId)
	if err != nil {
		return "", failure.New("failed to insert player", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return playerId, nil
}

func (s *store) findAccountByEmail(ctx context.Context, tx pgx.Tx, email string) (*AccountModel, error) {
	var dest AccountModel
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	row := q.QueryRow(ctx, s.queries.findAccountByEmail, email)
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
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	_, err := q.Exec(ctx, s.queries.insertRefreshToken, accountId, token, issuedAt, expiresAt)
	if err != nil {
		return failure.New("failed to insert refresh token", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return nil
}

func (s *store) revokeAccountRefreshTokens(ctx context.Context, tx pgx.Tx, accountId string) error {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	_, err := q.Exec(ctx, s.queries.revokeAccountRefreshTokens, accountId)
	if err != nil {
		return failure.New("failed to revoke account refresh tokens", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
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

	var dest RefreshTokenModel
	row := q.QueryRow(ctx, s.queries.findRefreshTokenByHash, rt)
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
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	// last updated
	lua := time.Now()

	_, err := q.Exec(ctx, s.queries.updateRefreshToken, lua, rtId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return failure.New("refresh token not found", fmt.Errorf("%w -> %v", failure.ErrNotFound, err))
		}
		return failure.New("unable to update refresh token", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return nil
}

func (s *store) revokeRefreshToken(ctx context.Context, tx pgx.Tx, rtId string) error {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	_, err := q.Exec(ctx, s.queries.revokeRefreshToken, rtId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return failure.New("refresh token not found", fmt.Errorf("%w -> %v", failure.ErrNotFound, err))
		}
		return failure.New("unable to revoke refresh token", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return nil
}
