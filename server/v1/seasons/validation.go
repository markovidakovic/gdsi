package seasons

import (
	"fmt"

	"github.com/markovidakovic/gdsi/server/response"
)

func validatePostSeason(input CreateSeasonRequestModel) []response.InvalidField {
	var inv []response.InvalidField
	fmt.Printf("input: %v\n", input)

	if len(inv) > 0 {
		return inv
	}

	return nil
}
