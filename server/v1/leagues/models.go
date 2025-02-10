package leagues

import (
	"time"

	"github.com/markovidakovic/gdsi/server/response"
)

type LeagueModel struct {
	Id          string       `json:"id"`
	Title       string       `json:"title"`
	Description *string      `json:"description"`
	Season      SeasonModel  `json:"season"`
	Creator     CreatorModel `json:"creator"`
	CreatedAt   time.Time    `json:"created_at"`
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

func (m CreateLeagueRequestModel) Validate() []response.InvalidField {
	var inv []response.InvalidField

	if m.Title == "" {
		inv = append(inv, response.InvalidField{
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

func (m UpdateLeagueRequestModel) Validate() []response.InvalidField {
	var inv []response.InvalidField

	if m.Title == "" {
		inv = append(inv, response.InvalidField{
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
