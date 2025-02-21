package leagues

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/failure"
)

type LeagueModel struct {
	Id          string       `json:"id"`
	Title       string       `json:"title"`
	Description *string      `json:"description"`
	Season      SeasonModel  `json:"season"`
	Creator     CreatorModel `json:"creator"`
	CreatedAt   time.Time    `json:"created_at"`
}

func (lm *LeagueModel) ScanRow(row pgx.Row) error {
	err := row.Scan(&lm.Id, &lm.Title, &lm.Description, &lm.Season.Id, &lm.Season.Title, &lm.Creator.Id, &lm.Creator.Name, &lm.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return failure.New("scanning league row", fmt.Errorf("%w -> %v", failure.ErrNotFound, err))
		}
		return failure.New("database error", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	return nil
}

func (lm *LeagueModel) ScanRows(rows pgx.Rows) error {
	err := rows.Scan(&lm.Id, &lm.Title, &lm.Description, &lm.Season.Id, &lm.Season.Title, &lm.Creator.Id, &lm.Creator.Name, &lm.CreatedAt)
	if err != nil {
		return failure.New("database error", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	return nil
}

type CreatorModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SeasonModel struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type CreateLeagueRequestModel struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	CreatorId   string  `json:"-"`
	SeasonId    string  `json:"-"`
}

func (m CreateLeagueRequestModel) Validate() []failure.InvalidField {
	var inv []failure.InvalidField

	if m.Title == "" {
		inv = append(inv, failure.InvalidField{
			Field:    "title",
			Message:  "Title is required",
			Location: "body",
		})
	}

	if len(inv) > 0 {
		return inv
	}

	return nil
}

type UpdateLeagueRequestModel struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	SeasonId    string  `json:"-"`
	LeagueId    string  `json:"-"`
}

func (m UpdateLeagueRequestModel) Validate() []failure.InvalidField {
	var inv []failure.InvalidField

	if m.Title == "" {
		inv = append(inv, failure.InvalidField{
			Field:    "title",
			Message:  "Title is required",
			Location: "body",
		})
	}

	if len(inv) > 0 {
		return inv
	}

	return nil
}
