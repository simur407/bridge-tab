package cmd

import (
	"github.com/spf13/cobra"
)

var tournamentManagementCmd = &cobra.Command{
	Use: "tournament",
	Short: "Resposible for managing Tournaments",
	Long: "Tournament Management allows organisers or umpires to manage Tournaments like: create, delete, add deal protocols, etc.",
}

func init() {
	tournamentManagementCmd.AddCommand(
		createTournamentCmd,
		removeTournamentCmd,
		listTournamentsCmd,
		startTournamentCmd,
		joinTournamentCmd,
		leaveTournamentCmd,
		listContestantsCmd,
		teamsCmd,
	)
}
