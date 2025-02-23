// Package validation provides a fluent api for validating domain entities and relationships.
//
// It uses a builder pattern to enable chainable validation calls while maintaining
// fail-fast behavior for critical errors. The package distinguishes between validation
// errors (invalid input) and system errors (e.g. database failures) and it also handles
// error accumulation logic.
package validation

import (
	"context"
	"fmt"

	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/failure"
)

type Validator struct {
	db *db.Conn
}

func NewValidator(db *db.Conn) *Validator {
	return &Validator{db: db}
}

// ValidationResult holds invalid fields and a possible internal failure
// when quering the db
type ValidationResult struct {
	invalidFields []failure.InvalidField
	failure       *failure.Failure
}

func (vr *ValidationResult) addInvalFld(field, message, source string) {
	vr.invalidFields = append(vr.invalidFields, failure.InvalidField{
		Field:    field,
		Message:  message,
		Location: source,
	})
}

func (vr *ValidationResult) result() error {
	if vr.failure != nil {
		return vr.failure
	}
	if len(vr.invalidFields) > 0 {
		return failure.NewValidation("invalid request parameters", vr.invalidFields)
	}
	return nil
}

func (v *Validator) courtExists(ctx context.Context, courtId string, source string) *ValidationResult {
	vr := &ValidationResult{}

	sql := `select exists(select 1 from court where id = $1)`
	var exists bool

	if err := v.db.QueryRow(ctx, sql, courtId).Scan(&exists); err != nil {
		vr.failure = failure.New("checking court existance", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
		return vr
	}

	if !exists {
		vr.addInvalFld("courtId", "court not found", source)
	}

	return vr
}

func (v *Validator) seasonExists(ctx context.Context, seasonId string, source string) *ValidationResult {
	vr := &ValidationResult{}

	sql := `select exists(select 1 from season where id = $1)`
	var exists bool

	if err := v.db.QueryRow(ctx, sql, seasonId).Scan(&exists); err != nil {
		vr.failure = failure.New("checking season existance", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
		return vr
	}

	if !exists {
		vr.addInvalFld("seasonId", "season not found", source)
	}

	return vr
}

func (v *Validator) leagueExists(ctx context.Context, leagueId string, source string) *ValidationResult {
	vr := &ValidationResult{}

	sql := `select exists(select 1 from league where id = $1)`
	var exists bool

	if err := v.db.QueryRow(ctx, sql, leagueId).Scan(&exists); err != nil {
		vr.failure = failure.New("checking league existance", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
		return vr
	}
	if !exists {
		vr.addInvalFld("leagueId", "league not found", source)
	}

	return vr
}

func (v *Validator) leagueInSeason(ctx context.Context, seasonId, leagueId string, source string) *ValidationResult {
	vr := &ValidationResult{}

	sql := `select exists(select 1 from league where id = $1 and season_id = $2)`
	var exists bool

	if err := v.db.QueryRow(ctx, sql, leagueId, seasonId).Scan(&exists); err != nil {
		vr.failure = failure.New("checking if league part of season", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
		return vr
	}
	if !exists {
		vr.addInvalFld("leagueId", "league not in season", source)
	}
	return vr
}

func (v *Validator) playerExists(ctx context.Context, playerId string, source string) *ValidationResult {
	vr := &ValidationResult{}

	sql := `select exists(select 1 from player where id = $1)`
	var exists bool

	if err := v.db.QueryRow(ctx, sql, playerId).Scan(&exists); err != nil {
		vr.failure = failure.New("checking player existance", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
		return vr
	}
	if !exists {
		vr.addInvalFld("playerId", "player not found", source)
	}
	return vr
}

// todo: later should be refactored to support a slice of player ids
func (v *Validator) playersInLeague(ctx context.Context, leagueId, playerOneId, playerTwoId string, source string) *ValidationResult {
	vr := &ValidationResult{}

	sql := `
		select exists (
			select 1 from player
			where id in ($1, $2) and current_league_id = $3
			having count(*) = 2
		)
	`
	var exists bool
	if err := v.db.QueryRow(ctx, sql, playerOneId, playerTwoId, leagueId).Scan(&exists); err != nil {
		vr.failure = failure.New("checking if player part of league", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
		return vr
	}
	if !exists {
		vr.addInvalFld("playerId", "players not in same league", source)
	}
	return vr
}

// matchExistsForPlayers checks if the two provided players have already played a match
// within the context of a season->league
// used for post and put match endpoints
func (v *Validator) matchExistsBetweenPlayers(ctx context.Context, seasonId, leagueId, playerOneId, playerTwoId string, source string) *ValidationResult {
	vr := &ValidationResult{}

	sql := `
		select exists (
			select 1 from match
			where season_id = $1 and league_id = $2 and player_one_id = $3 and player_two_Id = $4
		)
	`
	var exists bool
	if err := v.db.QueryRow(ctx, sql, seasonId, leagueId, playerOneId, playerTwoId).Scan(&exists); err != nil {
		vr.failure = failure.New("checking if players already played a match", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
		return vr
	}
	if exists {
		vr.addInvalFld("playerId", "match between players already exists", source)
	}
	return vr
}

// matchScheduledAtInSeason checks if the provided request parameter scheduled_at is
// within the seasons start_date - end_date range
// this is used for post and put match endpoints
func (v *Validator) matchScheduledInSeason(ctx context.Context, seasonId, scheduledAt string, source string) *ValidationResult {
	vr := &ValidationResult{}

	sql := `
		select exists (
			select 1 from season
			where id = $1 and $2 between start_date and end_date
		)
	`
	var exists bool
	if err := v.db.QueryRow(ctx, sql, seasonId, scheduledAt).Scan(&exists); err != nil {
		vr.failure = failure.New("checking if scheduled at in season playing range", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
		return vr
	}
	if !exists {
		vr.addInvalFld("scheduled_at", "must be within season date range", source)
	}
	return vr
}

type ValidationBuilder struct {
	validator *Validator
	ctx       context.Context
	result    *ValidationResult
}

func (v *Validator) NewValidation(ctx context.Context) *ValidationBuilder {
	return &ValidationBuilder{
		validator: v,
		ctx:       ctx,
		result:    &ValidationResult{},
	}
}

func (vb *ValidationBuilder) CourtExists(courtId, source string) *ValidationBuilder {
	// fail-fast if there is a critical error in prev validation
	if vb.result.failure != nil {
		return vb
	}

	vr := vb.validator.courtExists(vb.ctx, courtId, source)
	if vr.failure != nil {
		vb.result.failure = vr.failure
		return vb
	}
	vb.result.invalidFields = append(vb.result.invalidFields, vr.invalidFields...)
	return vb
}

func (vb *ValidationBuilder) SeasonExists(seasonId, source string) *ValidationBuilder {
	// fail-fast if there is a critical error in prev validation
	if vb.result.failure != nil {
		return vb
	}

	vr := vb.validator.seasonExists(vb.ctx, seasonId, source)
	if vr.failure != nil {
		vb.result.failure = vr.failure
		return vb
	}
	vb.result.invalidFields = append(vb.result.invalidFields, vr.invalidFields...)
	return vb
}

func (vb *ValidationBuilder) LeagueExists(leagueId, source string) *ValidationBuilder {
	// fail-fast if there is a critical error in prev validation
	if vb.result.failure != nil {
		return vb
	}

	vr := vb.validator.leagueExists(vb.ctx, leagueId, source)
	if vr.failure != nil {
		vb.result.failure = vr.failure
		return vb
	}
	vb.result.invalidFields = append(vb.result.invalidFields, vr.invalidFields...)
	return vb
}

func (vb *ValidationBuilder) LeagueInSeason(seasonId, leagueId, source string) *ValidationBuilder {
	// fal-fast if there is a critical error in prev validation
	if vb.result.failure != nil {
		return vb
	}

	vr := vb.validator.leagueInSeason(vb.ctx, seasonId, leagueId, source)
	if vr.failure != nil {
		vb.result.failure = vr.failure
		return vb
	}
	vb.result.invalidFields = append(vb.result.invalidFields, vr.invalidFields...)
	return vb
}

func (vb *ValidationBuilder) PlayerExists(playerId, source string) *ValidationBuilder {
	// fail fast if there is a critical error in prev validation
	if vb.result.failure != nil {
		return vb
	}

	vr := vb.validator.playerExists(vb.ctx, playerId, source)
	if vr.failure != nil {
		vb.result.failure = vr.failure
		return vb
	}
	vb.result.invalidFields = append(vb.result.invalidFields, vr.invalidFields...)
	return vb
}

func (vb *ValidationBuilder) PlayersInLeague(leagueId, playerOneId, playerTwoId, source string) *ValidationBuilder {
	if vb.result.failure != nil {
		return vb
	}

	vr := vb.validator.playersInLeague(vb.ctx, leagueId, playerOneId, playerTwoId, source)
	if vr.failure != nil {
		vb.result.failure = vr.failure
		return vb
	}
	vb.result.invalidFields = append(vb.result.invalidFields, vr.invalidFields...)
	return vb
}

func (vb *ValidationBuilder) MatchExistsBetweenPlayers(seasonId, leagueId, playerOneId, playerTwoId string, source string) *ValidationBuilder {
	if vb.result.failure != nil {
		return vb
	}
	vr := vb.validator.matchExistsBetweenPlayers(vb.ctx, seasonId, leagueId, playerOneId, playerTwoId, source)
	if vr.failure != nil {
		vb.result.failure = vr.failure
		return vb
	}
	vb.result.invalidFields = append(vb.result.invalidFields, vr.invalidFields...)
	return vb
}

func (vb *ValidationBuilder) MatchScheduledInSeason(seasonId, scheduledAt string, source string) *ValidationBuilder {
	if vb.result.failure != nil {
		return vb
	}
	vr := vb.validator.matchScheduledInSeason(vb.ctx, seasonId, scheduledAt, source)
	if vr.failure != nil {
		vb.result.failure = vr.failure
		return vb
	}
	vb.result.invalidFields = append(vb.result.invalidFields, vr.invalidFields...)
	return vb
}

func (vb *ValidationBuilder) Result() error {
	return vb.result.result()
}
