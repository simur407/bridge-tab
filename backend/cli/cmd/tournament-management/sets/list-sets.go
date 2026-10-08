package sets

import (
	application_query "bridge-tab/internal/tournament-management/application/query"
	domain "bridge-tab/internal/tournament-management/domain"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var listSetsCmd = func(SetReadRepository *domain.SetReadRepository) *cobra.Command {
	return &cobra.Command{
		Use:          "list",
		Short:        "Lists sets in a tournament",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			query := application_query.ListSetsQuery{TournamentId: setsTournamentId}
			sets, err := query.Execute(*SetReadRepository)
			if err != nil {
				return err
			}

			fmt.Printf("%-36v | %-10v | %-20v\n", "Id", "Label", "Boards")
			fmt.Println(strings.Repeat("-", 74))
			for _, set := range sets {
				boardParts := make([]string, len(set.BoardNos))
				for i, boardNo := range set.BoardNos {
					boardParts[i] = fmt.Sprintf("%d", boardNo)
				}
				fmt.Printf("%-36v | %-10v | %-20v\n", set.Id, set.Label, strings.Join(boardParts, ","))
				fmt.Println(strings.Repeat("-", 74))
			}
			return nil
		},
	}
}
