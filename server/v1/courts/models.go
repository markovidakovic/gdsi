package courts

import (
	"time"

	"github.com/markovidakovic/gdsi/server/response"
)

type CourtModel struct {
	Id        string       `json:"id"`
	Name      string       `json:"name"`
	Creator   CreatorModel `json:"creator"`
	CreatedAt time.Time    `json:"created_at"`
}

type CreatorModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CreateCourtRequestModel struct {
	Name      string `json:"name"`
	CreatorId string `json:"-"`
}

func (m CreateCourtRequestModel) Validate() []response.InvalidField {
	var inv []response.InvalidField

	if m.Name == "" {
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

type UpdateCourtRequestModel struct {
	Name string `json:"name"`
}

func (m UpdateCourtRequestModel) Validate() []response.InvalidField {
	var inv []response.InvalidField

	if m.Name == "" {
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
