package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type ListContestantsQuery struct {
	TournamentId 	string
}

func (c *ListContestantsQuery) Execute(repo domain.TournamentReadRepository) ([]domain.ContestantDto, error) {
	id := domain.TournamentId(c.TournamentId)
	t, err := repo.FindAllContestants(&id)
	if err != nil {
		return nil, err
	}

	return t, nil
}
