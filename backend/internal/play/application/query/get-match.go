package play

import (
	domain "bridge-tab/internal/play/domain"
)

type GetMatchQuery struct {
	PlaySessionId    string
	PlayerId         string
	VersusTeamNumber int
	DealNo           int
}

func (query *GetMatchQuery) Execute(playSessionRepository domain.PlaySessionRepository) (*ScheduledMatchDto, error) {
	playSessionId := domain.PlaySessionId(query.PlaySessionId)
	playerId := domain.PlayerId(query.PlayerId)

	session, err := playSessionRepository.Load(playSessionId)
	if err != nil {
		return nil, err
	}
	playerTeam, err := session.TeamByMember(playerId)
	if err != nil {
		return nil, err
	}
	versusTeam, err := session.TeamByNumber(query.VersusTeamNumber)
	if err != nil {
		return nil, err
	}
	match, err := session.Match(query.DealNo, playerTeam.Id, versusTeam.Id)
	if err != nil {
		return nil, err
	}
	dto := scheduledMatch(session, match)
	return &dto, nil
}
