package matches

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/response"
)

type service struct {
	cfg   *config.Config
	store *store
}

func newService(cfg *config.Config, store *store) *service {
	return &service{
		cfg,
		store,
	}
}

func (s *service) processCreateMatch(ctx context.Context, input CreateMatchRequestModel) (*MatchModel, error) {
	// validate params in the db layer
	courtExists, seasonExists, leagueExists, playerOneExists, playerTwoExists, playersInLeague, err := s.store.validateInsertUpdateMatch(ctx, input.CourtId, input.SeasonId, input.LeagueId, input.PlayerOneId, input.PlayerTwoId)
	if err != nil {
		return nil, err
	}

	if !courtExists {
		return nil, fmt.Errorf("finding court: %w", response.ErrNotFound)
	}
	if !seasonExists {
		return nil, fmt.Errorf("finding season: %w", response.ErrNotFound)
	}
	if !leagueExists {
		return nil, fmt.Errorf("finding league: %w", response.ErrNotFound)
	}
	if !playerOneExists {
		return nil, fmt.Errorf("finding player one: %w", response.ErrNotFound)
	}
	if !playerTwoExists {
		return nil, fmt.Errorf("finding player two: %w", response.ErrNotFound)
	}
	if !playersInLeague {
		return nil, fmt.Errorf("players not in league: %w", response.ErrBadRequest)
	}

	if input.Score != nil {
		fmt.Println("analize the score and set the winner id")
	} else {
		input.WinnerId = nil
	}

	// insert match
	cm, err := s.store.insertMatch(ctx, input)
	if err != nil {
		return nil, err
	}

	return &cm, nil
}

func (s *service) processGetMatches(ctx context.Context, seasonId, leagueId string) ([]MatchModel, error) {
	// validate params in the db layer
	seasonExists, leagueExists, err := s.store.validateFindMatches(ctx, seasonId, leagueId)
	if err != nil {
		return nil, err
	}

	if !seasonExists {
		return nil, fmt.Errorf("finding season: %w", response.ErrNotFound)
	}
	if !leagueExists {
		return nil, fmt.Errorf("finding league: %w", response.ErrNotFound)
	}

	// call the store
	mms, err := s.store.findMatches(ctx, seasonId, leagueId)
	if err != nil {
		return nil, err
	}

	return mms, nil
}

func (s *service) processGetMatch(ctx context.Context, seasonId, leagueId, matchId string) (*MatchModel, error) {
	// validate params in the db layer
	seasonExists, leagueExists, err := s.store.validateFindMatches(ctx, seasonId, leagueId)
	if err != nil {
		return nil, err
	}

	if !seasonExists {
		return nil, fmt.Errorf("finding season: %w", response.ErrNotFound)
	}
	if !leagueExists {
		return nil, fmt.Errorf("finding league: %w", response.ErrNotFound)
	}

	mm, err := s.store.findMatch(ctx, seasonId, leagueId, matchId)
	if err != nil {
		return nil, err
	}

	return mm, nil
}

func (s *service) processUpdateMatch(ctx context.Context, input UpdateMatchRequestModel) (*MatchModel, error) {
	// check if able to modify match
	hasScore, err := s.store.checkMatchScore(ctx, input.MatchId)
	if err != nil {
		return nil, err
	}

	if hasScore {
		return nil, fmt.Errorf("%w: not able to modify a match that has a score", response.ErrConflict)
	}

	// validate params in the db layer
	courtExists, seasonExists, leagueExists, playerOneExists, playerTwoExists, playersInLeague, err := s.store.validateInsertUpdateMatch(ctx, input.CourtId, input.SeasonId, input.LeagueId, input.PlayerOneId, input.PlayerTwoId)
	if err != nil {
		return nil, err
	}

	if !courtExists {
		return nil, fmt.Errorf("finding court: %w", response.ErrNotFound)
	}
	if !seasonExists {
		return nil, fmt.Errorf("finding season: %w", response.ErrNotFound)
	}
	if !leagueExists {
		return nil, fmt.Errorf("finding league: %w", response.ErrNotFound)
	}
	if !playerOneExists {
		return nil, fmt.Errorf("finding player one: %w", response.ErrNotFound)
	}
	if !playerTwoExists {
		return nil, fmt.Errorf("finding player two: %w", response.ErrNotFound)
	}
	if !playersInLeague {
		return nil, fmt.Errorf("players not in league: %w", response.ErrBadRequest)
	}

	mm, err := s.store.updateMatch(ctx, input)
	if err != nil {
		return nil, err
	}

	return mm, nil
}

func (s *service) processSubmitMatchScore(ctx context.Context, input SubmitMatchScoreRequestModel) (*MatchModel, error) {
	seasonExists, leagueExists, matchExists, err := s.store.validateSubmitMatchScore(ctx, input.SeasonId, input.LeagueId, input.MatchId)
	if err != nil {
		return nil, err
	}
	if !seasonExists {
		return nil, fmt.Errorf("finding season: %w", response.ErrNotFound)
	}
	if !leagueExists {
		return nil, fmt.Errorf("finding league: %w", response.ErrNotFound)
	}
	if !matchExists {
		return nil, fmt.Errorf("finding match: %w", response.ErrNotFound)
	}

	// find the existing match
	match, err := s.store.findMatch(ctx, input.SeasonId, input.LeagueId, input.MatchId)
	if err != nil {
		return nil, err
	}

	// check if able to submit result
	if match.Score != nil {
		return nil, fmt.Errorf("%w: not able to submit match result, score already exists", response.ErrConflict)
	}

	// determine the winner
	sets := strings.Split(input.Score, ",")
	var pl1SetsWon int

	for _, set := range sets {
		setSl := strings.Split(set, "-")
		pl1Score, _ := strconv.Atoi(setSl[0])
		pl2Score, _ := strconv.Atoi(setSl[1])

		if pl1Score > pl2Score {
			pl1SetsWon++
		}
	}

	if pl1SetsWon == 2 {
		input.WinnerId = match.PlayerOne.Id
	} else {
		input.WinnerId = match.PlayerTwo.Id
	}

	// also add the match p1 and p2 info to the input struct
	input.PlayerOneId = match.PlayerOne.Id
	input.PlayerTwoId = match.PlayerTwo.Id

	// call the store
	result, err := s.store.insertMatchScore(ctx, input)
	if err != nil {
		return nil, err
	}

	return result, nil
}
