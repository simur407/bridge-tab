package cmd

import (
	"fmt"
	"strings"

	tournament "bridge-tab/internal/tournament-management/application"

	"github.com/spf13/cobra"
)

var listTournamentsCmd = &cobra.Command{
	Use: "list",
	Short: "Lists existing tournaments",
	Long: `This command lists all existing tournaments.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		query := tournament.ListTournamentsQuery{}

		results, err := query.Execute(TournamentReadRepository)
		if err != nil {
			return err
		}

		fmt.Printf("%-36v | %-15v\n", "Id", "Name")
		fmt.Println(strings.Repeat("-", 55))
		for _, Tournament := range results {
			fmt.Printf("%-36v | %-15v\n", Tournament.Id, Tournament.Name)
			fmt.Println(strings.Repeat("-", 55))
		}
		return nil
	},
}
