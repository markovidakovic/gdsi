package matches

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/failure"
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

func (s *service) processCreateMatch(ctx context.Context, model CreateMatchRequestModel) (*MatchModel, error) {
	// validation
	err := s.validator.NewValidation(ctx).
		CourtExists(model.CourtId).
		SeasonExists(model.SeasonId).
		LeagueExists(model.LeagueId).LeagueInSeason(model.SeasonId, model.LeagueId).
		PlayerExists(model.PlayerOneId).PlayerExists(model.PlayerTwoId).PlayersInLeague(model.LeagueId, model.PlayerOneId, model.PlayerTwoId).
		Result()
	if err != nil {
		return nil, err
	}

	// begin tx
	tx, err := s.store.db.Begin(ctx)
	if err != nil {
		return nil, failure.New("unable to create a match", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err != pgx.ErrTxClosed {
			log.Printf("failed to rollback the create match tx: %v", err)
		}
	}()

	if model.Score != nil {
		winnerId := determineMatchWinner(*model.Score, model.PlayerOneId, model.PlayerTwoId)
		model.WinnerId = &winnerId
	}

	// insert match
	match, err := s.store.insertMatch(ctx, tx, model.CourtId, model.ScheduledAt, model.PlayerOneId, model.PlayerTwoId, model.WinnerId, model.Score, model.SeasonId, model.LeagueId)
	if err != nil {
		return nil, failure.New("unable to create a match", err)
	}

	// for cases where the score is submitted upon match creation
	if model.Score != nil {
		// calc player statistics
		pl1Stats := calcMatchStats(*model.Score, true)
		pl2Stats := calcMatchStats(*model.Score, false)

		// update player stats
		err = s.store.updatePlayerStatistics(ctx, tx, *model.WinnerId, model.PlayerOneId, model.PlayerTwoId)
		if err != nil {
			// the err here will be an internal and not a player not found because we've confirmend
			// in the validation that these two players exist
			return nil, failure.New("unable to create a match", err)
		}

		// update standings for pl1
		err = s.store.updateStanding(ctx, tx, model.SeasonId, model.LeagueId, model.PlayerOneId, pl1Stats)
		if err != nil {
			return nil, failure.New("unable to create a match", err)
		}

		// update standings for pl2
		err = s.store.updateStanding(ctx, tx, model.SeasonId, model.LeagueId, model.PlayerTwoId, pl2Stats)
		if err != nil {
			return nil, failure.New("unable to create a match", err)
		}
	}

	// increment matches scheduled for the player creating the match
	err = s.store.incrementPlayerMatchesScheduled(ctx, tx, model.PlayerOneId)
	if err != nil {
		return nil, failure.New("unable to create a match", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, failure.New("unable to create a match", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return &match, nil
}

func (s *service) processGetMatches(ctx context.Context, seasonId, leagueId string) ([]MatchModel, error) {
	// validation
	err := s.validator.NewValidation(ctx).
		SeasonExists(seasonId).
		LeagueExists(leagueId).LeagueInSeason(seasonId, leagueId).
		Result()
	if err != nil {
		return nil, err
	}

	// call the store
	mms, err := s.store.findMatches(ctx, seasonId, leagueId)
	if err != nil {
		return nil, err
	}

	return mms, nil
}

func (s *service) processGetMatch(ctx context.Context, seasonId, leagueId, matchId string) (*MatchModel, error) {
	// validate
	err := s.validator.NewValidation(ctx).
		SeasonExists(seasonId).
		LeagueExists(leagueId).LeagueInSeason(seasonId, leagueId).
		Result()
	if err != nil {
		return nil, err
	}

	mm, err := s.store.findMatch(ctx, seasonId, leagueId, matchId)
	if err != nil {
		return nil, err
	}

	return mm, nil
}

func (s *service) processUpdateMatch(ctx context.Context, model UpdateMatchRequestModel) (*MatchModel, error) {
	// validation
	err := s.validator.NewValidation(ctx).
		CourtExists(model.CourtId).
		SeasonExists(model.SeasonId).
		LeagueExists(model.LeagueId).LeagueInSeason(model.SeasonId, model.LeagueId).
		PlayerExists(model.PlayerOneId).PlayerExists(model.PlayerTwoId).PlayersInLeague(model.LeagueId, model.PlayerOneId, model.PlayerTwoId).
		Result()
	if err != nil {
		return nil, err
	}

	// check if able to modify match
	hasScore, err := s.store.checkMatchScore(ctx, model.MatchId)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("match for update not found", err)
		}
		return nil, failure.New("unable to update match", err)
	}

	if hasScore {
		return nil, failure.New("not able to modify a match that has a score", failure.ErrCantModify)
	}

	mm, err := s.store.updateMatch(ctx, nil, model.CourtId, model.ScheduledAt, model.PlayerTwoId, model.SeasonId, model.LeagueId, model.MatchId)
	if err != nil {
		return nil, err
	}

	return mm, nil
}

func (s *service) processSubmitMatchScore(ctx context.Context, model SubmitMatchScoreRequestModel) (*MatchModel, error) {
	// validation
	err := s.validator.NewValidation(ctx).
		SeasonExists(model.SeasonId).
		LeagueExists(model.LeagueId).LeagueInSeason(model.SeasonId, model.LeagueId).
		Result()
	if err != nil {
		return nil, err
	}

	// find the existing match
	match, err := s.store.findMatch(ctx, model.SeasonId, model.LeagueId, model.MatchId)
	if err != nil {
		return nil, err
	}

	// check if able to submit result
	if match.Score != nil {
		return nil, failure.New("not able to submit match result, score already exists", failure.ErrCantModify)
	}

	// also add the match p1 and p2 info to the model struct
	model.PlayerOneId = match.PlayerOne.Id
	model.PlayerTwoId = match.PlayerTwo.Id
	model.WinnerId = determineMatchWinner(model.Score, match.PlayerOne.Id, match.PlayerTwo.Id)

	// begin tx
	tx, err := s.store.db.Begin(ctx)
	if err != nil {
		// return nil, err
		return nil, failure.New("not able to submit match score", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err != pgx.ErrTxClosed {
			log.Printf("failed to rollback the submit match score tx: %v", err)
		}
	}()

	result, err := s.store.updateMatchScore(ctx, tx, model.SeasonId, model.LeagueId, model.MatchId, model.Score, model.WinnerId)
	if err != nil {
		// if we made it this far, the match exists and this error will be an internal error, so we format the message accordingly
		return nil, failure.New("not able to submit match score", err)
	}

	err = s.store.updatePlayerStatistics(ctx, tx, model.WinnerId, model.PlayerOneId, model.PlayerTwoId)
	if err != nil {
		// same as above, the winner, pl1 and pl2 are part of the match and the match exists
		return nil, failure.New("not able to submit match score", err)
	}

	// calc pl1 & pl2 match stats
	pl1MatchStats := calcMatchStats(model.Score, true)
	pl2MatchStats := calcMatchStats(model.Score, false)

	// update standing for pl1
	err = s.store.updateStanding(ctx, tx, model.SeasonId, model.LeagueId, model.PlayerOneId, pl1MatchStats)
	if err != nil {
		return nil, failure.New("not able to submit match score", err)
	}

	// update standing for pl2
	err = s.store.updateStanding(ctx, tx, model.SeasonId, model.LeagueId, model.PlayerTwoId, pl2MatchStats)
	if err != nil {
		return nil, failure.New("not able to submit match score", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, failure.New("not able to submit match score", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return result, nil
}

// determineMatchWinner takes the score, for which it expects to be previously validated
// pl1 and pl2 ids and returns the id of the winner
func determineMatchWinner(score, pl1Id, pl2Id string) string {
	sets := strings.Split(score, ",")
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
		return pl1Id
	} else {
		return pl2Id
	}
}

type MatchStats struct {
	WonMatches int
	Pts        int
	SetsWon    int
	SetsLost   int
	GamesWon   int
	GamesLost  int
}

func calcMatchStats(score string, isPl1 bool) MatchStats {
	stats := MatchStats{}
	sets := strings.Split(score, ",")

	for i, set := range sets {
		setSl := strings.Split(set, "-")
		pl1Score, _ := strconv.Atoi(setSl[0])
		pl2Score, _ := strconv.Atoi(setSl[1])

		// check if third "set" is super tie-break
		if i == 2 && (pl1Score >= 10 || pl2Score >= 10) {
			if isPl1 {
				if pl1Score > pl2Score {
					stats.SetsWon++
				} else {
					stats.SetsLost++
				}
			} else {
				if pl2Score > pl1Score {
					stats.SetsWon++
				} else {
					stats.SetsWon++
				}
			}

			// skip counting games for super tie-break
			continue
		}

		if isPl1 {
			stats.GamesWon += pl1Score
			stats.GamesLost += pl2Score
			if pl1Score > pl2Score {
				stats.SetsWon++
			} else {
				stats.SetsLost++
			}
		} else {
			stats.GamesWon += pl2Score
			stats.GamesLost += pl1Score
			if pl2Score > pl1Score {
				stats.SetsWon++
			} else {
				stats.SetsLost++
			}
		}
	}

	if stats.SetsWon == 2 {
		stats.WonMatches = 1
		stats.Pts = 2
	} else {
		stats.WonMatches = 0
		stats.Pts = 1
	}

	return stats
}
