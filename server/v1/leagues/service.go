package leagues

import (
	"context"

	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/params"
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
	err := s.validator.NewValidation(ctx).SeasonExists(model.SeasonId, "path").Result()
	if err != nil {
		return nil, err
	}

	lm, err := s.store.insertLeague(ctx, nil, model.Title, model.Description, model.CreatorId, model.SeasonId)
	if err != nil {
		return nil, err
	}

	return &lm, nil
}

func (s *service) processFindLeagues(ctx context.Context, seasonId string, query *params.Query) ([]LeagueModel, int, error) {
	err := s.validator.NewValidation(ctx).SeasonExists(seasonId, "path").Result()
	if err != nil {
		return nil, 0, err
	}

	count, err := s.store.countLeagues(ctx, seasonId)
	if err != nil {
		return nil, 0, failure.New("unable to find leagues", err)
	}

	limit, offset := query.CalcLimitAndOffset(count)

	lms, err := s.store.findLeagues(ctx, seasonId, limit, offset, query.OrderBy)
	if err != nil {
		return nil, 0, err
	}

	return lms, count, err
}

func (s *service) processFindLeague(ctx context.Context, seasonId, leagueId string) (*LeagueModel, error) {
	err := s.validator.NewValidation(ctx).SeasonExists(seasonId, "path").Result()
	if err != nil {
		return nil, err
	}

	lm, err := s.store.findLeague(ctx, seasonId, leagueId)
	if err != nil {
		return nil, err
	}

	return lm, nil
}

func (s *service) processUpdateLeague(ctx context.Context, model UpdateLeagueRequestModel) (*LeagueModel, error) {
	err := s.validator.NewValidation(ctx).SeasonExists(model.SeasonId, "path").Result()
	if err != nil {
		return nil, err
	}

	lm, err := s.store.updateLeague(ctx, nil, model.Title, model.Description, model.SeasonId, model.LeagueId)
	if err != nil {
		return nil, err
	}

	return &lm, nil
}

func (s *service) processDeleteLeague(ctx context.Context, seasonId, leagueId string) error {
	err := s.validator.NewValidation(ctx).SeasonExists(seasonId, "path").Result()
	if err != nil {
		return err
	}

	err = s.store.deleteLeague(ctx, nil, seasonId, leagueId)
	if err != nil {
		return err
	}

	return nil
}
