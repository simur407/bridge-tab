package score_registration

import (
	domain "bridge-tab/internal/score-registration/domain"
)

type RegisterPlayedResultCommand struct {
	TournamentId string
	Result       domain.PlayedResult
}

func (c *RegisterPlayedResultCommand) Execute(repo domain.PlayedResultRepository) error {
	if c.Result.Contract == "" {
		return domain.ErrRoundResultIncomplete
	}
	return repo.Save(domain.TournamentScoreId(c.TournamentId), c.Result)
}
