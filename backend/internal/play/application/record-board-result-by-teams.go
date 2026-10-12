package play

import (
	"errors"

	domain "bridge-tab/internal/play/domain"
	score_domain "bridge-tab/internal/score-registration/domain"
)

type RecordBoardResultByTeamsCommand struct {
	TournamentId     string
	FirstTeamNumber  int
	SecondTeamNumber int
	DealNo           int
	Contract         string
	Tricks           int
	Declarer         string
	OpeningLead      string
}

func (c *RecordBoardResultByTeamsCommand) Execute(
	boardResultRepository domain.BoardResultRepository,
	playSessionRepository domain.PlaySessionRepository,
	playedResultRepository score_domain.PlayedResultRepository,
) error {
	if c.FirstTeamNumber < 1 {
		return errors.New("1st team number is empty")
	}
	if c.SecondTeamNumber < 1 {
		return errors.New("2nd team number is empty")
	}
	if err := validateBoardResultInput(c.TournamentId, c.DealNo, c.Contract, c.Tricks, c.Declarer, c.OpeningLead); err != nil {
		return err
	}
	c.Contract = normalizeContract(c.Contract)

	session, err := playSessionRepository.Load(domain.TournamentId(c.TournamentId))
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
	match, err := session.Match(c.DealNo, firstTeam.Id, secondTeam.Id)
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
