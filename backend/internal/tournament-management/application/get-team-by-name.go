package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type GetTeamByNameQuery struct {
	TournamentId string
	Name         string
}

func (q *GetTeamByNameQuery) Execute(repo domain.TeamReadRepository) (*domain.TeamDto, error) {
	t, err := repo.FindByName(&q.TournamentId, &q.Name)
	if err != nil {
		return nil, err
	}

	return t, nil
}
