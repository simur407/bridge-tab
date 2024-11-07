package cmd

import (
	"github.com/spf13/cobra"
)

var removeTurnamentCmd = &cobra.Command{
	Use: "remove",
	Short: "Bridge Tab CLI to manage turnaments",
	Long: `Bridge Tab CLI is a tool to manage duplicate bridge turnaments. 
It allows organisers or umpires to manage turnaments before and durring the tournament.`,

	Run: func(cmd *cobra.Command, args []string) {
	},
}
