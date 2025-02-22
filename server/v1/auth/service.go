package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/sec"
)

type service struct {
	cfg   *config.Config
	store *store
}

func newService(cfg *config.Config, store *store) *service {
	var s = &service{
		cfg,
		store,
	}
	return s
}

func (s *service) processSignup(ctx context.Context, model SignupRequestModel) (string, string, error) {
	var err error

	// hash the password
	model.Password, err = sec.EncryptPwd(model.Password)
	if err != nil {
		return "", "", failure.New("signup failed", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	// check if email exists
	existingAccount, err := s.store.findAccountByEmail(ctx, nil, model.Email)
	if err != nil {
		if !errors.Is(err, failure.ErrNotFound) {
			return "", "", failure.New("signup failed", err)
		}
	}
	if existingAccount != nil {
		return "", "", failure.New("email already registered", failure.ErrDuplicate)
	}

	// begin tx
	tx, err := s.store.db.Begin(ctx)
	if err != nil {
		return "", "", failure.New("signup failed", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	defer func() {
		if tx != nil {
			err := tx.Rollback(ctx)
			if err != nil && err != pgx.ErrTxClosed {
				log.Printf("rolling back tx: %v", err)
			}
		}
	}()

	// insert account
	account, err := s.store.insertAccount(ctx, tx, model)
	if err != nil {
		return "", "", failure.New("signup failed", err)
	}

	// insert player
	playerId, err := s.store.insertPlayer(ctx, tx, account.Id)
	if err != nil {
		return "", "", failure.New("signup failed", err)
	}

	account.PlayerId = &playerId

	// generate jwts
	accessTkn, refreshTkn, err := generateAuthTokens(s.cfg.JwtAuth, s.cfg.JwtAccessExpiration, s.cfg.JwtRefreshExpiration, account.Id, account.Role, *account.PlayerId)
	if err != nil {
		return "", "", failure.New("signup failed", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	// hash the refresh token
	hashedRfrTkn := sec.HashToken(refreshTkn.val)

	// insert refresh token
	err = s.store.insertRefreshToken(ctx, tx, account.Id, hashedRfrTkn, refreshTkn.issAt, refreshTkn.expAt)
	if err != nil {
		return "", "", failure.New("signup failed", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", "", failure.New("signup failed", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return accessTkn.val, refreshTkn.val, nil
}

func (s *service) processLogin(ctx context.Context, model LoginRequestModel) (string, string, error) {
	// call the store
	account, err := s.store.findAccountByEmail(ctx, nil, model.Email)
	if err != nil {
		// special case here. the findAccountByEmail method returns failure.ErrNotFound or failure.ErrInternal
		// in the login endpoint we don't want to return the failure.ErrNotFound if the account has not been found.
		// rather, we want to return the failure.ErrBadRequest, so we disregard the previous error from the store method
		// the drawback is that the error msg from the store method will not be logged - for now this is ok
		if errors.Is(err, failure.ErrNotFound) {
			return "", "", failure.New("invalid email or password", failure.ErrBadRequest)
		}
		return "", "", failure.New("login failed", err)
	}

	// validate password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(model.Password))
	if err != nil {
		// todo: maybe refactor this err msg later
		return "", "", failure.New("invalid email or password", fmt.Errorf("%w -> %v", failure.ErrBadRequest, err))
	}

	// generate jwts
	accessTkn, refreshTkn, err := generateAuthTokens(s.cfg.JwtAuth, s.cfg.JwtAccessExpiration, s.cfg.JwtRefreshExpiration, account.Id, account.Role, *account.PlayerId)
	if err != nil {
		return "", "", failure.New("login failed", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	// begin tx
	tx, err := s.store.db.Begin(ctx)
	if err != nil {
		return "", "", failure.New("login failed", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	defer func() {
		if tx != nil {
			err := tx.Rollback(ctx)
			if err != nil && err != pgx.ErrTxClosed {
				log.Printf("rolling back tx: %v", err)
			}
		}
	}()

	// revoke previous refresh tkns
	err = s.store.revokeAccountRefreshTokens(ctx, tx, account.Id)
	if err != nil {
		return "", "", failure.New("login failed", err)
	}

	// insert refresh token
	err = s.store.insertRefreshToken(ctx, tx, account.Id, sec.HashToken(refreshTkn.val), refreshTkn.issAt, refreshTkn.expAt)
	if err != nil {
		return "", "", failure.New("login failed", err)
	}

	// commit tx
	err = tx.Commit(ctx)
	if err != nil {
		return "", "", failure.New("login failed", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return accessTkn.val, refreshTkn.val, nil
}

func (s *service) processRefreshTokens(ctx context.Context, model RefreshTokenRequestModel) (string, string, error) {
	// hash the incoming refresh token
	rtHash := sec.HashToken(model.RefreshToken)

	// begin tx
	tx, err := s.store.db.Begin(ctx)
	if err != nil {
		return "", "", failure.New("refresh tokens failed", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	defer func() {
		if tx != nil {
			err := tx.Rollback(ctx)
			if err != nil && err == pgx.ErrTxClosed {
				log.Printf("rollbacking tx: %v", err)
			}
		}
	}()

	// get the stored refresh token
	rt, err := s.store.findRefreshTokenByHash(ctx, nil, rtHash)
	if err != nil {
		return "", "", err
	}

	if rt.IsRevoked {
		// revoke all existing refresh tokens for the account
		err := s.store.revokeAccountRefreshTokens(ctx, nil, rt.AccountId)
		if err != nil {
			return "", "", err
		}
		return "", "", failure.New("refresh token revoked", failure.ErrUnauthorized)
	} else if time.Now().After(rt.ExpiresAt) {
		// revoke the expired refresh token
		err := s.store.revokeRefreshToken(ctx, nil, rt.Id)
		if err != nil {
			return "", "", err
		}
		return "", "", failure.New("refresh token expired", failure.ErrUnauthorized)
	}

	// update prev RT - it just updates the lat column
	err = s.store.updateRefreshToken(ctx, tx, rt.Id)
	if err != nil {
		return "", "", err
	}

	// revoke the previous refresh token
	err = s.store.revokeRefreshToken(ctx, tx, rt.Id)
	if err != nil {
		return "", "", err
	}

	// generate tokens
	accessTkn, refreshTkn, err := generateAuthTokens(s.cfg.JwtAuth, s.cfg.JwtAccessExpiration, s.cfg.JwtRefreshExpiration, rt.AccountId, rt.AccountRole, rt.PlayerId)
	if err != nil {
		return "", "", failure.New("refresh tokens failed", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	// insert refresh token
	err = s.store.insertRefreshToken(ctx, tx, rt.AccountId, sec.HashToken(refreshTkn.val), refreshTkn.issAt, refreshTkn.expAt)
	if err != nil {
		return "", "", failure.New("refresh tokens failed", err)
	}

	// commit tx
	err = tx.Commit(ctx)
	if err != nil {
		return "", "", failure.New("refresh tokens failed", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return accessTkn.val, refreshTkn.val, nil
}

type token struct {
	issAt time.Time
	expAt time.Time
	val   string
}

func generateAuthTokens(ja *jwtauth.JWTAuth, jwtAccessExp, jwtRefreshExp, accountId, role, playerId string) (accessTkn, refreshTkn token, err error) {
	// parse config vars
	durAccess, err := time.ParseDuration(jwtAccessExp)
	if err != nil {
		return
	}
	durRefresh, err := time.ParseDuration(jwtRefreshExp)
	if err != nil {
		return
	}

	var iss string = "gdsi api"
	var aud string = "gdsi app"

	now := time.Now()

	// expiration dates for tokens
	var expAccess time.Time = now.Add(durAccess)
	var expRefresh time.Time = now.Add(durRefresh)

	// jwt claims
	var claims = map[string]interface{}{
		"iss":       iss,
		"sub":       accountId,
		"aud":       aud,
		"exp":       expAccess.Unix(),
		"nbf":       now.Unix(),
		"iat":       now.Unix(),
		"role":      role,
		"player_id": playerId,
	}

	// encode access jwt
	_, accessTknEnc, err := ja.Encode(claims)
	if err != nil {
		return
	}

	// change the exp value for refresh token
	claims["exp"] = expRefresh.Unix()

	// encode refresh jwt
	_, refreshTknEnc, err := ja.Encode(claims)
	if err != nil {
		return
	}

	accessTkn = token{
		issAt: now,
		expAt: expAccess,
		val:   accessTknEnc,
	}
	refreshTkn = token{
		issAt: now,
		expAt: expRefresh,
		val:   refreshTknEnc,
	}

	return
}
