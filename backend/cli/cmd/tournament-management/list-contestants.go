package tournament_management

import (
	"fmt"
	"strings"

	tournament "bridge-tab/internal/tournament-management/application"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/spf13/cobra"
)

var contestantsTournamentId string

var listContestantsCmd = func(tournamentReadRepository *tournament_domain.TournamentReadRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "contestants",
		Short:        "Lists existing contestants in a tournament",
		Long:         `This command lists all existing contestants in a tournament.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			query := tournament.ListContestantsQuery{TournamentId: contestantsTournamentId}

			results, err := query.Execute(*tournamentReadRepository)
			if err != nil {
				return err
			}

			fmt.Printf("%-36v\n", "Id")
			fmt.Println(strings.Repeat("-", 36))
			for _, Contestant := range results {
				fmt.Printf("%-36v\n", Contestant.Id)
				fmt.Println(strings.Repeat("-", 36))
			}
			return nil
		},
	}

	command.Flags().StringVarP(&contestantsTournamentId, "id", "i", "", "tournament id")
	command.MarkFlagRequired("id")

	return command
}
