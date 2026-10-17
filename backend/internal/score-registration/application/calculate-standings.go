package score_registration

import (
	"errors"
	"fmt"

	domain "bridge-tab/internal/score-registration/domain"
)

var ErrPlayedResultsNotFound = errors.New("played results not found for tournament")

type CalculateStandingsCommand struct {
	TournamentId string
	Method       domain.ScoringMethodName
}

func (c *CalculateStandingsCommand) Execute(
	playedResultRepo domain.PlayedResultRepository,
	scoreRepo domain.TournamentScoreRepository,
) error {
	method := c.Method
	if method == "" {
		return fmt.Errorf("scoring method is required")
	}

	tournamentId := domain.TournamentScoreId(c.TournamentId)
	results, err := playedResultRepo.Find(tournamentId)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		return ErrPlayedResultsNotFound
	}

	score, err := domain.CalculateTournamentScore(
		tournamentId,
		domain.NewScoringMethod(method),
		results,
	)
	if err != nil {
		return err
	}
	return scoreRepo.Save(score)
}
