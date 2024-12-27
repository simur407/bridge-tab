package tournament_management

import (
	tournament_management "bridge-tab/internal/tournament-management/application"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
	"fmt"

	"github.com/spf13/cobra"
)

var startTournamentId string

var startTournamentCmd = func(TournamentRepository *tournament_domain.TournamentRepository) *cobra.Command {
	command := &cobra.Command{
		Use:   "start",
		Short: "Starts a tournament",
		Long: `This command starts a tournament. 
	Once tournament is started there will be no more modifications to the tournament setup available.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			command := &tournament_management.StartTurnamentCommand{TournamentId: startTournamentId}

			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			fmt.Println("Started Tournament { Id:", command.TournamentId, "}")
			return nil
		},
	}

	command.Flags().StringVarP(&startTournamentId, "id", "i", "", "tournament id")
	command.MarkFlagRequired("id")

	return command
}
