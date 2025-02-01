package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/response"
	"github.com/markovidakovic/gdsi/server/security"
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
	// Hash the password
	model.Password, err = security.EncryptPwd(model.Password)
	if err != nil {
		return "", "", fmt.Errorf("encrypting the password: %v", err)
	}

	// Check if email exists
	existingAccount, err := s.store.findAccountByEmail(ctx, model.Email)
	if err != nil {
		if !errors.Is(err, response.ErrNotFound) {
			return "", "", err
		}
	}
	if existingAccount != nil {
		return "", "", response.ErrDuplicateRecord
	}

	// Insert account
	newAccount, err := s.store.insertAccount(ctx, model)
	if err != nil {
		return "", "", err
	}

	// Generate jwt
	accessToken, refreshToken, err := s.getAuthTokens(ctx, newAccount.Id, newAccount.Role)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *service) processLogin(ctx context.Context, model LoginRequestModel) (string, string, error) {
	// Call the store
	account, err := s.store.findAccountByEmail(ctx, model.Email)
	if err != nil {
		return "", "", err
	}

	// Validate password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(model.Password))
	if err != nil {
		return "", "", response.ErrNotFound
	}

	// Generate jwt
	accessToken, refreshToken, err := s.getAuthTokens(ctx, account.Id, account.Role)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *service) processRefreshTokens(ctx context.Context, model RefreshTokenRequestModel) (string, string, error) {
	// hash the incoming refresh token
	rtHash := hashToken(model.RefreshToken)

	// get the stored refresh token
	rt, role, err := s.store.findRefreshToken(ctx, rtHash)
	if err != nil {
		return "", "", err
	}

	if rt.IsRevoked {
		// revoke all existing refresh tokens for the account
		err := s.store.revokeAllAccountRefreshTokens(ctx, rt.AccountId)
		if err != nil {
			return "", "", err
		}
		return "", "", fmt.Errorf("%w: refresh token is revoked", response.ErrUnauthorized)
	}

	if time.Now().After(rt.ExpiresAt) {
		// revoke the expired refresh token
		err := s.store.revokeRefreshToken(ctx, rt.Id)
		if err != nil {
			return "", "", err
		}
		return "", "", fmt.Errorf("%w: refresh token expired", response.ErrUnauthorized)
	}

	// generate tokens
	access, refresh, err := s.getAuthTokens(ctx, rt.AccountId, role)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *service) getAuthTokens(ctx context.Context, id, role string) (access, refresh string, err error) {
	var iss string = "gdsi api"
	var aud string = "gdsi app"

	// parse config vars
	durationAccess, err := time.ParseDuration(s.cfg.JwtAccessExpiration)
	if err != nil {
		return
	}
	durationRefresh, err := time.ParseDuration(s.cfg.JwtRefreshExpiration)
	if err != nil {
		return
	}

	now := time.Now()

	// expiration dates for tokens
	var expAccess time.Time = now.Add(durationAccess)
	var expRefresh time.Time = now.Add(durationRefresh)

	// jwt claims
	var claims = map[string]interface{}{
		"iss":  iss,
		"sub":  id,
		"aud":  aud,
		"exp":  expAccess.Unix(),
		"nbf":  now.Unix(),
		"iat":  now.Unix(),
		"role": role,
	}

	// create access jwt
	_, access, err = s.cfg.JwtAuth.Encode(claims)
	if err != nil {
		return
	}

	// change the exp value
	claims["exp"] = expRefresh.Unix()

	// create refresh token
	_, refresh, err = s.cfg.JwtAuth.Encode(claims)
	if err != nil {
		return
	}

	// hash the refresh token for storage
	refreshHashed := hashToken(refresh)

	// call the store
	err = s.store.insertRefreshToken(ctx, id, refreshHashed, now, expRefresh)
	if err != nil {
		return
	}

	return
}

func hashToken(val string) string {
	hash := sha256.Sum256([]byte(val))
	return hex.EncodeToString(hash[:])
}
