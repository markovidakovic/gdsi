package seasons

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/response"
	"github.com/markovidakovic/gdsi/server/types"
)

type SeasonModel struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Creator     struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"creator"`
	CreatedAt time.Time `json:"created_at"`
}

func (sm *SeasonModel) ScanRow(row pgx.Row) error {
	err := row.Scan(&sm.Id, &sm.Title, &sm.Description, &sm.StartDate, &sm.EndDate, &sm.Creator.Id, &sm.Creator.Name, &sm.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("scanning season row: %w", response.ErrNotFound)
		}
		return fmt.Errorf("scanning season row: %v", err)
	}

	return nil
}

func (sm *SeasonModel) ScanRows(rows pgx.Rows) error {
	err := rows.Scan(&sm.Id, &sm.Title, &sm.Description, &sm.StartDate, &sm.EndDate, &sm.Creator.Id, &sm.Creator.Name, &sm.CreatedAt)
	if err != nil {
		return fmt.Errorf("scanning season rows: %v", err)
	}

	return nil
}

type CreateSeasonRequestModel struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	StartDate   types.Date `json:"start_date"`
	EndDate     types.Date `json:"end_date"`
	CreatorId   string     `json:"-"`
}

func (m CreateSeasonRequestModel) Validate() []response.InvalidField {
	var inv []response.InvalidField

	if m.Title == "" {
		inv = append(inv, response.InvalidField{
			Field:    "title",
			Message:  "Title is required",
			Location: "body",
		})
	}
	if m.EndDate.Time().Before(m.StartDate.Time()) {
		inv = append(inv, response.InvalidField{
			Field:    "end_date",
			Message:  "End date must be after start date",
			Location: "body",
		})
	}

	if len(inv) > 0 {
		return inv
	}

	return nil
}

type UpdateSeasonRequestModel struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	StartDate   types.Date `json:"start_date"`
	EndDate     types.Date `json:"end_date"`
}

func (m UpdateSeasonRequestModel) Validate() []response.InvalidField {
	var inv []response.InvalidField

	if m.Title == "" {
		inv = append(inv, response.InvalidField{
			Field:    "title",
			Message:  "Title is required",
			Location: "body",
		})
	}
	if m.EndDate.Time().Before(m.StartDate.Time()) {
		inv = append(inv, response.InvalidField{
			Field:    "end_date",
			Message:  "End date must be after start date",
			Location: "body",
		})
	}

	if len(inv) > 0 {
		return inv
	}

	return nil
}
