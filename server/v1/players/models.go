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
	Id                  string       `json:"id"`
	Height              *float64     `json:"height"`
	Weight              *float64     `json:"weight"`
	Handedness          *string      `json:"handedness"`
	Racket              *string      `json:"racket"`
	MatchesExpected     int          `json:"matches_expected"`
	MatchesPlayed       int          `json:"matches_played"`
	MatchesWon          int          `json:"matches_won"`
	MatchesScheduled    int          `json:"matches_scheduled"`
	SeasonsPlayed       int          `json:"seasons_played"`
	Elo                 int          `json:"elo"`
	HighestElo          int          `json:"highest_elo"`
	IsEloProvisional    bool         `json:"is_elo_provisional"`
	Account             accountModel `json:"account"`
	CurrentLeague       *leagueModel `json:"current_league"`
	PreviousLeague      *leagueModel `json:"previous_league"`
	PreviousLeagueRank  *int32       `json:"previous_league_rank"`
	IsPlayingNextSeason bool         `json:"is_playing_next_season"`
	CreatedAt           time.Time    `json:"created_at"`
}

type accountModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type leagueModel struct {
	Id   string `json:"id"`
	Tier int32  `json:"tier"`
	Name string `json:"name"`
}

func (pm *PlayerModel) ScanRow(row pgx.Row) error {
	var currLeagueId, currLeagueName, prevLeagueId, prevLeagueName sql.NullString
	var currLeagueTier, prevLeagueTier, prevLeagueRank sql.NullInt32
	err := row.Scan(
		&pm.Id,
		&pm.Height,
		&pm.Weight,
		&pm.Handedness,
		&pm.Racket,
		&pm.MatchesExpected,
		&pm.MatchesPlayed,
		&pm.MatchesWon,
		&pm.MatchesScheduled,
		&pm.SeasonsPlayed,
		&pm.Elo,
		&pm.HighestElo,
		&pm.IsEloProvisional,
		&pm.Account.Id,
		&pm.Account.Name,
		&currLeagueId,
		&currLeagueTier,
		&currLeagueName,
		&prevLeagueId,
		&prevLeagueTier,
		&prevLeagueName,
		&prevLeagueRank,
		&pm.IsPlayingNextSeason,
		&pm.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return failure.New("scanning player row", fmt.Errorf("%w -> %v", failure.ErrNotFound, err))
		}
		return failure.New("database error scanning player row", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	if !currLeagueId.Valid {
		pm.CurrentLeague = nil
	} else {
		pm.CurrentLeague = &leagueModel{
			Id:   currLeagueId.String,
			Tier: currLeagueTier.Int32,
			Name: currLeagueName.String,
		}
	}
	if !prevLeagueId.Valid {
		pm.PreviousLeague = nil
	} else {
		pm.PreviousLeague = &leagueModel{
			Id:   prevLeagueId.String,
			Tier: prevLeagueRank.Int32,
			Name: prevLeagueName.String,
		}
	}
	if !prevLeagueRank.Valid {
		pm.PreviousLeagueRank = nil
	} else {
		pm.PreviousLeagueRank = &prevLeagueRank.Int32
	}

	return nil
}

func (pm *PlayerModel) ScanRows(rows pgx.Rows) error {
	var currLeagueId, currLeagueName, prevLeagueId, prevLeagueName sql.NullString
	var currLeagueTier, prevLeagueTier, prevLeagueRank sql.NullInt32
	err := rows.Scan(
		&pm.Id,
		&pm.Height,
		&pm.Weight,
		&pm.Handedness,
		&pm.Racket,
		&pm.MatchesExpected,
		&pm.MatchesPlayed,
		&pm.MatchesWon,
		&pm.MatchesScheduled,
		&pm.SeasonsPlayed,
		&pm.Elo,
		&pm.HighestElo,
		&pm.IsEloProvisional,
		&pm.Account.Id,
		&pm.Account.Name,
		&currLeagueId,
		&currLeagueTier,
		&currLeagueName,
		&prevLeagueId,
		&prevLeagueTier,
		&prevLeagueName,
		&prevLeagueRank,
		&pm.IsPlayingNextSeason,
		&pm.CreatedAt,
	)
	if err != nil {
		return failure.New("database error scanning player rows", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	if !currLeagueId.Valid {
		pm.CurrentLeague = nil
	} else {
		pm.CurrentLeague = &leagueModel{
			Id:   currLeagueId.String,
			Tier: currLeagueTier.Int32,
			Name: currLeagueName.String,
		}
	}
	if !prevLeagueId.Valid {
		pm.PreviousLeague = nil
	} else {
		pm.PreviousLeague = &leagueModel{
			Id:   prevLeagueId.String,
			Tier: prevLeagueRank.Int32,
			Name: prevLeagueName.String,
		}
	}
	if !prevLeagueRank.Valid {
		pm.PreviousLeagueRank = nil
	} else {
		pm.PreviousLeagueRank = &prevLeagueRank.Int32
	}

	return nil
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
