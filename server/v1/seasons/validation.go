package seasons

import (
	"github.com/markovidakovic/gdsi/server/response"
)

func validatePostSeason(input CreateSeasonRequestModel) []response.InvalidField {
	var inv []response.InvalidField

	if input.Title == "" {
		inv = append(inv, response.InvalidField{
			Field:    "title",
			Message:  "Title is required",
			Location: "body",
		})
	}
	if input.EndDate.Time().Before(input.StartDate.Time()) {
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

func validatePutSeason(input UpdateSeasonRequestModel) []response.InvalidField {
	var inv []response.InvalidField
	if input.Title == "" {
		inv = append(inv, response.InvalidField{
			Field:    "title",
			Message:  "Title is required",
			Location: "body",
		})
	}
	if input.EndDate.Time().Before(input.StartDate.Time()) {
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
