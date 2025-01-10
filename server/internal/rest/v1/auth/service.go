package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/pkg/response"
)

type service struct {
	cfg   *config.Config
	store *store
}

func (s *service) signup(ctx context.Context, model SignupRequestModel) (string, string, error) {
	// Hash the password
	pwdBytes, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", response.ErrInternal
	}
	model.Password = string(pwdBytes)

	// Check if email exists
	existingAccount, err := s.store.selectAccountByEmail(ctx, model.Email)
	if err != nil {
		if !errors.Is(err, response.ErrNotFound) {
			return "", "", err
		}
	}
	if existingAccount != nil {
		return "", "", response.ErrDuplicateRecord
	}

	// Call the store
	newAccount, err := s.store.insertAccount(ctx, model)
	if err != nil {
		return "", "", err
	}

	// Generate jwt
	accessToken, refreshToken, err := s.getAuthTokens(ctx, newAccount.Id)
	if err != nil {
		return "", "", response.ErrInternal
	}

	return accessToken, refreshToken, nil
}

func (s *service) login(ctx context.Context, model LoginRequestModel) (string, string, error) {
	// Call the store
	account, err := s.store.selectAccountByEmail(ctx, model.Email)
	if err != nil {
		return "", "", err
	}

	// Validate password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(model.Password))
	if err != nil {
		return "", "", response.ErrNotFound
	}

	// Generate jwt
	accessToken, refreshToken, err := s.getAuthTokens(ctx, account.Id)
	if err != nil {
		return "", "", response.ErrInternal
	}

	return accessToken, refreshToken, nil
}

func (s *service) getAuthTokens(ctx context.Context, accountId string) (access, refresh string, err error) {
	var iss string = "gdsi api"
	var aud string = "gdsi app"

	// Parse config vars
	durationAccess, err := time.ParseDuration(s.cfg.JwtAccessExpiration)
	if err != nil {
		return
	}
	durationRefresh, err := time.ParseDuration(s.cfg.JwtRefreshExpiration)
	if err != nil {
		return
	}

	now := time.Now()

	// Expiration dates for tokens
	var expAccess time.Time = now.Add(durationAccess)
	var expRefresh time.Time = now.Add(durationRefresh)

	// Jwt claims
	var claims = map[string]interface{}{
		"iss": iss,
		"sub": accountId,
		"aud": aud,
		"exp": expAccess.Unix(),
		"nbf": now.Unix(),
		"iat": now.Unix(),
	}

	// Create access jwt
	_, access, err = s.cfg.JwtAuth.Encode(claims)
	if err != nil {
		return
	}

	// Change the exp value
	claims["exp"] = expRefresh.Unix()

	// Create refresh token
	_, refresh, err = s.cfg.JwtAuth.Encode(claims)
	if err != nil {
		return
	}

	// Hash the refresh token for storage
	hash := sha256.Sum256([]byte(refresh))
	refreshHashed := hex.EncodeToString(hash[:])

	// Call the store
	_, err = s.store.insertRefreshToken(ctx, accountId, refreshHashed, now, expRefresh)
	if err != nil {
		return
	}

	return
}

func newService(cfg *config.Config, store *store) *service {
	var s = &service{
		cfg,
		store,
	}
	return s
}
