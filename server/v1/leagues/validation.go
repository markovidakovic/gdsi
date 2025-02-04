package leagues

import "github.com/markovidakovic/gdsi/server/response"

func validatePostLeague(input CreateLeagueRequestModel) []response.InvalidField {
	var inv []response.InvalidField

	if input.Title == "" {
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

func validatePutLeague(input UpdateLeagueRequestModel) []response.InvalidField {
	var inv []response.InvalidField

	if input.Title == "" {
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
