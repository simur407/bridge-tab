package play

import (
	"errors"

	domain "bridge-tab/internal/play/domain"
)

type GetMatchByTeamsQuery struct {
	TournamentId     string
	FirstTeamNumber  int
	SecondTeamNumber int
	DealNo           int
}

func (query *GetMatchByTeamsQuery) Execute(
	playSessionRepository domain.PlaySessionRepository,
	boardResultReadRepository domain.BoardResultReadRepository,
) (*ScheduledMatchDto, error) {
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
	match, err := session.Match(query.DealNo, firstTeam.Id, secondTeam.Id)
	if err != nil {
		return nil, err
	}
	dto := scheduledMatch(session, match)
	result, err := boardResultReadRepository.FindByTeams(query.TournamentId, match.DealNo, string(match.NsTeamId), string(match.EwTeamId))
	if errors.Is(err, domain.ErrBoardResultNotFound) {
		return &dto, nil
	}
	if err != nil {
		return nil, err
	}
	applyResult(&dto, *result)
	return &dto, nil
}
