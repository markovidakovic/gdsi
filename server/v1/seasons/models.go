package seasons

import (
	"time"

	"github.com/markovidakovic/gdsi/server/types"
)

type Season struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CreatorId   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreatorModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SeasonModel struct {
	Id          string       `json:"id"`
	Title       string       `json:"title"`
	Description *string      `json:"description"`
	StartDate   time.Time    `json:"start_date"`
	EndDate     time.Time    `json:"end_date"`
	Creator     CreatorModel `json:"creator"`
	CreatedAt   time.Time    `json:"created_at"`
}

type CreateSeasonRequestModel struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	StartDate   types.Date `json:"start_date"`
	EndDate     types.Date `json:"end_date"`
	CreatorId   string     `json:"-"`
}

type UpdateSeasonRequestModel struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	StartDate   types.Date `json:"start_date"`
	EndDate     types.Date `json:"end_date"`
}
