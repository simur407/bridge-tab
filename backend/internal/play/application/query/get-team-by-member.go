package play

import (
	domain "bridge-tab/internal/play/domain"
)

type GetTeamNumberByMemberQuery struct {
	PlaySessionId string
	PlayerId      string
}

func (query *GetTeamNumberByMemberQuery) Execute(playSessionRepository domain.PlaySessionRepository) (int, error) {
	playSessionId := domain.PlaySessionId(query.PlaySessionId)
	playerId := domain.PlayerId(query.PlayerId)

	session, err := playSessionRepository.Load(playSessionId)
	if err != nil {
		return 0, err
	}
	team, err := session.TeamByMember(playerId)
	if err != nil {
		return 0, err
	}
	return team.Number, nil
}
