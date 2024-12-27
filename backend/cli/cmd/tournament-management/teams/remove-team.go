package teams

import (
	tournament_management "bridge-tab/internal/tournament-management/application"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
	"fmt"

	"github.com/spf13/cobra"
)

var removeTeamId string

var removeTeamCmd = func(TournamentRepository *tournament_domain.TournamentRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "remove",
		Short:        "Removes team from the tournament",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			command := &tournament_management.RemoveTeamCommand{TournamentId: teamsTournamentId, TeamId: removeTeamId}

			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			fmt.Println("Team { Id:", command.TeamId, " } removed from Tournament { Id:", command.TournamentId, "}")
			return nil
		},
	}

	command.Flags().StringVarP(&removeTeamId, "id", "i", "", "team id")
	command.MarkFlagRequired("id")

	return command
}
