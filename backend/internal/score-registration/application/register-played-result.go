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
	tournamentId := domain.TournamentScoreId(c.TournamentId)
	return repo.Save(tournamentId, c.Result)
}
