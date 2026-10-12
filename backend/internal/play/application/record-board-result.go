package play

import (
	"errors"

	domain "bridge-tab/internal/play/domain"
	score_domain "bridge-tab/internal/score-registration/domain"
)

type RecordBoardResultCommand struct {
	TournamentId     string
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
	if err := validateBoardResultInput(c.TournamentId, c.DealNo, c.Contract, c.Tricks, c.Declarer, c.OpeningLead); err != nil {
		return err
	}
	c.Contract = normalizeContract(c.Contract)

	session, err := playSessionRepository.Load(domain.TournamentId(c.TournamentId))
	if err != nil {
		return err
	}
	playerTeam, err := session.TeamByMember(domain.PlayerId(c.PlayerId))
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

	existing, err := boardResultRepository.Load(session.State.TournamentId, c.DealNo, match.NsTeamId, match.EwTeamId)
	if err != nil && !errors.Is(err, domain.ErrBoardResultNotFound) {
		return err
	}
	if existing != nil {
		return domain.ErrBoardResultAlreadyRecorded
	}

	result := domain.RecordBoardResult(session.State.TournamentId, c.DealNo, match.NsTeamId, match.EwTeamId, c.Contract, c.Tricks, c.Declarer, c.OpeningLead)
	if err := boardResultRepository.Save(result); err != nil {
		return err
	}
	return registerPlayedResult(c.TournamentId, result.State, match.Vulnerable, playedResultRepository)
}
