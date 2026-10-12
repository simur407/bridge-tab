package play

import (
	domain "bridge-tab/internal/play/domain"
)

type GetTeamNumberByMemberQuery struct {
	TournamentId string
	PlayerId     string
}

func (query *GetTeamNumberByMemberQuery) Execute(playSessionRepository domain.PlaySessionRepository) (int, error) {
	session, err := playSessionRepository.Load(domain.TournamentId(query.TournamentId))
	if err != nil {
		return 0, err
	}
	team, err := session.TeamByMember(domain.PlayerId(query.PlayerId))
	if err != nil {
		return 0, err
	}
	return team.Number, nil
}
