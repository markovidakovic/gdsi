package players

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/failure"
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

func (pm *PlayerModel) ScanRow(row pgx.Row) error {
	var leagueId, leagueTitle sql.NullString
	err := row.Scan(&pm.Id, &pm.Height, &pm.Weight, &pm.Handedness, &pm.Racket, &pm.MatchesExpected, &pm.MatchesPlayed, &pm.MatchesWon, &pm.MatchesScheduled, &pm.SeasonsPlayed, &pm.Account.Id, &pm.Account.Name, &leagueId, &leagueTitle, &pm.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return failure.New("scanning player row", fmt.Errorf("%w -> %v", failure.ErrNotFound, err))
		}
		return failure.New("database error scanning player row", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	if !leagueId.Valid {
		pm.CurrentLeague = nil
	} else {
		pm.CurrentLeague = &CurrentLeagueModel{
			Id:    leagueId.String,
			Title: leagueTitle.String,
		}
	}

	return nil
}

func (pm *PlayerModel) ScanRows(rows pgx.Rows) error {
	var leagueId, leagueTitle sql.NullString
	err := rows.Scan(&pm.Id, &pm.Height, &pm.Weight, &pm.Handedness, &pm.Racket, &pm.MatchesExpected, &pm.MatchesPlayed, &pm.MatchesWon, &pm.MatchesScheduled, &pm.SeasonsPlayed, &pm.Account.Id, &pm.Account.Name, &leagueId, &leagueTitle, &pm.CreatedAt)
	if err != nil {
		return failure.New("database error scanning player rows", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	if !leagueId.Valid {
		pm.CurrentLeague = nil
	} else {
		pm.CurrentLeague = &CurrentLeagueModel{
			Id:    leagueId.String,
			Title: leagueTitle.String,
		}
	}
	return nil
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
func (m UpdatePlayerRequestModel) Validate() []failure.InvalidField {
	var inv []failure.InvalidField

	if len(inv) > 0 {
		return inv
	}

	return nil
}
