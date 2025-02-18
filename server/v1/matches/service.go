package matches

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
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

func (s *service) processCreateMatch(ctx context.Context, model CreateMatchRequestModel) (*MatchModel, error) {
	// validate params in the db layer
	courtExists, seasonExists, leagueExists, playerOneExists, playerTwoExists, playersInLeague, err := s.store.validateInsertUpdateMatch(ctx, model.CourtId, model.SeasonId, model.LeagueId, model.PlayerOneId, model.PlayerTwoId)
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

	// begin tx
	tx, err := s.store.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning tx: %v", err)
	}

	defer func() {
		if err != nil && err != pgx.ErrTxClosed {
			log.Printf("failed to rollback the insert match tx: %v", err)
		}
	}()

	// insert match
	match, err := s.store.insertMatch(ctx, tx, model)
	if err != nil {
		return nil, err
	}

	// for cases where the score is submitted upon match creation
	if model.Score != nil {
		winnerId := determineMatchWinner(*model.Score, model.PlayerOneId, model.PlayerTwoId)
		model.WinnerId = &winnerId

		// calc player statistics
		pl1Stats := calcMatchStats(*model.Score, true)
		pl2Stats := calcMatchStats(*model.Score, true)

		// update player stats
		err = s.store.updatePlayerStatistics(ctx, tx, *model.WinnerId, model.PlayerOneId, model.PlayerTwoId)
		if err != nil {
			return nil, err
		}

		// update standings for pl1
		err = s.store.updateStanding(ctx, tx, model.SeasonId, model.LeagueId, model.PlayerOneId, pl1Stats)
		if err != nil {
			return nil, err
		}

		// update standings for pl2
		err = s.store.updateStanding(ctx, tx, model.SeasonId, model.LeagueId, model.PlayerTwoId, pl2Stats)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commiting tx: %v", err)
	}

	return &match, nil
}

func (s *service) processGetMatches(ctx context.Context, seasonId, leagueId string) ([]MatchModel, error) {
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

func (s *service) processUpdateMatch(ctx context.Context, model UpdateMatchRequestModel) (*MatchModel, error) {
	// check if able to modify match
	hasScore, err := s.store.checkMatchScore(ctx, model.MatchId)
	if err != nil {
		return nil, err
	}

	if hasScore {
		return nil, fmt.Errorf("%w: not able to modify a match that has a score", response.ErrConflict)
	}

	// validate params in the db layer
	courtExists, seasonExists, leagueExists, playerOneExists, playerTwoExists, playersInLeague, err := s.store.validateInsertUpdateMatch(ctx, model.CourtId, model.SeasonId, model.LeagueId, model.PlayerOneId, model.PlayerTwoId)
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

	mm, err := s.store.updateMatch(ctx, nil, model)
	if err != nil {
		return nil, err
	}

	return mm, nil
}

func (s *service) processSubmitMatchScore(ctx context.Context, model SubmitMatchScoreRequestModel) (*MatchModel, error) {
	seasonExists, leagueExists, matchExists, err := s.store.validateSubmitMatchScore(ctx, model.SeasonId, model.LeagueId, model.MatchId)
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
	match, err := s.store.findMatch(ctx, model.SeasonId, model.LeagueId, model.MatchId)
	if err != nil {
		return nil, err
	}

	// check if able to submit result
	if match.Score != nil {
		return nil, fmt.Errorf("%w: not able to submit match result, score already exists", response.ErrConflict)
	}

	// also add the match p1 and p2 info to the model struct
	model.PlayerOneId = match.PlayerOne.Id
	model.PlayerTwoId = match.PlayerTwo.Id
	model.WinnerId = determineMatchWinner(model.Score, match.PlayerOne.Id, match.PlayerTwo.Id)

	// begin tx
	tx, err := s.store.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err != pgx.ErrTxClosed {
			log.Printf("failed to rollback the insert match tx: %v", err)
		}
	}()

	result, err := s.store.updateMatchScore(ctx, tx, model.SeasonId, model.LeagueId, model.MatchId, model.Score, model.WinnerId)
	if err != nil {
		return nil, err
	}

	err = s.store.updatePlayerStatistics(ctx, tx, model.WinnerId, model.PlayerOneId, model.PlayerTwoId)
	if err != nil {
		return nil, err
	}

	// calc pl1 & pl2 match stats
	pl1MatchStats := calcMatchStats(model.Score, true)
	pl2MatchStats := calcMatchStats(model.Score, false)

	// update standing for pl1
	err = s.store.updateStanding(ctx, tx, model.SeasonId, model.LeagueId, model.PlayerOneId, pl1MatchStats)
	if err != nil {
		return nil, err
	}

	// update standing for pl2
	err = s.store.updateStanding(ctx, tx, model.SeasonId, model.LeagueId, model.PlayerTwoId, pl2MatchStats)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
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
