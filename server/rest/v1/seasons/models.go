package seasons

import "time"

type Season struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
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
	Creator     CreatorModel `json:"creator"`
	CreatedAt   time.Time    `json:"created_at"`
}

type CreateSeasonModel struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type UpdateSeasonModel struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}
