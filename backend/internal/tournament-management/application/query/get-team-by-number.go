package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type GetTeamByNumberQuery struct {
	TournamentId string
	Number       int
}

func (q *GetTeamByNumberQuery) Execute(repo domain.TeamReadRepository) (*domain.TeamDto, error) {
	t, err := repo.FindByNumber(&q.TournamentId, q.Number)
	if err != nil {
		return nil, err
	}

	return t, nil
}
