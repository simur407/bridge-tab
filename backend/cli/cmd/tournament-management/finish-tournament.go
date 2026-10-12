package tournament_management

import (
	play_domain "bridge-tab/internal/play/domain"
	score_domain "bridge-tab/internal/score-registration/domain"
	tournament_management "bridge-tab/internal/tournament-management/application/command"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
	"fmt"

	"github.com/spf13/cobra"
)

var finishTournamentId string

var finishTournamentCmd = func(
	TournamentRepository *tournament_domain.TournamentRepository,
	PlaySessionRepository *play_domain.PlaySessionRepository,
	BoardResultReadRepository *play_domain.BoardResultReadRepository,
	PlayedResultRepository *score_domain.PlayedResultRepository,
	TournamentScoreRepository *score_domain.TournamentScoreRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:   "finish",
		Short: "Finishes a tournament and calculates matchpoint standings",
		Long: `This command finishes a tournament after every board result has been recorded.
	Played results are already registered one by one. Finishing calculates matchpoint standings.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			command := &tournament_management.FinishTournamentCommand{TournamentId: finishTournamentId}

			if err := command.Execute(
				*TournamentRepository,
				*PlaySessionRepository,
				*BoardResultReadRepository,
				*PlayedResultRepository,
				*TournamentScoreRepository,
			); err != nil {
				return err
			}

			fmt.Println("Finished Tournament { Id:", command.TournamentId, "} and calculated matchpoint standings")
			return nil
		},
	}

	command.Flags().StringVarP(&finishTournamentId, "id", "i", "", "tournament id")
	command.MarkFlagRequired("id")

	return command
}
