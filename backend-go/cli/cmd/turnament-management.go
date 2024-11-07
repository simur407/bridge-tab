package cmd

import (
	turnament "bridge-tab/internal/turnament-management/domain"

	"github.com/spf13/cobra"
)

var turnamentManagementCmd = &cobra.Command{
	Use: "turnament",
	Short: "Resposible for managing turnaments",
	Long: "Turnament Management allows organisers or umpires to manage turnaments like: create, delete, add deal protocols, etc.",
}

var turnamentRepository turnament.TurnamentRepository = &turnament.InMemoryTurnamentRepository{}

func init() {
	turnamentManagementCmd.AddCommand(createTurnamentCmd)
	turnamentManagementCmd.AddCommand(removeTurnamentCmd)
	turnamentManagementCmd.AddCommand(listTurnamentsCmd)
}
