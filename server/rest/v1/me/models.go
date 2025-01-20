package me

import "time"

type MeModel struct {
	Id            string             `json:"id"`
	Name          string             `json:"name"`
	Email         string             `json:"email"`
	Dob           time.Time          `json:"dob"`
	Gender        string             `json:"gender"`
	PhoneNumber   string             `json:"phone_number"`
	CreatedAt     time.Time          `json:"created_at"`
	PlayerProfile PlayerProfileModel `json:"player_profile"`
}

type PlayerProfileModel struct {
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
	WinningRation    float64             `json:"winning_ratio"`
	ActivityRatio    float64             `json:"activity_ratio"`
	Ranking          *int                `json:"ranking"`
	Elo              *int                `json:"elo"`
	CurrentLeague    *CurrentLeagueModel `json:"current_league"`
	CreatedAt        time.Time           `json:"created_at"`
}

type CurrentLeagueModel struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateMeModel struct {
	Name string `json:"name"`
}
