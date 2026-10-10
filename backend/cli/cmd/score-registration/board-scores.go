package score_registration

import (
	score_query "bridge-tab/internal/score-registration/application/query"
	score_domain "bridge-tab/internal/score-registration/domain"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var boardScoresTournamentId string
var boardScoresMethod string
var boardScoresDealNo int

var boardScoresCmd = func(tournamentScoreReadRepository *score_domain.TournamentScoreReadRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "board-scores",
		Short:        "Lists per-board scores for a tournament and scoring method",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			method := score_domain.ScoringMethodName(boardScoresMethod)
			switch method {
			case score_domain.MethodMatchpoints, score_domain.MethodIMPs:
			default:
				return fmt.Errorf("unsupported scoring method %q (allowed: %s, %s)",
					boardScoresMethod, score_domain.MethodMatchpoints, score_domain.MethodIMPs)
			}

			query := score_query.ListBoardScoresQuery{
				TournamentId: boardScoresTournamentId,
				Method:       method,
			}
			if boardScoresDealNo > 0 {
				dealNo := boardScoresDealNo
				query.DealNo = &dealNo
			}

			scores, err := query.Execute(*tournamentScoreReadRepository)
			if err != nil {
				return err
			}

			fmt.Printf("%-6v | %-6v | %-6v | %-10v | %-10v | %-10v\n",
				"Deal", "NS", "EW", "NS Points", "NS Award", "EW Award")
			fmt.Println(strings.Repeat("-", 64))
			for _, score := range scores {
				fmt.Printf("%-6v | %-6v | %-6v | %-10v | %-10.1f | %-10.1f\n",
					score.DealNo,
					score.NsTeamNumber,
					score.EwTeamNumber,
					score.NsPoints,
					score.NsAward,
					score.EwAward,
				)
				fmt.Println(strings.Repeat("-", 64))
			}
			return nil
		},
	}

	command.Flags().StringVarP(&boardScoresTournamentId, "tournament", "t", "", "tournament id")
	command.Flags().StringVarP(&boardScoresMethod, "method", "m", string(score_domain.MethodMatchpoints), "scoring method (matchpoints|imps)")
	command.MarkFlagRequired("tournament")
	command.Flags().IntVarP(&boardScoresDealNo, "deal", "d", 0, "optional deal number filter")

	return command
}
