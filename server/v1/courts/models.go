package courts

import "time"

type Court struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatorId string    `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
}

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

type CreateCourtModel struct {
	Name      string `json:"name"`
	CreatorId string `json:"-"`
}

type UpdateCourtModel struct {
	Name string `json:"name"`
}
