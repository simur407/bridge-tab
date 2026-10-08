package score_registration

import (
	score_domain "bridge-tab/internal/score-registration/domain"

	"github.com/spf13/cobra"
)

var ScoreRegistrationCmd = func(
	playedResultRepository *score_domain.PlayedResultRepository,
	tournamentScoreRepository *score_domain.TournamentScoreRepository,
	tournamentScoreReadRepository *score_domain.TournamentScoreReadRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:   "score",
		Short: "Calculate and browse tournament scores and standings",
		Long:  "Score Registration calculates standings from transferred played results and lets organizers browse them.",
	}

	command.AddCommand(
		calculateCmd(playedResultRepository, tournamentScoreRepository),
		standingsCmd(tournamentScoreReadRepository),
		boardScoresCmd(tournamentScoreReadRepository),
	)

	return command
}
