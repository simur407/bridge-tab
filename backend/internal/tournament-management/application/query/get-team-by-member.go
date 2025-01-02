package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type GetTeamByMemberQuery struct {
	TournamentId string
	MemberId     string
}

func (q *GetTeamByMemberQuery) Execute(repo domain.TeamReadRepository) (*domain.TeamDto, error) {
	t, err := repo.FindByMemberId(&q.TournamentId, &q.MemberId)
	if err != nil {
		return nil, err
	}

	return t, nil
}
