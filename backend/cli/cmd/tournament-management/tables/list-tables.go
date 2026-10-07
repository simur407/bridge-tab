package tables

import (
	"fmt"
	"strings"

	tournament "bridge-tab/internal/tournament-management/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/spf13/cobra"
)

var listTablesCmd = func(TableReadRepository *tournament_domain.TableReadRepository) *cobra.Command {
	return &cobra.Command{
		Use:          "list",
		Short:        "Lists existing tables in a tournament",
		Long:         `This command lists all existing tables in a tournament.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			query := tournament.ListTablesQuery{TournamentId: tablesTournamentId}

			results, err := query.Execute(*TableReadRepository)
			if err != nil {
				return err
			}

			fmt.Printf("%-36v | %-8v\n", "Id", "Number")
			fmt.Println(strings.Repeat("-", 47))
			for _, table := range results {
				fmt.Printf("%-36v | %-8v\n", table.Id, table.Number)
				fmt.Println(strings.Repeat("-", 47))
			}
			return nil
		},
	}
}
