package rounds_registration

import (
	domain "bridge-tab/internal/rounds-registration/domain"
	tournament_management "bridge-tab/internal/tournament-management/application/query"
	tournament_management_domain "bridge-tab/internal/tournament-management/domain"
)

type GetRoundQuery struct {
	GameSessionId    string
	PlayerId         string
	VersusTeamNumber int
	DealNo           int
}

func (q *GetRoundQuery) Execute(repository domain.GameSessionReadRepository, teamRepository tournament_management_domain.TeamReadRepository) (*domain.RoundDto, error) {
	// find player team
	getPlayerTeam := tournament_management.GetTeamByMemberQuery{TournamentId: q.GameSessionId, MemberId: q.PlayerId}
	playerTeam, err := getPlayerTeam.Execute(teamRepository)

	if err != nil {
		return nil, err
	}

	getTeamByNumber := tournament_management.GetTeamByNumberQuery{TournamentId: q.GameSessionId, Number: q.VersusTeamNumber}
	versusTeam, err := getTeamByNumber.Execute(teamRepository)

	if err != nil {
		return nil, err
	}

	round, err := repository.FindRound(&q.GameSessionId, q.DealNo, playerTeam.Id, versusTeam.Id)

	if err != nil {
		return nil, err
	}

	return round, nil
}
