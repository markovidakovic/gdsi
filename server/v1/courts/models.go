package courts

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/failure"
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
			return failure.New("scanning court row", fmt.Errorf("%w -> %v", failure.ErrNotFound, err))
		}
		return failure.New("database error", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	return nil
}

func (cm *CourtModel) ScanRows(rows pgx.Rows) error {
	err := rows.Scan(&cm.Id, &cm.Name, &cm.Creator.Id, &cm.Creator.Name, &cm.CreatedAt)
	if err != nil {
		return failure.New("database error", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	return nil
}

type CreateCourtRequestModel struct {
	Name      string `json:"name"`
	CreatorId string `json:"-"`
}

func (m CreateCourtRequestModel) Validate() []failure.InvalidField {
	var inv []failure.InvalidField

	if m.Name == "" {
		inv = append(inv, failure.InvalidField{
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

func (m UpdateCourtRequestModel) Validate() []failure.InvalidField {
	var inv []failure.InvalidField

	if m.Name == "" {
		inv = append(inv, failure.InvalidField{
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
