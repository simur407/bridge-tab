package rounds_registration

import (
	domain "bridge-tab/internal/rounds-registration/domain"
	tournament_management "bridge-tab/internal/tournament-management/application/query"
	tournament_management_domain "bridge-tab/internal/tournament-management/domain"
)

type GetRoundByTeamsQuery struct {
	GameSessionId    string
	FirstTeamNumber  int
	SecondTeamNumber int
	DealNo           int
}

func (q *GetRoundByTeamsQuery) Execute(repository domain.GameSessionReadRepository, teamRepository tournament_management_domain.TeamReadRepository) (*domain.RoundDto, error) {
	getFirstTeam := tournament_management.GetTeamByNumberQuery{TournamentId: q.GameSessionId, Number: q.FirstTeamNumber}
	firstTeam, err := getFirstTeam.Execute(teamRepository)
	if err != nil {
		return nil, err
	}

	getSecondTeam := tournament_management.GetTeamByNumberQuery{TournamentId: q.GameSessionId, Number: q.SecondTeamNumber}
	secondTeam, err := getSecondTeam.Execute(teamRepository)
	if err != nil {
		return nil, err
	}

	round, err := repository.FindRound(&q.GameSessionId, q.DealNo, firstTeam.Id, secondTeam.Id)
	if err != nil {
		return nil, err
	}

	return round, nil
}
