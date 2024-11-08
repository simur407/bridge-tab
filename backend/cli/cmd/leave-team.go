package cmd

import (
	tournament_management "bridge-tab/internal/tournament-management/application"
	"fmt"

	"github.com/spf13/cobra"
)

var leaveTeamId string
var leaveTeamContestantId string

var leaveTeamCmd = &cobra.Command{
	Use: "leave",
	Short: "Leaves team from the tournament",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		command := &tournament_management.LeaveTeamCommand{ 
			TournamentId: teamsTournamentId, 
			TeamId: leaveTeamId, 
			ContestantId: leaveTeamContestantId,
		}

		if err := command.Execute(TournamentRepository); err != nil {
			return err
		}

		fmt.Println("Contestant { Id:", command.ContestantId, " } leave Team { Id:", command.TeamId, " } in Tournament { Id:", command.TournamentId, "}")
		return nil
	},
}

func init() {
	leaveTeamCmd.Flags().StringVarP(&leaveTeamId, "id", "i", "", "team id")
	leaveTeamCmd.MarkFlagRequired("id")
	leaveTeamCmd.Flags().StringVarP(&leaveTeamContestantId, "contestantId", "c", "", "contestant id")
	leaveTeamCmd.MarkFlagRequired("contestantId")
}
