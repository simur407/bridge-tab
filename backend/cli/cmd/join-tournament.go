package cmd

import (
	tournament_management "bridge-tab/internal/tournament-management/application"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var joinTournamentId string
var joinContestantId string

var joinTournamentCmd = &cobra.Command{
	Use: "join",
	Short: "Joins contestant to a tournament",
	Long: `This command joins contestant to a tournament.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if joinContestantId == "" {
			joinContestantId = uuid.New().String()
		}
		command := &tournament_management.JoinTournamentCommand{ TournamentId: joinTournamentId, ContestantId: joinContestantId }

		if err := command.Execute(TournamentRepository); err != nil {
			return err
		}

		fmt.Println("Contestant { Id:", command.ContestantId, " } joined Tournament { Id:", command.TournamentId, "}")
		return nil
	},
}

func init() {
	joinTournamentCmd.Flags().StringVarP(&joinTournamentId, "id", "i", "", "tournament id")
	joinTournamentCmd.MarkFlagRequired("id")
	joinTournamentCmd.Flags().StringVarP(&joinContestantId, "contestantId", "c", "", "contestant id")
}
