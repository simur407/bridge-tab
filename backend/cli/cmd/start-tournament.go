package cmd

import (
	tournament_management "bridge-tab/internal/tournament-management/application"
	"fmt"

	"github.com/spf13/cobra"
)

var startTournamentId string

var startTournamentCmd = &cobra.Command{
	Use: "start",
	Short: "Starts a tournament",
	Long: `This command starts a tournament. 
	Once tournament is started there will be no more modifications to the tournament setup available.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		command := &tournament_management.StartTurnamentCommand{ TournamentId: startTournamentId }

		if err := command.Execute(TournamentRepository); err != nil {
			return err
		}

		fmt.Println("Started Tournament { Id:", command.TournamentId, "}")
		return nil
	},
}

func init() {
	startTournamentCmd.Flags().StringVarP(&startTournamentId, "id", "i", "", "tournament id")
	startTournamentCmd.MarkFlagRequired("id")
}
