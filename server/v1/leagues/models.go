package leagues

import (
	"time"
)

type League struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	SeasonId    string    `json:"season_id"`
	CreatorId   string    `json:"creator_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreatorModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SeasonModel struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type LeagueModel struct {
	Id          string       `json:"id"`
	Title       string       `json:"title"`
	Description *string      `json:"description"`
	Season      SeasonModel  `json:"season"`
	Creator     CreatorModel `json:"creator"`
	CreatedAt   time.Time    `json:"created_at"`
}

type CreateLeagueRequestModel struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type UpdateLeagueRequestModel struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}
