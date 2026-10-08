package tournament_management

import (
	"errors"
	"fmt"

	rounds_registration_domain "bridge-tab/internal/rounds-registration/domain"
	score_registration "bridge-tab/internal/score-registration/application"
	score_registration_domain "bridge-tab/internal/score-registration/domain"
	query "bridge-tab/internal/tournament-management/application/query"
	domain "bridge-tab/internal/tournament-management/domain"
)

var ErrRoundsNotAllPlayed = errors.New("not all rounds have been played")

type FinishTournamentCommand struct {
	TournamentId string
}

func (c *FinishTournamentCommand) Execute(
	tournamentRepo domain.TournamentRepository,
	gameSessionRepo rounds_registration_domain.GameSessionRepository,
	boardProtocolReadRepo domain.BoardProtocolReadRepository,
	playedResultRepo score_registration_domain.PlayedResultRepository,
) error {
	id := domain.TournamentId(c.TournamentId)
	t, err := tournamentRepo.Load(&id)
	if err != nil {
		return err
	}

	gameSessionId := rounds_registration_domain.GameSessionId(c.TournamentId)
	session, err := gameSessionRepo.Load(&gameSessionId)
	if err != nil {
		return err
	}

	for _, round := range session.State.Rounds {
		if !round.IsPlayed() {
			return ErrRoundsNotAllPlayed
		}
	}

	listBoardProtocols := query.ListBoardProtocolsQuery{TournamentId: c.TournamentId}
	boardProtocols, err := listBoardProtocols.Execute(boardProtocolReadRepo)
	if err != nil {
		return err
	}

	vulByBoard := map[int]domain.Vulnerable{}
	for _, protocol := range boardProtocols {
		vulByBoard[protocol.BoardNo] = vulnerableFromString(protocol.Vulnerable)
	}

	results := make([]score_registration_domain.PlayedResult, 0, len(session.State.Rounds))
	for _, round := range session.State.Rounds {
		vul, ok := vulByBoard[round.DealNo]
		if !ok {
			return fmt.Errorf("missing board protocol for deal %d", round.DealNo)
		}
		nsVul, ewVul := vulnerabilityFlags(vul)
		results = append(results, score_registration_domain.PlayedResult{
			DealNo:       round.DealNo,
			NsTeamId:     string(*round.NsTeam),
			EwTeamId:     string(*round.EwTeam),
			Contract:     round.Contract,
			Tricks:       round.Tricks,
			Declarer:     round.Declarer,
			NsVulnerable: nsVul,
			EwVulnerable: ewVul,
		})
	}

	register := score_registration.RegisterPlayedResultsCommand{
		TournamentId: c.TournamentId,
		Results:      results,
	}
	if err := register.Execute(playedResultRepo); err != nil {
		return err
	}

	if err := t.Finish(); err != nil {
		return err
	}
	return tournamentRepo.Save(t)
}

func vulnerableFromString(value string) domain.Vulnerable {
	switch value {
	case "NS":
		return domain.NS
	case "EW":
		return domain.EW
	case "Both":
		return domain.Both
	default:
		return domain.None
	}
}

func vulnerabilityFlags(vul domain.Vulnerable) (nsVul, ewVul bool) {
	switch vul {
	case domain.NS:
		return true, false
	case domain.EW:
		return false, true
	case domain.Both:
		return true, true
	default:
		return false, false
	}
}
