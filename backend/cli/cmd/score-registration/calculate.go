package score_registration

import (
	score_app "bridge-tab/internal/score-registration/application"
	score_domain "bridge-tab/internal/score-registration/domain"
	"fmt"

	"github.com/spf13/cobra"
)

var calculateTournamentId string
var calculateMethod string

var calculateCmd = func(
	playedResultRepository *score_domain.PlayedResultRepository,
	tournamentScoreRepository *score_domain.TournamentScoreRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:   "calculate",
		Short: "Calculates standings for a scoring method",
		Long: `Calculates board awards and standings from transferred played results
	using the chosen scoring method (matchpoints or imps).`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			method := score_domain.ScoringMethodName(calculateMethod)
			switch method {
			case score_domain.MethodMatchpoints, score_domain.MethodIMPs:
			default:
				return fmt.Errorf("unsupported scoring method %q (allowed: %s, %s)",
					calculateMethod, score_domain.MethodMatchpoints, score_domain.MethodIMPs)
			}

			command := &score_app.CalculateStandingsCommand{
				TournamentId: calculateTournamentId,
				Method:       method,
			}
			if err := command.Execute(*playedResultRepository, *tournamentScoreRepository); err != nil {
				return err
			}

			fmt.Println("Calculated standings { Tournament:", calculateTournamentId, "Method:", method, "}")
			return nil
		},
	}

	command.Flags().StringVarP(&calculateTournamentId, "tournament", "t", "", "tournament id")
	command.Flags().StringVarP(&calculateMethod, "method", "m", string(score_domain.MethodMatchpoints), "scoring method (matchpoints|imps)")
	command.MarkFlagRequired("tournament")

	return command
}
