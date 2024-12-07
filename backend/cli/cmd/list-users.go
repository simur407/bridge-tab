package cmd

import (
	"fmt"
	"strings"

	user "bridge-tab/internal/user/application"

	"github.com/spf13/cobra"
)

var listUsersCmd = &cobra.Command{
	Use:          "list",
	Short:        "Lists existing users in the system",
	Long:         `This command lists all existing users in the system.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		query := user.GetUsersCommand{}

		results, err := query.Execute(UserReadRepository)
		if err != nil {
			return err
		}

		fmt.Printf("%-36v | %-15v\n", "Id", "Name")
		fmt.Println(strings.Repeat("-", 55))
		for _, User := range results {
			fmt.Printf("%-36v | %-15v\n", User.Id, User.Name)
			fmt.Println(strings.Repeat("-", 55))
		}
		return nil
	},
}
