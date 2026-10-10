package score_registration

import (
	"errors"

	domain "bridge-tab/internal/score-registration/domain"
)

var ErrNoPlayedResults = errors.New("no played results to register")

type RegisterPlayedResultsCommand struct {
	TournamentId string
	Results      []domain.PlayedResult
}

func (c *RegisterPlayedResultsCommand) Execute(repo domain.PlayedResultRepository) error {
	if len(c.Results) == 0 {
		return ErrNoPlayedResults
	}
	for _, result := range c.Results {
		if result.Contract == "" {
			return domain.ErrRoundResultIncomplete
		}
	}
	return repo.Replace(domain.TournamentScoreId(c.TournamentId), c.Results)
}
