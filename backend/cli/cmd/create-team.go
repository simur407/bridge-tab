package cmd

import (
	tournament_management "bridge-tab/internal/tournament-management/application"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var teamName string

var createTeamCmd = &cobra.Command{
	Use: "create",
	Short: "Creates team in a tournament",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		command := &tournament_management.CreateTeamCommand{ TournamentId: teamsTournamentId, TeamId: uuid.New().String(), Name: teamName }

		if err := command.Execute(TournamentRepository); err != nil {
			return err
		}

		fmt.Println("Team { Id:", command.TeamId, ", Name:", command.Name, " } created in Tournament { Id:", command.TournamentId, "}")
		return nil
	},
}

func init () {
	createTeamCmd.Flags().StringVarP(&teamName, "name", "n", "", "unique Team name")
	createTeamCmd.MarkFlagRequired("name")
}
