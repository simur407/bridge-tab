package cmd

import (
	"github.com/spf13/cobra"
)

var teamsTournamentId string

var teamsCmd = &cobra.Command{
	Use: "team",
	Short: "Resposible for managing Teams in Tournaments",
	Long: "Team Management allows organizers or umpires to manage Teams in Tournaments like: create, delete, join, leave, etc.",
}

func init() {
	teamsCmd.PersistentFlags().StringVarP(&teamsTournamentId, "tournamentId", "t", "", "tournament id")
	teamsCmd.MarkPersistentFlagRequired("tournamentId")
	teamsCmd.AddCommand(
		createTeamCmd,
		removeTeamCmd,
		joinTeamCmd,
		leaveTeamCmd,
		listTeamsCmd,
	)
}
