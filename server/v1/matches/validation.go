package matches

import "github.com/markovidakovic/gdsi/server/response"

func validatePostMatch(input CreateMatchRequestModel) []response.InvalidField {
	var inv []response.InvalidField

	if input.CourtId == "" {
		inv = append(inv, response.InvalidField{
			Field:    "court_id",
			Message:  "Court id is requried",
			Location: "body",
		})
	}
	if input.ScheduledAt == "" {
		inv = append(inv, response.InvalidField{
			Field:    "scheduled_at",
			Message:  "Scheduled at is required",
			Location: "body",
		})
	}
	if input.PlayerOneId == "" {
		inv = append(inv, response.InvalidField{
			Field:    "player_one_id",
			Message:  "Player one id is required",
			Location: "body",
		})
	}
	if input.PlayerTwoId == "" {
		inv = append(inv, response.InvalidField{
			Field:    "player_two_id",
			Message:  "Player two id is required",
			Location: "body",
		})
	}

	if len(inv) > 0 {
		return inv
	}
	return nil
}
