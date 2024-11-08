package cmd

import (
	tournament_management "bridge-tab/internal/tournament-management/application"
	"fmt"

	"github.com/spf13/cobra"
)

var removeTournamentId string

var removeTournamentCmd = &cobra.Command{
	Use: "remove",
	Short: "Bridge Tab CLI to manage tournaments",
	Long: `Bridge Tab CLI is a tool to manage duplicate bridge tournaments. 
It allows organizers or umpires to manage tournaments before and during the Tournament.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		command := &tournament_management.RemoveTournamentCommand{ TournamentId: removeTournamentId }

		if err := command.Execute(TournamentRepository); err != nil {
			return err
		}

		fmt.Println("Removed Tournament { Id:", command.TournamentId, "}")
		return nil
	},
}

func init () {
	removeTournamentCmd.Flags().StringVarP(&removeTournamentId, "id", "i", "", "tournament id")
	removeTournamentCmd.MarkFlagRequired("id")
}
