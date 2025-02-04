package me

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/response"
)

type store struct {
	db *db.Conn
}

func newStore(db *db.Conn) *store {
	return &store{
		db,
	}
}

func (s *store) findMe(ctx context.Context, accountId string) (*MeModel, error) {
	sql1 := `
		select 
			account.id as account_id,
			account.name as account_name,
			account.email as account_email,
			account.dob as account_dob,
			account.gender as account_gender,
			account.phone_number as account_phone_number,
			account.role as account_role,
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
		from account
		left join player on account.id = player.account_id
		left join league on player.current_league_id = league.id
		where account.id = $1
	`

	var mm MeModel
	mm.PlayerProfile = PlayerProfileModel{}

	var leagueId, leagueTitle sql.NullString
	var leagueCreatedAt sql.NullTime

	err := s.db.QueryRow(ctx, sql1, accountId).Scan(
		&mm.Id,
		&mm.Name,
		&mm.Email,
		&mm.Dob,
		&mm.Gender,
		&mm.PhoneNumber,
		&mm.Role,
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

func (s *store) updateMe(ctx context.Context, accountId string, input UpdateMeRequestModel) (*MeModel, error) {
	var dest MeModel
	dest.PlayerProfile = PlayerProfileModel{}
	var leagueId, leagueTitle sql.NullString
	var leagueCreatedAt sql.NullTime

	sql1 := `
		with updated_account as (
			update account 
			set name = $1
			where id = $2
			returning id, name, email, dob, gender, phone_number, role, created_at
		)
		select 
			ua.id as account_id,
			ua.name as account_name,
			ua.email as account_email,
			ua.dob as account_dob,
			ua.gender as account_gender,
			ua.phone_number as account_phone_number,
			ua.role as account_role,
			ua.created_at as account_created_at,
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
		from updated_account ua
		left join player on ua.id = player.account_id
		left join league on player.current_league_id = league.id
	`

	err := s.db.QueryRow(ctx, sql1, input.Name, accountId).Scan(
		&dest.Id,
		&dest.Name,
		&dest.Email,
		&dest.Dob,
		&dest.Gender,
		&dest.PhoneNumber,
		&dest.Role,
		&dest.CreatedAt,
		&dest.PlayerProfile.Id,
		&dest.PlayerProfile.Height,
		&dest.PlayerProfile.Weight,
		&dest.PlayerProfile.Handedness,
		&dest.PlayerProfile.Racket,
		&dest.PlayerProfile.MatchesExpected,
		&dest.PlayerProfile.MatchesPlayed,
		&dest.PlayerProfile.MatchesWon,
		&dest.PlayerProfile.MatchesScheduled,
		&dest.PlayerProfile.SeasonsPlayed,
		&dest.PlayerProfile.WinningRation,
		&dest.PlayerProfile.ActivityRatio,
		&dest.PlayerProfile.Ranking,
		&dest.PlayerProfile.Elo,
		&dest.PlayerProfile.CreatedAt,
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
		dest.PlayerProfile.CurrentLeague = nil
	} else {
		dest.PlayerProfile.CurrentLeague = &CurrentLeagueModel{
			Id:        leagueId.String,
			Title:     leagueTitle.String,
			CreatedAt: leagueCreatedAt.Time,
		}
	}

	return &dest, nil
}
