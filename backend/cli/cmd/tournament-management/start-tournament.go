package tournament_management

import (
	rounds_domain "bridge-tab/internal/rounds-registration/domain"
	tournament_management "bridge-tab/internal/tournament-management/application/command"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
	"fmt"

	"github.com/spf13/cobra"
)

var startTournamentId string

var startTournamentCmd = func(
	TournamentRepository *tournament_domain.TournamentRepository,
	TeamRepository *tournament_domain.TeamReadRepository,
	BoardProtocolRepository *tournament_domain.BoardProtocolReadRepository,
	GameSessionRepository *rounds_domain.GameSessionRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:   "start",
		Short: "Starts a tournament",
		Long: `This command starts a tournament. 
	Once tournament is started there will be no more modifications to the tournament setup available.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			command := &tournament_management.StartTurnamentCommand{TournamentId: startTournamentId}

			if err := command.Execute(*TournamentRepository, *TeamRepository, *BoardProtocolRepository, *GameSessionRepository); err != nil {
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
