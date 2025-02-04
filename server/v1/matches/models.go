package matches

import "time"

// db table
type Match struct {
	Id          string    `json:"id"`
	CourtId     string    `json:"court_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	PlayerOneId string    `json:"player_one_id"`
	PlayerTwoId string    `json:"player_two_id"`
	WinnerId    string    `json:"winner_id"`
	Score       string    `json:"score"`
	SeasonId    string    `json:"season_id"`
	LeagueId    string    `json:"league_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type MatchModel struct {
	Id          string      `json:"id"`
	Court       CourtModel  `json:"court"`
	ScheduledAt time.Time   `json:"scheduled_at"`
	PlayerOne   PlayerModel `json:"player_one"`
	PlayerTwo   PlayerModel `json:"player_two"`
	Winner      PlayerModel `json:"winner"`
	Score       string      `json:"score"`
	Season      SeasonModel `json:"season"`
	League      LeagueModel `json:"league"`
	CreatedAt   time.Time   `json:"created_at"`
}

type CourtModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type PlayerModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SeasonModel struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type LeagueModel struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type CreateMatchRequestModel struct {
	CourtId     string  `json:"court_id"`
	ScheduledAt string  `json:"scheduled_at"`
	PlayerOneId string  `json:"player_one_id"`
	PlayerTwoId string  `json:"player_two_id"`
	WinnerId    string  `json:"-"`
	Score       *string `json:"score"`
	SeasonId    string  `json:"-"`
	LeagueId    string  `json:"-"`
}

// todo: think about what is allowed on an update of a match
// should it be allowed?
type UpdateMatchRequestModel struct {
	CourtId     string  `json:"court_id"`
	ScheduledAt string  `json:"scheduled_at"`
	PlayerOneId string  `json:"player_one_id"`
	PlayerTwoId string  `json:"player_two_id"`
	Score       *string `json:"score"`
}
