package me

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/internal/db"
	"github.com/markovidakovic/gdsi/server/pkg/response"
)

type store struct {
	db *db.Conn
}

func (s *store) queryMe(ctx context.Context, accountId string) (*MeModel, error) {
	query := `
		select 
			account.id as account_id,
			account.name as account_name,
			account.email as account_email,
			account.dob as account_dob,
			account.gender as account_gender,
			account.phone_number as account_phone_number,
			account.created_at as account_created_at,
			player.id as player_id,
			player.height as player_height,
			player.weight as player_weight,
			player.handedness as player_handedness,
			player.racket as player_racket,
			player.matches_expected as player_matches_expected,
			player.matches_played as player_matches_played,
			player.matches_won as player_matches_won,
			player.matches_scheduled as player_matches_scheduled,
			player.seasons_played as player_seasons_played,
			player.winning_ratio as player_winning_ratio,
			player.activity_ratio as player_activity_ratio,
			player.ranking as player_ranking,
			player.elo as player_elo,
			player.created_at as player_created_at,
			league.id as league_id,
			league.title as league_title,
			league.created_at as league_created_at
		FROM account
		left join player on account.id = player.account_id
		left join league on player.current_league_id = league.id
		where account.id = $1
	`

	var mm MeModel
	mm.PlayerProfile = PlayerProfileModel{}

	var leagueId sql.NullString
	var leagueTitle sql.NullString
	var leagueCreatedAt sql.NullTime

	err := s.db.QueryRow(ctx, query, accountId).Scan(
		&mm.Id,
		&mm.Name,
		&mm.Email,
		&mm.Dob,
		&mm.Gender,
		&mm.PhoneNumber,
		&mm.CreatedAt,
		&mm.PlayerProfile.Id,
		&mm.PlayerProfile.Height,
		&mm.PlayerProfile.Weight,
		&mm.PlayerProfile.Handedness,
		&mm.PlayerProfile.Racket,
		&mm.PlayerProfile.MatchesExpected,
		&mm.PlayerProfile.MatchesPlayed,
		&mm.PlayerProfile.MatchesWon,
		&mm.PlayerProfile.MatchesScheduled,
		&mm.PlayerProfile.SeasonsPlayed,
		&mm.PlayerProfile.WinningRation,
		&mm.PlayerProfile.ActivityRatio,
		&mm.PlayerProfile.Ranking,
		&mm.PlayerProfile.Elo,
		&mm.PlayerProfile.CreatedAt,
		&leagueId,
		&leagueTitle,
		&leagueCreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, response.ErrNotFound
		}
		return nil, err
	}

	if !leagueId.Valid {
		mm.PlayerProfile.CurrentLeague = nil
	} else {
		mm.PlayerProfile.CurrentLeague = &CurrentLeagueModel{
			Id:        leagueId.String,
			Title:     leagueTitle.String,
			CreatedAt: leagueCreatedAt.Time,
		}
	}

	return &mm, nil
}

func newStore(db *db.Conn) *store {
	return &store{
		db,
	}
}
