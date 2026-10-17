package play

import (
	"errors"

	domain "bridge-tab/internal/play/domain"
	score_domain "bridge-tab/internal/score-registration/domain"
)

type RecordBoardResultCommand struct {
	PlaySessionId    string
	PlayerId         string
	VersusTeamNumber int
	DealNo           int
	Contract         string
	Tricks           int
	Declarer         string
	OpeningLead      string
}

func (c *RecordBoardResultCommand) Execute(
	boardResultRepository domain.BoardResultRepository,
	playSessionRepository domain.PlaySessionRepository,
	playedResultRepository score_domain.PlayedResultRepository,
) error {
	if c.PlayerId == "" {
		return errors.New("player id is empty")
	}
	if c.VersusTeamNumber < 1 {
		return errors.New("versus team number is empty")
	}
	if err := validateBoardResultInput(c.PlaySessionId, c.DealNo, c.Contract, c.Tricks, c.Declarer, c.OpeningLead); err != nil {
		return err
	}
	c.Contract = normalizeContract(c.Contract)

	playSessionId := domain.PlaySessionId(c.PlaySessionId)
	playerId := domain.PlayerId(c.PlayerId)

	session, err := playSessionRepository.Load(playSessionId)
	if err != nil {
		return err
	}
	playerTeam, err := session.TeamByMember(playerId)
	if err != nil {
		return err
	}
	versusTeam, err := session.TeamByNumber(c.VersusTeamNumber)
	if err != nil {
		return err
	}
	match, err := session.Match(c.DealNo, playerTeam.Id, versusTeam.Id)
	if err != nil {
		return err
	}

	existing, err := boardResultRepository.Load(session.State.PlaySessionId, c.DealNo, match.NsTeamId, match.EwTeamId)
	if err != nil && !errors.Is(err, domain.ErrBoardResultNotFound) {
		return err
	}
	if existing != nil {
		return domain.ErrBoardResultAlreadyRecorded
	}

	result := domain.RecordBoardResult(session.State.PlaySessionId, c.DealNo, match.NsTeamId, match.EwTeamId, c.Contract, c.Tricks, c.Declarer, c.OpeningLead)
	if err := boardResultRepository.Save(result); err != nil {
		return err
	}
	return registerPlayedResult(string(session.State.PlaySessionId), result.State, match.Vulnerable, playedResultRepository)
}
