package rounds_registration

import (
	"errors"

	domain "bridge-tab/internal/rounds-registration/domain"
	tournament_management "bridge-tab/internal/tournament-management/application/query"
	tournament_management_domain "bridge-tab/internal/tournament-management/domain"
)

type RecordRoundCommand struct {
	GameSessionId    string
	FirstTeamNumber  int
	SecondTeamNumber int
	DealNo           int
	Contract         string
	Tricks           int
	Declarer         string
	OpeningLead      string
}

func (c *RecordRoundCommand) Validate() error {
	if c.GameSessionId == "" {
		return errors.New("game session id is empty")
	}
	if c.FirstTeamNumber < 1 {
		return errors.New("1st team number is empty")
	}
	if c.SecondTeamNumber < 1 {
		return errors.New("2nd team number is empty")
	}
	if c.DealNo == 0 {
		return errors.New("deal no is empty")
	}
	if err := validateRoundScore(c.Contract, c.Tricks, c.Declarer, c.OpeningLead); err != nil {
		return err
	}

	c.Contract = normalizeContract(c.Contract)
	return nil
}

func (c *RecordRoundCommand) Execute(repository domain.GameSessionRepository, teamRepository tournament_management_domain.TeamReadRepository) error {
	if err := c.Validate(); err != nil {
		return err
	}

	gameSessionId := domain.GameSessionId(c.GameSessionId)

	t, err := repository.Load(&gameSessionId)
	if err != nil {
		return err
	}

	getFirstTeam := tournament_management.GetTeamByNumberQuery{TournamentId: c.GameSessionId, Number: c.FirstTeamNumber}
	firstTeam, err := getFirstTeam.Execute(teamRepository)
	if err != nil {
		return err
	}

	getSecondTeam := tournament_management.GetTeamByNumberQuery{TournamentId: c.GameSessionId, Number: c.SecondTeamNumber}
	secondTeam, err := getSecondTeam.Execute(teamRepository)
	if err != nil {
		return err
	}

	err = t.AddRoundScore(c.DealNo, domain.TeamId(firstTeam.Id), domain.TeamId(secondTeam.Id), c.Contract, c.Tricks, c.Declarer, c.OpeningLead)
	if err != nil {
		return err
	}

	return repository.Save(t)
}
