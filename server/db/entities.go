package db

import (
	"database/sql"
	"time"
)

// db table account
type Account struct {
	Id          string
	Name        string
	Email       string
	Dob         time.Time
	Gender      string
	PhoneNumber string
	Password    string
	Role        string
	CreatedAt   time.Time
}

// db table refresh_token
type RefreshToken struct {
	Id         string
	AccountId  string // fk to account
	TokenHash  string
	DeviceId   sql.NullString
	IpAddress  sql.NullString
	UserAgent  sql.NullString
	IssuedAt   time.Time
	ExpiresAt  time.Time
	LastUsedAt sql.NullTime
	IsRevoked  bool
}

// db table court
type Court struct {
	Id        string
	Name      string
	CreatorId string // fk to account
	CreatedAt time.Time
}

// db table season
type Season struct {
	Id          string
	Title       string
	Description sql.NullString
	StartDate   time.Time
	EndDate     time.Time
	CreatorId   string // fk to account
	CreatedAt   string
}

// db table player
type Player struct {
	Id               string
	Height           sql.NullInt32
	Weight           sql.NullInt32
	Handedness       sql.NullString
	Racket           sql.NullString
	MatchesExpected  int // amount of expected matches played from each season
	MatchesPlayed    int // actual amount of matches played
	MatchesWon       int
	MatchesScheduled int // amount of matches created
	SeasonsPlayed    int
	AccountId        string // fk to account
	CurrentLeagueId  string // fk to league
	CreatedAt        time.Time
}

// db table match
type Match struct {
	Id          string
	CourtId     string // fk to court
	ScheduledAt time.Time
	PlayerOneId string         // fk to player
	PlayerTwoId string         // fk to player
	WinnerId    sql.NullString // fk to player
	Score       sql.NullString
	SeasonId    string // fk to season
	LeagueId    string // fk to league
	CreatorId   string // fk to player
	CreatedAt   time.Time
}

// db table standing
type Standing struct {
	Id            string
	Points        int
	MatchesPlayed int
	MatchesWon    int
	SetsWon       int
	SetsLost      int
	GamesWon      int
	GamesLost     int
	SeasonId      string // fk to season
	LeagueId      string // fk to league
	PlayerId      string // fk to player
	CreatedAt     time.Time
}
