package tournament_management

import (
	"fmt"
	"strings"

	tournament "bridge-tab/internal/tournament-management/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/spf13/cobra"
)

var listTournamentsCmd = func(TournamentReadRepository *tournament_domain.TournamentReadRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "list",
		Short:        "Lists existing tournaments",
		Long:         `This command lists all existing tournaments.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			query := tournament.ListTournamentsQuery{}

			results, err := query.Execute(*TournamentReadRepository)
			if err != nil {
				return err
			}

			fmt.Printf("%-36v | %-15v | %-30v\n", "Id", "Name", "Started At")
			fmt.Println(strings.Repeat("-", 88))
			for _, Tournament := range results {
				fmt.Printf("%-36v | %-15v | %-30v\n", Tournament.Id, Tournament.Name, Tournament.StartedAt)
				fmt.Println(strings.Repeat("-", 88))
			}
			return nil
		},
	}

	return command
}
