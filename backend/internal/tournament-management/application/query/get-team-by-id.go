package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type GetTeamByIdQuery struct {
	TournamentId string
	TeamId       string
}

func (q *GetTeamByIdQuery) Execute(repo domain.TeamReadRepository) (*domain.TeamDto, error) {
	t, err := repo.FindById(&q.TournamentId, &q.TeamId)
	if err != nil {
		return nil, err
	}

	return t, nil
}
