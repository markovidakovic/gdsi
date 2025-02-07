package matches

import (
	"time"

	"github.com/markovidakovic/gdsi/server/response"
)

// db table
type Match struct {
	Id          string    `json:"id"`
	CourtId     string    `json:"court_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	PlayerOneId string    `json:"player_one_id"`
	PlayerTwoId string    `json:"player_two_id"`
	WinnerId    *string   `json:"winner_id"`
	Score       *string   `json:"score"`
	SeasonId    string    `json:"season_id"`
	LeagueId    string    `json:"league_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type MatchModel struct {
	Id          string       `json:"id"`
	Court       CourtModel   `json:"court"`
	ScheduledAt time.Time    `json:"scheduled_at"`
	PlayerOne   PlayerModel  `json:"player_one"`
	PlayerTwo   PlayerModel  `json:"player_two"`
	Winner      *PlayerModel `json:"winner"`
	Score       *string      `json:"score"`
	Season      SeasonModel  `json:"season"`
	League      LeagueModel  `json:"league"`
	CreatedAt   time.Time    `json:"created_at"`
}

type CourtModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type PlayerModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SeasonModel struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type LeagueModel struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

// create match
type CreateMatchRequestModel struct {
	CourtId     string  `json:"court_id"`
	ScheduledAt string  `json:"scheduled_at"`
	PlayerOneId string  `json:"-"`
	PlayerTwoId string  `json:"player_two_id"`
	WinnerId    *string `json:"-"`
	Score       *string `json:"score"`
	SeasonId    string  `json:"-"`
	LeagueId    string  `json:"-"`
}

func (m CreateMatchRequestModel) Validate() []response.InvalidField {
	var inv []response.InvalidField

	if m.CourtId == "" {
		inv = append(inv, response.InvalidField{
			Field:    "court_id",
			Message:  "Court id is required",
			Location: "body",
		})
	}
	if m.ScheduledAt == "" {
		inv = append(inv, response.InvalidField{
			Field:    "scheduled_at",
			Message:  "Scheduled at is required",
			Location: "body",
		})
	}
	if m.PlayerTwoId == "" {
		inv = append(inv, response.InvalidField{
			Field:    "player_two_id",
			Message:  "Player two id is required",
			Location: "body",
		})
	}
	if *m.Score == "" {
		inv = append(inv, response.InvalidField{
			Field:    "score",
			Message:  "Invalid score value",
			Location: "body",
		})
	}

	if len(inv) > 0 {
		return inv
	}

	return nil
}

// update match
type UpdateMatchRequestModel struct {
	CourtId     string `json:"court_id"`
	ScheduledAt string `json:"scheduled_at"`
	PlayerOneId string `json:"-"`
	PlayerTwoId string `json:"player_two_id"`
	SeasonId    string `json:"-"`
	LeagueId    string `json:"-"`
	MatchId     string `json:"-"`
}

// todo:
func (m UpdateMatchRequestModel) Validate() []response.InvalidField {
	var inv []response.InvalidField

	if len(inv) > 0 {
		return inv
	}

	return nil
}

// submit score
type SubmitMatchScoreRequestModel struct {
	Score    string `json:"score"`
	SeasonId string `json:"-"`
	LeagueId string `json:"-"`
	MatchId  string `json:"-"`
}
