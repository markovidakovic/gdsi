package me

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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

func (mm *MeModel) ScanRow(row pgx.Row) error {
	var leagueId, leagueTitle sql.NullString
	err := row.Scan(&mm.Id, &mm.Name, &mm.Email, &mm.Dob, &mm.Gender, &mm.PhoneNumber, &mm.Role, &mm.CreatedAt, &mm.Player.Id, &mm.Player.Height, &mm.Player.Weight, &mm.Player.Handedness, &mm.Player.Racket, &mm.Player.MatchesExpected, &mm.Player.MatchesPlayed, &mm.Player.MatchesWon, &mm.Player.MatchesScheduled, &mm.Player.SeasonsPlayed, &leagueId, &leagueTitle, &mm.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("scanning me row: %w", response.ErrNotFound)
		}
		return fmt.Errorf("scanning me row: %v", err)
	}

	if !leagueId.Valid {
		mm.Player.CurrentLeague = nil
	} else {
		mm.Player.CurrentLeague = &CurrentLeagueModel{
			Id:    leagueId.String,
			Title: leagueTitle.String,
		}
	}

	return nil
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
	Id    string `json:"id"`
	Title string `json:"title"`
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
