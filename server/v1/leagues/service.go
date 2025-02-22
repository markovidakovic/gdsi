package leagues

import (
	"context"

	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/validation"
)

type service struct {
	cfg       *config.Config
	store     *store
	validator *validation.Validator
}

func newService(cfg *config.Config, store *store, validator *validation.Validator) *service {
	return &service{
		cfg,
		store,
		validator,
	}
}

func (s *service) processCreateLeague(ctx context.Context, model CreateLeagueRequestModel) (*LeagueModel, error) {
	// validation
	err := s.validator.NewValidation(ctx).SeasonExists(model.SeasonId, "path").Result()
	if err != nil {
		return nil, err
	}

	// create league
	lm, err := s.store.insertLeague(ctx, nil, model.Title, model.Description, model.CreatorId, model.SeasonId)
	if err != nil {
		return nil, err
	}

	return &lm, nil
}

func (s *service) processFindLeagues(ctx context.Context, seasonId string) ([]LeagueModel, error) {
	// validation
	err := s.validator.NewValidation(ctx).SeasonExists(seasonId, "path").Result()
	if err != nil {
		return nil, err
	}

	// find leagues
	lms, err := s.store.findLeagues(ctx, seasonId)
	if err != nil {
		return nil, err
	}

	return lms, err
}

func (s *service) processFindLeague(ctx context.Context, seasonId, leagueId string) (*LeagueModel, error) {
	// validation
	err := s.validator.NewValidation(ctx).SeasonExists(seasonId, "path").Result()
	if err != nil {
		return nil, err
	}

	// find league
	lm, err := s.store.findLeague(ctx, seasonId, leagueId)
	if err != nil {
		return nil, err
	}

	return lm, nil
}

func (s *service) processUpdateLeague(ctx context.Context, model UpdateLeagueRequestModel) (*LeagueModel, error) {
	// validation
	err := s.validator.NewValidation(ctx).SeasonExists(model.SeasonId, "path").Result()
	if err != nil {
		return nil, err
	}

	// update league
	lm, err := s.store.updateLeague(ctx, nil, model.Title, model.Description, model.SeasonId, model.LeagueId)
	if err != nil {
		return nil, err
	}

	return &lm, nil
}

func (s *service) processDeleteLeague(ctx context.Context, seasonId, leagueId string) error {
	// validation
	err := s.validator.NewValidation(ctx).SeasonExists(seasonId, "path").Result()
	if err != nil {
		return err
	}

	// delete league
	err = s.store.deleteLeague(ctx, nil, seasonId, leagueId)
	if err != nil {
		return err
	}

	return nil
}
