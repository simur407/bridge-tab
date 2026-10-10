package rounds_registration

import (
	domain "bridge-tab/internal/rounds-registration/domain"
	tournament_management "bridge-tab/internal/tournament-management/application/query"
	tournament_management_domain "bridge-tab/internal/tournament-management/domain"
)

type PlayRoundCommand struct {
	GameSessionId    string
	PlayerId         string
	VersusTeamNumber int
	DealNo           int
	Contract         string
	Tricks           int
	Declarer         string
	OpeningLead      string
}

func (c *PlayRoundCommand) Execute(repository domain.GameSessionRepository, teamRepository tournament_management_domain.TeamReadRepository) error {
	if err := validateRoundInput(c.GameSessionId, c.PlayerId, c.VersusTeamNumber, c.DealNo, c.Contract, c.Tricks, c.Declarer, c.OpeningLead); err != nil {
		return err
	}

	gameSessionId := domain.GameSessionId(c.GameSessionId)

	t, err := repository.Load(&gameSessionId)
	if err != nil {
		return err
	}
	// find player team
	getPlayerTeam := tournament_management.GetTeamByMemberQuery{TournamentId: c.GameSessionId, MemberId: c.PlayerId}
	playerTeam, err := getPlayerTeam.Execute(teamRepository)

	if err != nil {
		return err
	}

	getTeamByNumber := tournament_management.GetTeamByNumberQuery{TournamentId: c.GameSessionId, Number: c.VersusTeamNumber}
	versusTeam, err := getTeamByNumber.Execute(teamRepository)

	if err != nil {
		return err
	}

	c.Contract = normalizeContract(c.Contract)

	err = t.AddRoundScore(c.DealNo, domain.TeamId(playerTeam.Id), domain.TeamId(versusTeam.Id), c.Contract, c.Tricks, c.Declarer, c.OpeningLead)

	if err != nil {
		return err
	}

	return repository.Save(t)
}
