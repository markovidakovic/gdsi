package seasons

import "time"

type Season struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type SeasonModel struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateSeasonModel struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type UpdateSeasonModel struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}
