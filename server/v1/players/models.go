package players

import "time"

// db table
type Player struct {
	Id              string    `json:"id"`
	Height          *float64  `json:"height"`
	Weight          *float64  `json:"weight"`
	Handedness      *string   `json:"handedness"`
	Racket          *string   `json:"racket"`
	MatchesExpected int       `json:"matches_expected"`
	MatchesPlayed   int       `json:"matches_played"`
	MatchesWon      int       `json:"matches_won"`
	SeasonsPlayed   int       `json:"seasons_played"`
	WinningRatio    float64   `json:"winning_ratio"`
	ActivityRatio   float64   `json:"activity_ratio"`
	Ranking         *int      `json:"ranking"`
	Elo             *int      `json:"elo"`
	AccountId       string    `json:"account_id"`
	CurrentLeagueId *string   `json:"current_league_id"`
	CreatedAt       time.Time `json:"created_at"`
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
	WinningRatio     float64             `json:"winning_ratio"`
	ActivityRatio    float64             `json:"activity_ratio"`
	Ranking          *int                `json:"ranking"`
	Elo              *int                `json:"elo"`
	Account          AccountModel        `json:"account"`
	CurrentLeague    *CurrentLeagueModel `json:"current_league"`
	CreatedAt        time.Time           `json:"created_at"`
}

type AccountModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CurrentLeagueModel struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type UpdatePlayerRequestModel struct {
	Height     *float64 `json:"height"`
	Weight     *float64 `json:"weight"`
	Handedness *string  `json:"handedness"`
	Racket     *string  `json:"racket"`
}
