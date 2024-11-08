package cmd

import (
	tournament_management "bridge-tab/internal/tournament-management/application"
	"fmt"

	"github.com/spf13/cobra"
)

var joinTeamId string
var joinTeamContestantId string

var joinTeamCmd = &cobra.Command{
	Use: "join",
	Short: "Joins contestant to a team in the tournament",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		command := &tournament_management.JoinTeamCommand{ 
			TournamentId: teamsTournamentId, 
			TeamId: joinTeamId, 
			ContestantId: joinTeamContestantId,
		}

		if err := command.Execute(TournamentRepository); err != nil {
			return err
		}

		fmt.Println("Contestant { Id:", command.ContestantId, " } joined Team { Id:", command.TeamId, " } in Tournament { Id:", command.TournamentId, "}")
		return nil
	},
}

func init() {
	joinTeamCmd.Flags().StringVarP(&joinTeamId, "id", "i", "", "team id")
	joinTeamCmd.MarkFlagRequired("id")
	joinTeamCmd.Flags().StringVarP(&joinTeamContestantId, "contestantId", "c", "", "contestant id")
	joinTeamCmd.MarkFlagRequired("contestantId")
}
