package matches

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/failure"
)

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

func (mm *MatchModel) ScanRow(row pgx.Row) error {
	var winnerId, winnerName sql.NullString
	err := row.Scan(&mm.Id, &mm.Court.Id, &mm.Court.Name, &mm.ScheduledAt, &mm.PlayerOne.Id, &mm.PlayerOne.Name, &mm.PlayerTwo.Id, &mm.PlayerTwo.Name, &winnerId, &winnerName, &mm.Score, &mm.Season.Id, &mm.Season.Title, &mm.League.Id, &mm.League.Title, &mm.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return failure.New("scanning match row", fmt.Errorf("%w -> %v", failure.ErrNotFound, err))
		}
		return failure.New("database error scanning match row", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	if !winnerId.Valid {
		mm.Winner = nil
	} else {
		mm.Winner = &PlayerModel{
			Id:   winnerId.String,
			Name: winnerName.String,
		}
	}

	return nil
}

func (mm *MatchModel) ScanRows(rows pgx.Rows) error {
	var winnerId, winnerName sql.NullString
	err := rows.Scan(&mm.Id, &mm.Court.Id, &mm.Court.Name, &mm.ScheduledAt, &mm.PlayerOne.Id, &mm.PlayerOne.Name, &mm.PlayerTwo.Id, &mm.PlayerTwo.Name, &winnerId, &winnerName, &mm.Score, &mm.Season.Id, &mm.Season.Title, &mm.League.Id, &mm.League.Title, &mm.CreatedAt)
	if err != nil {
		return failure.New("database error scanning match rows", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	if !winnerId.Valid {
		mm.Winner = nil
	} else {
		mm.Winner = &PlayerModel{
			Id:   winnerId.String,
			Name: winnerName.String,
		}
	}
	return nil
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

func (m CreateMatchRequestModel) Validate() []failure.InvalidField {
	var inv []failure.InvalidField

	if m.CourtId == "" {
		inv = append(inv, failure.InvalidField{
			Field:    "court_id",
			Message:  "Court id is required",
			Location: "body",
		})
	}
	if m.ScheduledAt == "" {
		inv = append(inv, failure.InvalidField{
			Field:    "scheduled_at",
			Message:  "Scheduled at is required",
			Location: "body",
		})
	}
	if m.PlayerTwoId == "" {
		inv = append(inv, failure.InvalidField{
			Field:    "player_two_id",
			Message:  "Player two id is required",
			Location: "body",
		})
	}
	if m.Score != nil {
		if *m.Score == "" {
			inv = append(inv, failure.InvalidField{
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
func (m UpdateMatchRequestModel) Validate() []failure.InvalidField {
	var inv []failure.InvalidField

	if m.CourtId == "" {
		inv = append(inv, failure.InvalidField{
			Field:    "court_id",
			Message:  "Court id is required",
			Location: "body",
		})
	}
	if m.ScheduledAt == "" {
		inv = append(inv, failure.InvalidField{
			Field:    "scheduled_at",
			Message:  "Scheduled at is required",
			Location: "body",
		})
	}
	if m.PlayerTwoId == "" {
		inv = append(inv, failure.InvalidField{
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
	Score       string `json:"score"`
	SeasonId    string `json:"-"`
	LeagueId    string `json:"-"`
	MatchId     string `json:"-"`
	WinnerId    string `json:"-"`
	PlayerOneId string `json:"-"`
	PlayerTwoId string `json:"-"`
}

func (m SubmitMatchScoreRequestModel) Validate() []failure.InvalidField {
	var inv []failure.InvalidField

	if m.Score == "" || !isValidScore(m.Score) {
		inv = append(inv, failure.InvalidField{
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

	if len(sets) < 2 {
		return false
	}

	var pl1SetsWon, pl2SetsWon int

	// check first two sets
	for i := 0; i < 2; i++ {
		setSl := strings.Split(sets[i], "-")
		if len(setSl) != 2 {
			return false
		}

		pl1Games, err1 := strconv.Atoi(setSl[0])
		pl2Games, err2 := strconv.Atoi(setSl[1])
		if err1 != nil || err2 != nil {
			return false
		}

		if !isValidSet(pl1Games, pl2Games) {
			return false
		}

		if pl1Games > pl2Games {
			pl1SetsWon++
		} else {
			pl2SetsWon++
		}
	}

	// after two sets, check if someone won both
	if pl1SetsWon == 2 || pl2SetsWon == 2 {
		return len(sets) == 2 // can't have more than two sets in the score
	}

	// must be 1-1 at this point and must have a 3rd set or super tie-break
	if len(sets) != 3 {
		return false
	}

	// validate 3rd set or tie-break
	setSl := strings.Split(sets[2], "-")
	if len(setSl) != 2 {
		return false
	}

	pl1Games, err1 := strconv.Atoi(setSl[0])
	pl2Games, err2 := strconv.Atoi(setSl[1])
	if err1 != nil || err2 != nil {
		return false
	}

	if pl1Games >= 10 || pl2Games >= 10 {
		if !isValidTieBreak(pl1Games, pl2Games) {
			return false
		}
	} else if !isValidSet(pl1Games, pl2Games) {
		return false
	}

	if pl1Games > pl2Games {
		pl1SetsWon++
	} else {
		pl2SetsWon++
	}

	// final score must be 2-1 either way
	return (pl1SetsWon == 2 && pl2SetsWon == 1) || (pl2SetsWon == 2 && pl1SetsWon == 1)
}

func isValidSet(games1, games2 int) bool {
	if games1 == 7 {
		return games2 == 5 || games2 == 6 // 7-5 or 7-6 (tie-break)
	}
	if games2 == 7 {
		return games1 == 5 || games1 == 6 // 5-7 or 6-7 (tie-break)
	}
	if games1 == 6 {
		return games2 <= 4 // 6-0 through 6-4
	}
	if games2 == 6 {
		return games1 <= 4 // 0-6 through 4-6
	}
	return false
}

func isValidTieBreak(score1, score2 int) bool {
	// super tie-break
	if score1 == 10 {
		// if excatly 10, opponent must have 8 or less
		return score2 <= 8
	}
	if score2 == 10 {
		return score1 <= 8
	}
	// if more than 10, must win by 2
	if score1 > 10 {
		return score1-score2 == 2
	}
	if score2 > 10 {
		return score2-score1 == 2
	}
	return false
}
