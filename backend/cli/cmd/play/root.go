package play

import (
	play_domain "bridge-tab/internal/play/domain"
	score_domain "bridge-tab/internal/score-registration/domain"

	"github.com/spf13/cobra"
)

var PlayCmd = func(
	boardResultRepository *play_domain.BoardResultRepository,
	boardResultReadRepository *play_domain.BoardResultReadRepository,
	playSessionRepository *play_domain.PlaySessionRepository,
	playedResultRepository *score_domain.PlayedResultRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:   "play",
		Short: "Records board results and registers them for scoring",
		Long:  "Play records the result of one board at one table. Each board result is registered for scoring as soon as it is saved.",
	}

	command.AddCommand(
		recordBoardResultCmd(boardResultRepository, playSessionRepository, playedResultRepository),
		editBoardResultCmd(boardResultRepository, boardResultReadRepository, playSessionRepository, playedResultRepository),
		listBoardResultsCmd(boardResultReadRepository),
	)

	return command
}
