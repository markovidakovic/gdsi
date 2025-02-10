package players

import (
	"time"

	"github.com/markovidakovic/gdsi/server/response"
)

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

// todo:
func (m UpdatePlayerRequestModel) Validate() []response.InvalidField {
	var inv []response.InvalidField

	if len(inv) > 0 {
		return inv
	}

	return nil
}
