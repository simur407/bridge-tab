package cmd

import (
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Responsible for managing Users",
	Long:  "User allows admins to manage Users like: list",
}

func init() {
	userCmd.AddCommand(
		listUsersCmd,
	)
}
