package play

import (
	domain "bridge-tab/internal/play/domain"
)

type GetBoardResultByTeamsQuery struct {
	TournamentId     string
	FirstTeamNumber  int
	SecondTeamNumber int
	DealNo           int
}

func (query *GetBoardResultByTeamsQuery) Execute(
	boardResultReadRepository domain.BoardResultReadRepository,
	playSessionRepository domain.PlaySessionRepository,
) (*domain.BoardResultDto, error) {
	session, err := playSessionRepository.Load(domain.TournamentId(query.TournamentId))
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
	return boardResultReadRepository.FindByTeams(query.TournamentId, query.DealNo, string(firstTeam.Id), string(secondTeam.Id))
}
