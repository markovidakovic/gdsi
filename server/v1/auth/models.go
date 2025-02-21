package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/sec"
)

type AccountModel struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Dob         time.Time `json:"dob"`
	Gender      string    `json:"gender"`
	PhoneNumber string    `json:"phone_number"`
	Password    string    `json:"-"`
	Role        string    `json:"role"`
	PlayerId    *string   `json:"player_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func (am *AccountModel) ScanRow(row pgx.Row) error {
	err := row.Scan(&am.Id, &am.Name, &am.Email, &am.Dob, &am.Gender, &am.PhoneNumber, &am.Password, &am.Role, &am.PlayerId, &am.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return failure.New("scanning account row", fmt.Errorf("%w: %v", failure.ErrNotFound, err))
		}
		return failure.New("database error", fmt.Errorf("%w: %v", failure.ErrInternal, err))
	}

	return nil
}

type RefreshTokenModel struct {
	Id          string     `json:"id"`
	AccountId   string     `json:"account_id"`
	AccountRole string     `json:"account_role"`
	TokenHash   string     `json:"token_hash"`
	DeviceId    *string    `json:"device_id"`
	IpAddress   *string    `json:"ip_address"`
	UserAgent   *string    `json:"user_agent"`
	IssuedAt    time.Time  `json:"issued_at"`
	ExpiresAt   time.Time  `json:"expires_at"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	IsRevoked   bool       `json:"is_revoked"`
	PlayerId    string     `json:"player_id"`
}

func (rtm *RefreshTokenModel) ScanRow(row pgx.Row) error {
	err := row.Scan(&rtm.Id, &rtm.AccountId, &rtm.AccountRole, &rtm.TokenHash, &rtm.DeviceId, &rtm.IpAddress, &rtm.UserAgent, &rtm.IssuedAt, &rtm.ExpiresAt, &rtm.LastUsedAt, &rtm.IsRevoked, &rtm.PlayerId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return failure.New("scanning refresh token row", fmt.Errorf("%w: %v", failure.ErrNotFound, err))
		}
		return failure.New("database error", fmt.Errorf("%w: %v", failure.ErrInternal, err))
	}
	return nil
}

// signup request body model
type SignupRequestModel struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Dob         string `json:"dob"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func (m SignupRequestModel) Validate() []failure.InvalidField {
	var inv []failure.InvalidField

	if m.Name == "" {
		inv = append(inv, failure.InvalidField{
			Field:    "name",
			Message:  "Name field is required",
			Location: "body",
		})
	}
	if m.Email == "" {
		inv = append(inv, failure.InvalidField{
			Field:    "email",
			Message:  "Email field is required",
			Location: "body",
		})
	} else if !sec.IsValidEmail(m.Email) {
		inv = append(inv, failure.InvalidField{
			Field:    "email",
			Message:  "Invalid email",
			Location: "body",
		})
	}
	if m.Dob == "" {
		inv = append(inv, failure.InvalidField{
			Field:    "dob",
			Message:  "Date of birth field is required",
			Location: "body",
		})
	} else {
		if _, err := time.Parse("2006-01-02", m.Dob); err != nil {
			inv = append(inv, failure.InvalidField{
				Field:    "dob",
				Message:  "Invalid date format",
				Location: "body",
			})
		}
	}
	if m.Gender == "" {
		inv = append(inv, failure.InvalidField{
			Field:    "gender",
			Message:  "Gender field is required",
			Location: "body",
		})
	} else if m.Gender != "male" && m.Gender != "female" {
		inv = append(inv, failure.InvalidField{
			Field:    "gender",
			Message:  "Invalid gender, expected male or female",
			Location: "body",
		})
	}
	if m.PhoneNumber == "" {
		inv = append(inv, failure.InvalidField{
			Field:    "phone_number",
			Message:  "Phone number field is required",
			Location: "body",
		})
	} else if !sec.IsValidPhone(m.PhoneNumber) {
		inv = append(inv, failure.InvalidField{
			Field:    "phone_number",
			Message:  "Invalid phone number",
			Location: "body",
		})
	}
	if m.Password == "" {
		inv = append(inv, failure.InvalidField{
			Field:    "password",
			Message:  "Password field required",
			Location: "body",
		})
	}

	if len(inv) > 0 {
		return inv
	}

	return nil
}

// login request body model
type LoginRequestModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m LoginRequestModel) Validate() []failure.InvalidField {
	var inv []failure.InvalidField

	if m.Email == "" {
		inv = append(inv, failure.InvalidField{
			Field:    "email",
			Message:  "Email field is required",
			Location: "body",
		})
	} else if !sec.IsValidEmail(m.Email) {
		inv = append(inv, failure.InvalidField{
			Field:    "email",
			Message:  "Invalid email",
			Location: "body",
		})
	}
	if m.Password == "" {
		inv = append(inv, failure.InvalidField{
			Field:    "password",
			Message:  "Password field is required",
			Location: "body",
		})
	}

	if len(inv) > 0 {
		return inv
	}

	return nil
}

// tokens response model
// v1/signup
// v1/tokens/access
type TokensResponseModel struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// forgotten password request body model
type ForgottenPasswordRequestModel struct {
	Email string `json:"email"`
}

func (m ForgottenPasswordRequestModel) Validate() []failure.InvalidField {
	var inv []failure.InvalidField

	if len(inv) > 0 {
		return inv
	}

	return nil
}

// forgotten password response body model
type ForgottenPasswordResponseModel struct {
	Message string `json:"message"`
}

// change forgotten password request body model
type ChangeForgottenPasswordRequestModel struct {
	Code            string `json:"code"`
	Email           string `json:"email"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (m ChangeForgottenPasswordRequestModel) Validate() []failure.InvalidField {
	var inv []failure.InvalidField

	if len(inv) > 0 {
		return inv
	}

	return nil
}

// change forgotten password response body model
type ChangeForgottenPasswordResponseModel struct {
	Message string `json:"message"`
}

// refresh token request
type RefreshTokenRequestModel struct {
	RefreshToken string `json:"refresh_token"`
}

func (m RefreshTokenRequestModel) Validate() []failure.InvalidField {
	var inv []failure.InvalidField

	if m.RefreshToken == "" {
		inv = append(inv, failure.InvalidField{
			Field:    "refresh_token",
			Message:  "Refresh token field is required",
			Location: "body",
		})
	}

	if len(inv) > 0 {
		return inv
	}

	return nil
}
