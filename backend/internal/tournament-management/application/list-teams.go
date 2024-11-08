package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type ListTeamsQuery struct {
	TournamentId 	string
}

func (c *ListTeamsQuery) Execute(repo domain.TeamReadRepository) ([]domain.TeamDto, error) {
	t, err := repo.FindAll(&c.TournamentId)
	if err != nil {
		return nil, err
	}

	return t, nil
}
