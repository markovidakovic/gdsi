package me

import (
	"time"

	"github.com/markovidakovic/gdsi/server/response"
)

type MeModel struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Email       string      `json:"email"`
	Dob         time.Time   `json:"dob"`
	Gender      string      `json:"gender"`
	PhoneNumber string      `json:"phone_number"`
	Role        string      `json:"role"`
	CreatedAt   time.Time   `json:"created_at"`
	Player      PlayerModel `json:"player"`
}

type PlayerModel struct {
	Id               string              `json:"id"`
	Height           *float64            `json:"height"`
	Weight           *float64            `json:"weight"`
	Handedness       *string             `json:"handedness"`
	Racket           *string             `json:"racket"`
	MatchesExpected  int                 `json:"matches_expected"`
	MatchesPlayed    int                 `json:"matches_played"`
	MatchesWon       int                 `json:"matches_won"`
	MatchesScheduled int                 `json:"matches_scheduled"`
	SeasonsPlayed    int                 `json:"seasons_played"`
	CurrentLeague    *CurrentLeagueModel `json:"current_league"`
	CreatedAt        time.Time           `json:"created_at"`
}

type CurrentLeagueModel struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateMeRequestModel struct {
	Name string `json:"name"`
}

// todo:
func (m UpdateMeRequestModel) Validate() []response.InvalidField {
	var inv []response.InvalidField

	if len(inv) > 0 {
		return inv
	}

	return nil
}

type UpdatePasswordRequestModel struct {
	OldPassword         string `json:"old_password"`
	NewPassword         string `json:"new_password"`
	RepeatedNewPassword string `json:"repeated_new_password"`
}

// todo:
func (m UpdatePasswordRequestModel) Validate() []response.InvalidField {
	var inv []response.InvalidField

	if len(inv) > 0 {
		return inv
	}

	return nil
}

type UpdatePasswordResponseModel struct {
	Message string `json:"message"`
}
