package db

import "time"

// db table
type Account struct {
	Id          string
	Name        string
	Email       string
	Dob         time.Time
	Gender      string
	PhoneNumber string
	Password    string
	Role        string
	CreatedAt   time.Time
}
