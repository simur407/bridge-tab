package teams

import (
	tournament_management "bridge-tab/internal/tournament-management/application/command"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var teamName string

var createTeamCmd = func(TournamentRepository *tournament_domain.TournamentRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "create",
		Short:        "Creates team in a tournament",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			command := &tournament_management.CreateTeamCommand{TournamentId: teamsTournamentId, TeamId: uuid.New().String(), Name: teamName}

			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			fmt.Println("Team { Id:", command.TeamId, ", Name:", command.Name, " } created in Tournament { Id:", command.TournamentId, "}")
			return nil
		},
	}

	command.Flags().StringVarP(&teamName, "name", "n", "", "unique Team name")
	command.MarkFlagRequired("name")

	return command
}
