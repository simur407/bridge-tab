package rounds_registration

import (
	domain "bridge-tab/internal/rounds-registration/domain"
	tournament_management "bridge-tab/internal/tournament-management/application/query"
	tournament_management_domain "bridge-tab/internal/tournament-management/domain"
	"errors"
	"regexp"
)

type PlayRoundCommand struct {
	GameSessionId  string
	PlayerId       string
	VersusTeamName string
	DealNo         int
	Contract       string
	Tricks         int
	Declarer       string
	OpeningLead    string
}

func (c *PlayRoundCommand) Execute(repository domain.GameSessionRepository, teamRepository tournament_management_domain.TeamReadRepository) error {
	if err := validate(c); err != nil {
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

	// find other team by name
	getTeamByName := tournament_management.GetTeamByNameQuery{TournamentId: c.GameSessionId, Name: c.VersusTeamName}
	versusTeam, err := getTeamByName.Execute(teamRepository)

	if err != nil {
		return err
	}

	err = t.AddRoundScore(c.DealNo, domain.TeamId(playerTeam.Id), domain.TeamId(versusTeam.Id), c.Contract, c.Tricks, c.Declarer, c.OpeningLead)

	if err != nil {
		return err
	}

	return repository.Save(t)
}

func validate(c *PlayRoundCommand) error {
	if c.GameSessionId == "" {
		return errors.New("game session id is empty")
	}
	if c.PlayerId == "" {
		return errors.New("player id is empty")
	}
	if c.VersusTeamName == "" {
		return errors.New("versus team name is empty")
	}
	if c.DealNo == 0 {
		return errors.New("deal no is empty")
	}
	if c.Contract == "" {
		return errors.New("contract is empty")
	}
	if c.Contract != "Pass" && c.Tricks == 0 {
		return errors.New("tricks is empty")
	}
	if c.Contract != "Pass" && c.Declarer == "" {
		return errors.New("declarer is empty")
	}
	if c.Contract != "Pass" && c.OpeningLead == "" {
		return errors.New("opening lead is empty")
	}

	match, err := regexp.MatchString("[1-7][CDHSN]x{0,2}|Pass", c.Contract)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("invalid contract")
	}

	if c.Contract != "Pass" && (c.Tricks < 0 || c.Tricks > 13) {
		return errors.New("invalid tricks")
	}

	if c.Contract != "Pass" && c.Declarer != "N" && c.Declarer != "E" && c.Declarer != "S" && c.Declarer != "W" {
		return errors.New("invalid declarer")
	}

	if c.Contract != "Pass" {
		match, err = regexp.MatchString("[2-9]|10|[AKQJ][CDHSN]", c.OpeningLead)
		if err != nil {
			return err
		}
		if !match {
			return errors.New("invalid opening lead")
		}
	}

	return nil
}
