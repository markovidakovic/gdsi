package standings

import "time"

type Standing struct {
	Id            string    `json:"id"`
	Points        int       `json:"points"`
	MatchesPlayed int       `json:"matches_played"`
	MatchesWon    int       `json:"matches_won"`
	SetsWon       int       `json:"sets_won"`
	SetsLost      int       `json:"sets_lost"`
	GamesWon      int       `json:"games_won"`
	GamesLost     int       `json:"games_lost"`
	SeasonId      string    `json:"season_id"`
	LeagueId      string    `json:"league_id"`
	PlayerId      string    `json:"player_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type SeasonModel struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type LeagueModel struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type PlayerModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type StandingModel struct {
	Id            string      `json:"id"`
	Points        int         `json:"points"`
	MatchesPlayed int         `json:"matches_played"`
	MatchesWon    int         `json:"matches_won"`
	SetsWon       int         `json:"sets_won"`
	SetsLost      int         `json:"sets_lost"`
	GamesWon      int         `json:"games_won"`
	GamesLost     int         `json:"games_lost"`
	Season        SeasonModel `json:"season"`
	League        LeagueModel `json:"league"`
	Player        PlayerModel `json:"player"`
	CreatedAt     time.Time   `json:"created_at"`
}
