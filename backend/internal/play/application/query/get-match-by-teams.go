package play

import (
	"errors"

	domain "bridge-tab/internal/play/domain"
)

type GetMatchByTeamsQuery struct {
	PlaySessionId    string
	FirstTeamNumber  int
	SecondTeamNumber int
	DealNo           int
}

func (query *GetMatchByTeamsQuery) Execute(
	playSessionRepository domain.PlaySessionRepository,
	boardResultReadRepository domain.BoardResultReadRepository,
) (*ScheduledMatchDto, error) {
	playSessionId := domain.PlaySessionId(query.PlaySessionId)
	session, err := playSessionRepository.Load(playSessionId)
	if err != nil {
		return nil, err
	}
	firstTeam, err := session.TeamByNumber(query.FirstTeamNumber)
	if err != nil {
		return nil, err
	}
	secondTeam, err := session.TeamByNumber(query.SecondTeamNumber)
	if err != nil {
		return nil, err
	}
	match, err := session.Match(query.DealNo, firstTeam.Id, secondTeam.Id)
	if err != nil {
		return nil, err
	}
	dto := scheduledMatch(session, match)
	nsTeamId := string(match.NsTeamId)
	ewTeamId := string(match.EwTeamId)
	result, err := boardResultReadRepository.FindByTeams(string(playSessionId), match.DealNo, nsTeamId, ewTeamId)
	if errors.Is(err, domain.ErrBoardResultNotFound) {
		return &dto, nil
	}
	if err != nil {
		return nil, err
	}
	applyResult(&dto, *result)
	return &dto, nil
}
