package cmd

import (
	tournament_management "bridge-tab/internal/tournament-management/application"
	"fmt"

	"github.com/spf13/cobra"
)

var removeTeamId string

var removeTeamCmd = &cobra.Command{
	Use: "remove",
	Short: "Removes team from the tournament",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		command := &tournament_management.RemoveTeamCommand{ TournamentId: teamsTournamentId, TeamId: removeTeamId }

		if err := command.Execute(TournamentRepository); err != nil {
			return err
		}

		fmt.Println("Team { Id:", command.TeamId, " } removed from Tournament { Id:", command.TournamentId, "}")
		return nil
	},
}

func init() {
	removeTeamCmd.Flags().StringVarP(&removeTeamId, "id", "i", "", "team id")
	removeTeamCmd.MarkFlagRequired("id")
}
