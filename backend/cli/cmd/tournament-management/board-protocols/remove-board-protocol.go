package board_protocols

import (
	application "bridge-tab/internal/tournament-management/application/command"
	domain "bridge-tab/internal/tournament-management/domain"
	"fmt"

	"github.com/spf13/cobra"
)

var boardProtocolNo int

var removeBoardProtocolCmd = func(TournamentRepository *domain.TournamentRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "remove",
		Short:        "Removes board protocol from the tournament",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			command := &application.RemoveBoardProtocolCommand{TournamentId: boardProtocolsTournamentId, BoardNo: boardProtocolNo}

			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			fmt.Println("Board Protocol { BoardNo:", command.BoardNo, " } removed from Tournament { Id:", command.TournamentId, "}")
			return nil
		},
	}

	command.Flags().IntVarP(&boardProtocolNo, "no", "n", 0, "board protocol number")
	command.MarkFlagRequired("no")

	return command
}
