package auth

import (
	"context"
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

func (s *service) signupNewAccount(ctx context.Context, model SignupRequestModel) (string, error) {
	// Hash the password
	pwdBytes, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", response.ErrInternal
	}
	model.Password = string(pwdBytes)

	// Check if email exists
	existingAccount, err := s.store.selectAccountByEmail(ctx, model.Email)
	if err != nil {
		if !errors.Is(err, response.ErrNotFound) {
			return "", err
		}
	}
	if existingAccount != nil {
		return "", response.ErrDuplicateRecord
	}

	// Call the store
	newAccount, err := s.store.insertAccount(ctx, model)
	if err != nil {
		return "", err
	}

	// Generate jwt
	now := time.Now()
	jwtDuration, err := time.ParseDuration(s.cfg.JwtExpiration)
	if err != nil {
		return "", response.ErrInternal
	}

	_, token, err := s.cfg.JwtAuth.Encode(map[string]interface{}{
		"iss": "gdsi api",
		"sub": newAccount.Id,
		"aud": "gdsi app",
		"exp": now.Add(jwtDuration).Unix(),
		"nbf": now.Unix(),
		"iat": now.Unix(),
	})
	if err != nil {
		return "", response.ErrInternal
	}

	return token, nil
}

func (s *service) getAccessToken(ctx context.Context, model LoginRequestModel) (string, error) {
	// Call the store
	account, err := s.store.selectAccountByEmail(ctx, model.Email)
	if err != nil {
		return "", err
	}

	// Validate password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(model.Password))
	if err != nil {
		return "", response.ErrNotFound
	}

	// Generate jwt
	now := time.Now()
	jwtDuration, err := time.ParseDuration(s.cfg.JwtExpiration)
	if err != nil {
		return "", response.ErrInternal
	}

	_, token, err := s.cfg.JwtAuth.Encode(map[string]interface{}{
		"iss": "gdsi api",
		"sub": account.Id,
		"aud": "gdsi app",
		"exp": now.Add(jwtDuration).Unix(),
		"nbf": now.Unix(),
		"iat": now.Unix(),
	})
	if err != nil {
		return "", response.ErrInternal
	}

	return token, nil
}

func newService(cfg *config.Config, store *store) *service {
	var s = &service{
		cfg,
		store,
	}
	return s
}
