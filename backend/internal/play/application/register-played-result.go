package play

import (
	"fmt"

	domain "bridge-tab/internal/play/domain"
	score_app "bridge-tab/internal/score-registration/application"
	score_domain "bridge-tab/internal/score-registration/domain"
)

func registerPlayedResult(
	playSessionId string,
	state domain.BoardResultState,
	vulnerable string,
	playedResultRepository score_domain.PlayedResultRepository,
) error {
	nsTeamId := string(state.NsTeamId)
	ewTeamId := string(state.EwTeamId)
	tournamentId := playSessionId
	nsVulnerable, ewVulnerable := vulnerabilityFlags(vulnerable)
	command := score_app.RegisterPlayedResultCommand{
		TournamentId: tournamentId,
		Result: score_domain.PlayedResult{
			DealNo:       state.DealNo,
			NsTeamId:     nsTeamId,
			EwTeamId:     ewTeamId,
			Contract:     state.Contract,
			Tricks:       state.Tricks,
			Declarer:     state.Declarer,
			NsVulnerable: nsVulnerable,
			EwVulnerable: ewVulnerable,
		},
	}
	if err := command.Execute(playedResultRepository); err != nil {
		return fmt.Errorf("register played result: %w", err)
	}
	return nil
}

func vulnerabilityFlags(vulnerable string) (nsVulnerable bool, ewVulnerable bool) {
	switch vulnerable {
	case "NS":
		return true, false
	case "EW":
		return false, true
	case "Both":
		return true, true
	default:
		return false, false
	}
}
