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

type AccessTokenResponseModel struct {
	AccessToken string `json:"access_token"`
}
