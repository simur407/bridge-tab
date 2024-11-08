package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type ListTournamentsQuery struct {}

func (c *ListTournamentsQuery) Execute(repo domain.TournamentReadRepository) ([]domain.TournamentDto, error) {
	t, err := repo.FindAll()
	if err != nil {
		return nil, err
	}

	return t, nil
}
