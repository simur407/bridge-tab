package score_registration

import (
	score_query "bridge-tab/internal/score-registration/application/query"
	score_domain "bridge-tab/internal/score-registration/domain"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var standingsTournamentId string
var standingsMethod string

var standingsCmd = func(tournamentScoreReadRepository *score_domain.TournamentScoreReadRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "standings",
		Short:        "Lists standings for a tournament and scoring method",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			method := score_domain.ScoringMethodName(standingsMethod)
			switch method {
			case score_domain.MethodMatchpoints, score_domain.MethodIMPs:
			default:
				return fmt.Errorf("unsupported scoring method %q (allowed: %s, %s)",
					standingsMethod, score_domain.MethodMatchpoints, score_domain.MethodIMPs)
			}

			query := score_query.GetStandingsQuery{
				TournamentId: standingsTournamentId,
				Method:       method,
			}
			standings, err := query.Execute(*tournamentScoreReadRepository)
			if err != nil {
				return err
			}

			fmt.Printf("%-6v | %-8v | %-20v | %-10v | %-10v\n",
				"Rank", "Number", "Name", "Total", "Percentage")
			fmt.Println(strings.Repeat("-", 66))
			for _, standing := range standings {
				fmt.Printf("%-6v | %-8v | %-20v | %-10.1f | %-10.2f\n",
					standing.Rank,
					standing.TeamNumber,
					standing.TeamName,
					standing.TotalScore,
					standing.Percentage,
				)
				fmt.Println(strings.Repeat("-", 66))
			}
			return nil
		},
	}

	command.Flags().StringVarP(&standingsTournamentId, "tournament", "t", "", "tournament id")
	command.Flags().StringVarP(&standingsMethod, "method", "m", string(score_domain.MethodMatchpoints), "scoring method (matchpoints|imps)")
	command.MarkFlagRequired("tournament")

	return command
}
