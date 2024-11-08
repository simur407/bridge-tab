package cmd

import (
	tournament_management "bridge-tab/internal/tournament-management/application"
	"fmt"

	"github.com/spf13/cobra"
)

var leaveTournamentId string
var leaveContestantId string

var leaveTournamentCmd = &cobra.Command{
	Use: "leave",
	Short: "Leaves contestant from a tournament",
	Long: `This command allows to leave contestant from a tournament.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		command := &tournament_management.LeaveTournamentCommand{ TournamentId: leaveTournamentId, ContestantId: leaveContestantId }

		if err := command.Execute(TournamentRepository); err != nil {
			return err
		}

		fmt.Println("Contestant { Id:", command.ContestantId, " } left Tournament { Id:", command.TournamentId, "}")
		return nil
	},
}

func init() {
	leaveTournamentCmd.Flags().StringVarP(&leaveTournamentId, "id", "i", "", "tournament id")
	leaveTournamentCmd.MarkFlagRequired("id")
	leaveTournamentCmd.Flags().StringVarP(&leaveContestantId, "contestantId", "c", "", "contestant id")
	leaveTournamentCmd.MarkFlagRequired("contestantId")
}
