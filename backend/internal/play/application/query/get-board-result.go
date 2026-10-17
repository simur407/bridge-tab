package play

import (
	domain "bridge-tab/internal/play/domain"
)

type GetBoardResultByTeamsQuery struct {
	PlaySessionId    string
	FirstTeamNumber  int
	SecondTeamNumber int
	DealNo           int
}

func (query *GetBoardResultByTeamsQuery) Execute(
	boardResultReadRepository domain.BoardResultReadRepository,
	playSessionRepository domain.PlaySessionRepository,
) (*domain.BoardResultDto, error) {
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
	firstTeamId := string(firstTeam.Id)
	secondTeamId := string(secondTeam.Id)
	return boardResultReadRepository.FindByTeams(string(playSessionId), query.DealNo, firstTeamId, secondTeamId)
}
