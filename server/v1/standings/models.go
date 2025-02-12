package standings

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/response"
)

type StandingModel struct {
	Id            string `json:"id"`
	Points        int    `json:"points"`
	MatchesPlayed int    `json:"matches_played"`
	MatchesWon    int    `json:"matches_won"`
	SetsWon       int    `json:"sets_won"`
	SetsLost      int    `json:"sets_lost"`
	GamesWon      int    `json:"games_won"`
	GamesLost     int    `json:"games_lost"`
	Season        struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"season"`
	League struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"league"`
	Player struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"player"`
	CreatedAt time.Time `json:"created_at"`
}

func (sm *StandingModel) ScanRow(row pgx.Row) error {
	err := row.Scan(&sm.Id, &sm.Points, &sm.MatchesPlayed, &sm.MatchesWon, &sm.SetsWon, &sm.SetsLost, &sm.GamesWon, &sm.GamesLost, &sm.Season.Id, &sm.Season.Name, &sm.League.Id, &sm.League.Name, &sm.Player.Id, &sm.Player.Name, &sm.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("scanning row: %w", response.ErrNotFound)
		}
		return fmt.Errorf("scanning row: %v", err)
	}

	return nil
}

func (sm *StandingModel) ScanRows(rows pgx.Rows) error {
	err := rows.Scan(&sm.Id, &sm.Points, &sm.MatchesPlayed, &sm.MatchesWon, &sm.SetsWon, &sm.SetsLost, &sm.GamesWon, &sm.GamesLost, &sm.Season.Id, &sm.Season.Name, &sm.League.Id, &sm.League.Name, &sm.Player.Id, &sm.Player.Name, &sm.CreatedAt)
	if err != nil {
		return fmt.Errorf("scanning rows: %v", err)
	}

	return nil
}
