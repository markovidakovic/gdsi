package courts

import "github.com/markovidakovic/gdsi/server/response"

func validatePostCourt(input CreateCourtRequestModel) []response.InvalidField {
	var inv []response.InvalidField

	if input.Name == "" {
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

func validatePutCourt(input UpdateCourtRequestModel) []response.InvalidField {
	var inv []response.InvalidField
	if input.Name == "" {
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
