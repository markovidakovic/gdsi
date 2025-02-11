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
	"github.com/markovidakovic/gdsi/server/response"
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
		return "", "", fmt.Errorf("encrypting the password: %v", err)
	}

	// check if email exists
	existingAccount, err := s.store.findAccountByEmail(ctx, nil, model.Email)
	if err != nil {
		if !errors.Is(err, response.ErrNotFound) {
			return "", "", err
		}
	}
	if existingAccount != nil {
		return "", "", response.ErrDuplicateRecord
	}

	// begin tx
	tx, err := s.store.db.Begin(ctx)
	if err != nil {
		return "", "", err
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
		return "", "", err
	}

	// insert player
	playerId, err := s.store.insertPlayer(ctx, tx, account.Id)
	if err != nil {
		return "", "", err
	}

	account.PlayerId = playerId

	// generate jwts
	accessTkn, refreshTkn, err := generateAuthTokens(s.cfg.JwtAuth, s.cfg.JwtAccessExpiration, s.cfg.JwtRefreshExpiration, account.Id, account.Role, account.PlayerId)
	if err != nil {
		return "", "", fmt.Errorf("generating auth tokens: %v", err)
	}

	// hash the refresh token
	hashedRfrTkn := sec.HashToken(refreshTkn.val)

	// insert refresh token
	err = s.store.insertRefreshToken(ctx, tx, account.Id, hashedRfrTkn, refreshTkn.issAt, refreshTkn.expAt)
	if err != nil {
		return "", "", err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", "", fmt.Errorf("commiting tx: %v", err)
	}

	return accessTkn.val, refreshTkn.val, nil
}

func (s *service) processLogin(ctx context.Context, model LoginRequestModel) (string, string, error) {
	// call the store
	account, err := s.store.findAccountByEmail(ctx, nil, model.Email)
	if err != nil {
		return "", "", err
	}

	// validate password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(model.Password))
	if err != nil {
		return "", "", response.ErrNotFound
	}

	// generate jwts
	accessTkn, refreshTkn, err := generateAuthTokens(s.cfg.JwtAuth, s.cfg.JwtAccessExpiration, s.cfg.JwtRefreshExpiration, account.Id, account.Role, account.PlayerId)
	if err != nil {
		return "", "", err
	}

	// begin tx
	tx, err := s.store.db.Begin(ctx)
	if err != nil {
		return "", "", fmt.Errorf("beginning tx: %v", err)
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
		return "", "", err
	}

	// insert refresh token
	err = s.store.insertRefreshToken(ctx, tx, account.Id, sec.HashToken(refreshTkn.val), refreshTkn.issAt, refreshTkn.expAt)
	if err != nil {
		return "", "", err
	}

	// commit tx
	err = tx.Commit(ctx)
	if err != nil {
		return "", "", fmt.Errorf("commiting tx: %v", err)
	}

	return accessTkn.val, refreshTkn.val, nil
}

func (s *service) processRefreshTokens(ctx context.Context, model RefreshTokenRequestModel) (string, string, error) {
	// hash the incoming refresh token
	rtHash := sec.HashToken(model.RefreshToken)

	// begin tx
	tx, err := s.store.db.Begin(ctx)
	if err != nil {
		return "", "", fmt.Errorf("begining tx: %v", err)
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
		return "", "", fmt.Errorf("%w: refresh token is revoked", response.ErrUnauthorized)
	} else if time.Now().After(rt.ExpiresAt) {
		// revoke the expired refresh token
		err := s.store.revokeRefreshToken(ctx, nil, rt.Id)
		if err != nil {
			return "", "", err
		}
		return "", "", fmt.Errorf("%w: refresh token has expired", response.ErrUnauthorized)
	}

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
		return "", "", err
	}

	// insert refresh token
	err = s.store.insertRefreshToken(ctx, tx, rt.AccountId, sec.HashToken(refreshTkn.val), refreshTkn.issAt, refreshTkn.expAt)
	if err != nil {
		return "", "", err
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
