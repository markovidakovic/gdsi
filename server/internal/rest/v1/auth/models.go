package auth

import "time"

type Account struct {
	Id          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Dob         time.Time `json:"dob"`
	Gender      string    `json:"gender"`
	PhoneNumber string    `json:"phone_number"`
	Password    string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
}

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

type SignupRequestModel struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Dob         string `json:"dob"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginRequestModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokensResponseModel struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
