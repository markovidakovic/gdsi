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

func (vr *ValidationResult) addInvalFld(field, message, location string) {
	vr.invalidFields = append(vr.invalidFields, failure.InvalidField{
		Field:    field,
		Message:  message,
		Location: location,
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

func (v *Validator) courtExists(ctx context.Context, courtId, location string) *ValidationResult {
	vr := &ValidationResult{}

	sql := `select exists(select 1 from court where id = $1)`
	var exists bool

	if err := v.db.QueryRow(ctx, sql, courtId).Scan(&exists); err != nil {
		vr.failure = failure.New("checking court existance", fmt.Errorf("%w: %v", failure.ErrInternal, err))
		return vr
	}

	if !exists {
		vr.addInvalFld("courtId", "court not found", location)
	}

	return vr
}

func (v *Validator) seasonExists(ctx context.Context, seasonId, location string) *ValidationResult {
	vr := &ValidationResult{}

	sql := `select exists(select 1 from season where id = $1)`
	var exists bool

	if err := v.db.QueryRow(ctx, sql, seasonId).Scan(&exists); err != nil {
		vr.failure = failure.New("checking season existance", fmt.Errorf("%w: %v", failure.ErrInternal, err))
		return vr
	}

	if !exists {
		vr.addInvalFld("seasonId", "season not found", location)
	}

	return vr
}

func (v *Validator) leagueExists(ctx context.Context, leagueId, location string) *ValidationResult {
	vr := &ValidationResult{}

	sql := `select exists(select 1 from season where id = $1)`
	var exists bool

	if err := v.db.QueryRow(ctx, sql, leagueId).Scan(&exists); err != nil {
		vr.failure = failure.New("checking league existance", fmt.Errorf("%w: %v", failure.ErrInternal, err))
		return vr
	}
	if !exists {
		vr.addInvalFld("leagueId", "league not found", location)
	}
	return vr
}

func (v *Validator) leagueInSeason(ctx context.Context, seasonId, leagueId, location string) *ValidationResult {
	vr := &ValidationResult{}

	sql := `select exists(select 1 from league where id = $1 and season_id = $2)`
	var exists bool

	if err := v.db.QueryRow(ctx, sql, leagueId, seasonId).Scan(&exists); err != nil {
		vr.failure = failure.New("checking if league part of season", fmt.Errorf("%w: %v", failure.ErrInternal, err))
		return vr
	}
	if !exists {
		vr.addInvalFld("leagueId", "league not in season", location)
	}
	return vr
}

func (v *Validator) playerExists(ctx context.Context, playerId, location string) *ValidationResult {
	vr := &ValidationResult{}

	sql := `select exists(select 1 from player where id = $1)`
	var exists bool

	if err := v.db.QueryRow(ctx, sql, playerId).Scan(&exists); err != nil {
		vr.failure = failure.New("checking player existance", fmt.Errorf("%w: %v", failure.ErrInternal, err))
		return vr
	}
	if !exists {
		vr.addInvalFld("playerId", "player not found", location)
	}
	return vr
}

// todo: later should be refactored to support a slice of player ids
func (v *Validator) playersInLeague(ctx context.Context, leagueId, playerOneId, playerTwoId, location string) *ValidationResult {
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
		vr.failure = failure.New("checking if player part of league", fmt.Errorf("%w: %v", failure.ErrInternal, err))
		return vr
	}
	if !exists {
		vr.addInvalFld("playerId", "players not in same league", location)
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

func (vb *ValidationBuilder) CourtExists(courtId, location string) *ValidationBuilder {
	// fail-fast if there is a critical error in prev validation
	if vb.result.failure != nil {
		return vb
	}

	vr := vb.validator.courtExists(vb.ctx, courtId, location)
	if vr.failure != nil {
		vb.result.failure = vr.failure
		return vb
	}
	vb.result.invalidFields = append(vb.result.invalidFields, vr.invalidFields...)
	return vb
}

func (vb *ValidationBuilder) SeasonExists(seasonId, location string) *ValidationBuilder {
	// fail-fast if there is a critical error in prev validation
	if vb.result.failure != nil {
		return vb
	}

	vr := vb.validator.seasonExists(vb.ctx, seasonId, location)
	if vr.failure != nil {
		vb.result.failure = vr.failure
		return vb
	}
	vb.result.invalidFields = append(vb.result.invalidFields, vr.invalidFields...)
	return vb
}

func (vb *ValidationBuilder) LeagueExists(leagueId, location string) *ValidationBuilder {
	// fail-fast if there is a critical error in prev validation
	if vb.result.failure != nil {
		return vb
	}

	vr := vb.validator.leagueExists(vb.ctx, leagueId, location)
	if vr.failure != nil {
		vb.result.failure = vr.failure
		return vb
	}
	vb.result.invalidFields = append(vb.result.invalidFields, vr.invalidFields...)
	return vb
}

func (vb *ValidationBuilder) LeagueInSeason(seasonId, leagueId, location string) *ValidationBuilder {
	// fal-fast if there is a critical error in prev validation
	if vb.result.failure != nil {
		return vb
	}

	vr := vb.validator.leagueInSeason(vb.ctx, seasonId, leagueId, location)
	if vr.failure != nil {
		vb.result.failure = vr.failure
		return vb
	}
	vb.result.invalidFields = append(vb.result.invalidFields, vr.invalidFields...)
	return vb
}

func (vb *ValidationBuilder) PlayerExists(playerId, location string) *ValidationBuilder {
	// fail fast if there is a critical error in prev validation
	if vb.result.failure != nil {
		return vb
	}

	vr := vb.validator.playerExists(vb.ctx, playerId, location)
	if vr.failure != nil {
		vb.result.failure = vr.failure
		return vb
	}
	vb.result.invalidFields = append(vb.result.invalidFields, vr.invalidFields...)
	return vb
}

func (vb *ValidationBuilder) PlayersInLeague(leagueId, playerOneId, playerTwoId, location string) *ValidationBuilder {
	if vb.result.failure != nil {
		return vb
	}

	vr := vb.validator.playersInLeague(vb.ctx, leagueId, playerOneId, playerTwoId, location)
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
