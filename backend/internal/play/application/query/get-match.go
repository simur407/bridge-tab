package play

import (
	domain "bridge-tab/internal/play/domain"
)

type GetMatchQuery struct {
	TournamentId     string
	PlayerId         string
	VersusTeamNumber int
	DealNo           int
}

func (query *GetMatchQuery) Execute(playSessionRepository domain.PlaySessionRepository) (*ScheduledMatchDto, error) {
	session, err := playSessionRepository.Load(domain.TournamentId(query.TournamentId))
	if err != nil {
		return nil, err
	}
	playerTeam, err := session.TeamByMember(domain.PlayerId(query.PlayerId))
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
