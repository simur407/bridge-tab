package play

import (
	"errors"

	domain "bridge-tab/internal/play/domain"
	score_domain "bridge-tab/internal/score-registration/domain"
)

type EditBoardResultCommand struct {
	PlaySessionId    string
	FirstTeamNumber  int
	SecondTeamNumber int
	DealNo           int
	Contract         string
	Tricks           int
	Declarer         string
	OpeningLead      string
}

func (c *EditBoardResultCommand) Validate() error {
	if c.FirstTeamNumber < 1 {
		return errors.New("1st team number is empty")
	}
	if c.SecondTeamNumber < 1 {
		return errors.New("2nd team number is empty")
	}
	if err := validateBoardResultInput(c.PlaySessionId, c.DealNo, c.Contract, c.Tricks, c.Declarer, c.OpeningLead); err != nil {
		return err
	}
	c.Contract = normalizeContract(c.Contract)
	return nil
}

func (c *EditBoardResultCommand) Execute(
	boardResultRepository domain.BoardResultRepository,
	playSessionRepository domain.PlaySessionRepository,
	playedResultRepository score_domain.PlayedResultRepository,
) error {
	if err := c.Validate(); err != nil {
		return err
	}

	playSessionId := domain.PlaySessionId(c.PlaySessionId)
	session, err := playSessionRepository.Load(playSessionId)
	if err != nil {
		return err
	}
	firstTeam, err := session.TeamByNumber(c.FirstTeamNumber)
	if err != nil {
		return err
	}
	secondTeam, err := session.TeamByNumber(c.SecondTeamNumber)
	if err != nil {
		return err
	}

	result, err := boardResultRepository.Load(session.State.PlaySessionId, c.DealNo, firstTeam.Id, secondTeam.Id)
	if err != nil {
		return err
	}
	match, err := session.Match(result.State.DealNo, result.State.NsTeamId, result.State.EwTeamId)
	if err != nil {
		return err
	}

	result.Edit(c.Contract, c.Tricks, c.Declarer, c.OpeningLead)
	if err := boardResultRepository.Save(result); err != nil {
		return err
	}
	return registerPlayedResult(string(session.State.PlaySessionId), result.State, match.Vulnerable, playedResultRepository)
}
