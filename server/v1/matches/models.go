package matches

import (
	"fmt"
	"strconv"
	"strings"
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
	if m.Score != nil {
		if *m.Score == "" {
			inv = append(inv, response.InvalidField{
				Field:    "score",
				Message:  "Invalid score value",
				Location: "body",
			})
		} else {

		}
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

// todo
func (m SubmitMatchScoreRequestModel) Validate() []response.InvalidField {
	var inv []response.InvalidField

	if m.Score == "" || !isValidScore(m.Score) {
		inv = append(inv, response.InvalidField{
			Field:    "score",
			Message:  "Score is not valid",
			Location: "body",
		})
	}

	if len(inv) > 0 {
		return inv
	}
	return nil
}

func isValidScore(score string) bool {
	sets := strings.Split(score, ",")

	fmt.Printf("sets: %v\n", sets)

	if len(sets) < 2 || len(sets) > 3 {
		return false
	}

	var pl1SetsWon, pl2SetsWon int

	for i, set := range sets {
		setSl := strings.Split(set, "-")
		if len(setSl) != 2 {
			return false
		}

		fmt.Printf("setSl: %v\n", setSl)

		gamesPl1, err1 := strconv.Atoi(setSl[0])
		gamesPl2, err2 := strconv.Atoi(setSl[1])
		if err1 != nil || err2 != nil {
			return false
		}

		fmt.Printf("gamesPl1: %v\n", gamesPl1)
		fmt.Printf("gamesPl2: %v\n", gamesPl2)

		// set games validation (first two sets, and a possible third set)
		if i < 2 {
			if !isValidRegularSet(gamesPl1, gamesPl2) {
				return false
			}
		} else {
			if !isValidThirdSet(gamesPl1, gamesPl2) {
				return false
			}
		}

		if gamesPl1 > gamesPl2 {
			pl1SetsWon++
		} else {
			pl2SetsWon++
		}
	}

	fmt.Printf("pl1SetsWon: %v\n", pl1SetsWon)
	fmt.Printf("pl2SetsWon: %v\n", pl2SetsWon)

	return false
}

func isValidRegularSet(gamesPl1, gamesPl2 int) bool {
	return false
}

func isValidThirdSet(scorePl1, scorePl2 int) bool {
	return false
}
