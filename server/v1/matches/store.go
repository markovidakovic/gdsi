package matches

import (
	"context"

	"github.com/markovidakovic/gdsi/server/db"
)

type store struct {
	db *db.Conn
}

func newStore(db *db.Conn) *store {
	return &store{
		db,
	}
}

func (s *store) insertMatch(ctx context.Context, input CreateMatchRequestModel) (MatchModel, error) {
	sql1 := `
		with inserted_match as (
			insert into match (court_id, scheduled_at, player_one_id, player_two_id, winner_id, score, season_id, league_id)
			values ($1, $2, $3, $4, $5, $6, $7, $8)
			returning id, court_id, scheduled_at, player_one_id, player_two_id, winner_id, score, season_id, league_id, created_at
		)
		select
			im.id
			court.id as court_id,
			court.name as court_name,
			player1.id as player_one_id,
			account1.name as player_one_name,
			player2.id as player_two_id,
			account2.name as player_two_name,
			winner.id as winner_id,
			account3.name as winner_name,
			im.score,
			season.id as season_id,
			season.title as season_title,
			league.id as league_id,
			league.title as league_title,
			im.created_at
		join court on im.court_id = court.id
		join player player1 on im.player_one_id = player1.id
		join account account1 on player1.account_id = account1.id
		join player player2 on im.player_two_id = player2.id
		join account account2 on player2.account_id = account2.id
		left join player winner on im.winner_id = winner.id
		left join account account3 on winner.account_id = account3.id
		join season on im.season_id = season.id
		join league on im.league_id = league.id
	`

	var dest MatchModel

	err := s.db.QueryRow(ctx, sql1, input.CourtId, input.ScheduledAt, input.PlayerOneId, input.PlayerTwoId, input.WinnerId, input.Score, input.SeasonId, input.LeagueId).Scan(
		&dest.Id,
		&dest.Court.Id,
		&dest.Court.Name,
		&dest.PlayerOne.Id,
		&dest.PlayerOne.Name,
		&dest.PlayerTwo.Id,
		&dest.PlayerTwo.Name,
		&dest.Winner.Id,
		&dest.Winner.Name,
		&dest.Score,
		&dest.Season.Id,
		&dest.Season.Title,
		&dest.League.Id,
		&dest.League.Title,
		&dest.CreatedAt,
	)
	if err != nil {
		return dest, err
	}

	return dest, nil
}

// validateInsertMatch takes the season id, league id, player ids and checks if: season, league, players exists and if
// league is part of the season and players are part of the league
func (s *store) validateInsertMatch(ctx context.Context, seasonId, leagueId, player1Id, player2Id string) (seasonExists bool, leagueExists bool, leagueInSeason bool, playerOneExists bool, playerTwoExists bool, playersInLeague bool, err error) {
	sql1 := `
		select
			exists ( 
				select 1 from season where id = $1 
			) as season_exists,
			exists (
				select 1 from league where id = $2
			) as league_exists,
			exists (
				select 1 from league where id = $2 and season_id = $1
				) as league_in_season,
			exists (
				select 1 from player where id = $3
			) as player_one_exists,
			exists (
				select 1 from player where id = $4
			) as player_two_exists,
			exists (
				select 1 from player
				where id in ($3, $4)
				and current_league_id = $2
				having count(*) = 2
			) as players_in_league
	`

	err = s.db.QueryRow(ctx, sql1, seasonId, leagueId, player1Id, player2Id).Scan(&seasonExists, &leagueExists, &leagueInSeason, &playerOneExists, &playerTwoExists, &playersInLeague)
	if err != nil {
		return
	}

	return
}
