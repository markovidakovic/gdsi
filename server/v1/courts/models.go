package courts

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/response"
)

type CourtModel struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Creator struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"creator"`
	CreatedAt time.Time `json:"created_at"`
}

func (cm *CourtModel) ScanRow(row pgx.Row) error {
	err := row.Scan(&cm.Id, &cm.Name, &cm.Creator.Id, &cm.Creator.Name, &cm.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("scanning court row: %w", response.ErrNotFound)
		}
		return fmt.Errorf("scanning court row: %v", err)
	}
	return nil
}

func (cm *CourtModel) ScanRows(rows pgx.Rows) error {
	err := rows.Scan(&cm.Id, &cm.Name, &cm.Creator.Id, &cm.Creator.Name, &cm.CreatedAt)
	if err != nil {
		return fmt.Errorf("scanning court rows: %v", err)
	}
	return nil
}

type CreateCourtRequestModel struct {
	Name      string `json:"name"`
	CreatorId string `json:"-"`
}

func (m CreateCourtRequestModel) Validate() []response.InvalidField {
	var inv []response.InvalidField

	if m.Name == "" {
		inv = append(inv, response.InvalidField{
			Field:    "name",
			Message:  "Name field is required",
			Location: "body",
		})
	}

	if len(inv) > 0 {
		return inv
	}

	return nil
}

type UpdateCourtRequestModel struct {
	Name string `json:"name"`
}

func (m UpdateCourtRequestModel) Validate() []response.InvalidField {
	var inv []response.InvalidField

	if m.Name == "" {
		inv = append(inv, response.InvalidField{
			Field:    "name",
			Message:  "Name field is required",
			Location: "body",
		})
	}

	if len(inv) > 0 {
		return inv
	}

	return nil
}
