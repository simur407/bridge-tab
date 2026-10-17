package board_protocols

import (
	application "bridge-tab/internal/tournament-management/application/command"
	domain "bridge-tab/internal/tournament-management/domain"
	"fmt"

	"github.com/spf13/cobra"
)

var boardNo int
var vulnerable string

var createBoardProtocolCmd = func(TournamentRepository *domain.TournamentRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "create",
		Short:        "Creates board protocol",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if boardNo <= 0 {
				return fmt.Errorf("board number is required")
			}

			var vulnerability domain.Vulnerable
			switch vulnerable {
			case "None":
				vulnerability = domain.None
			case "NS":
				vulnerability = domain.NS
			case "EW":
				vulnerability = domain.EW
			case "Both":
				vulnerability = domain.Both
			default:
				return fmt.Errorf("invalid vulnerability")
			}

			command := &application.CreateBoardProtocol{
				TournamentId: boardProtocolsTournamentId,
				BoardNo:      boardNo,
				Vulnerable:   int(vulnerability),
			}

			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			fmt.Println("Board Protocol { BoardNo:", command.BoardNo, ", Vulnerable:", vulnerable, "} created in Tournament { Id:", command.TournamentId, "}")
			return nil
		},
	}

	command.Flags().IntVarP(&boardNo, "boardNo", "n", 0, "unique Board number")
	command.MarkFlagRequired("boardNo")
	command.Flags().StringVarP(&vulnerable, "vulnerability", "v", "None",
		"the state of vulnerability of the board. Accepted values: None, NS, EW, Both. Default is None")

	return command
}
