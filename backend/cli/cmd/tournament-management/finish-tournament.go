package tournament_management

import (
	rounds_domain "bridge-tab/internal/rounds-registration/domain"
	score_domain "bridge-tab/internal/score-registration/domain"
	tournament_management "bridge-tab/internal/tournament-management/application/command"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
	"fmt"

	"github.com/spf13/cobra"
)

var finishTournamentId string

var finishTournamentCmd = func(
	TournamentRepository *tournament_domain.TournamentRepository,
	GameSessionRepository *rounds_domain.GameSessionRepository,
	BoardProtocolRepository *tournament_domain.BoardProtocolReadRepository,
	PlayedResultRepository *score_domain.PlayedResultRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:   "finish",
		Short: "Finishes a tournament and transfers results for scoring",
		Long: `This command finishes a tournament after all rounds are played.
	It transfers played results into score registration. Use score calculate to produce standings.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			command := &tournament_management.FinishTournamentCommand{TournamentId: finishTournamentId}

			if err := command.Execute(
				*TournamentRepository,
				*GameSessionRepository,
				*BoardProtocolRepository,
				*PlayedResultRepository,
			); err != nil {
				return err
			}

			fmt.Println("Finished Tournament { Id:", command.TournamentId, "}")
			return nil
		},
	}

	command.Flags().StringVarP(&finishTournamentId, "id", "i", "", "tournament id")
	command.MarkFlagRequired("id")

	return command
}
