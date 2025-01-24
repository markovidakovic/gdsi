package auth

import "time"

// db table model
type Account struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Dob         time.Time `json:"dob"`
	Gender      string    `json:"gender"`
	PhoneNumber string    `json:"phone_number"`
	Password    string    `json:"-"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}

// db table model
type RefreshToken struct {
	Id         string     `json:"id"`
	AccountId  string     `json:"account_id"`
	TokenHash  string     `json:"token_hash"`
	DeviceId   *string    `json:"device_id"`
	IpAddress  *string    `json:"ip_address"`
	UserAgent  *string    `json:"user_agent"`
	IssuedAt   time.Time  `json:"issued_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
	IsRevoked  bool       `json:"is_revoked"`
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

// login request body model
type LoginRequestModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

// change forgotten password response body model
type ChangeForgottenPasswordResponseModel struct {
	Message string `json:"message"`
}

// refresh token request
type RefreshTokenRequestModel struct {
	RefreshToken string `json:"refresh_token"`
}
