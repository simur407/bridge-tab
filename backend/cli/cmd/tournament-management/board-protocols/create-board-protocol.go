package board_protocols

import (
	application "bridge-tab/internal/tournament-management/application"
	domain "bridge-tab/internal/tournament-management/domain"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var boardNo int
var vulnerable string

var createBoardProtocolCmd = func(
	TournamentRepository *domain.TournamentRepository,
	TeamReadRepository *domain.TeamReadRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:          "create",
		Short:        "Creates board protocol",
		SilenceUsage: true,
		Args:         cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var vulnerability domain.Vulnerable
			if boardNo <= 0 {
				return fmt.Errorf("board number is required")
			}

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

			var teamPairs []struct {
				NS string
				EW string
			}
			for _, arg := range args {
				teamNames := strings.SplitN(arg, ";", 2)
				if len(teamNames) != 2 {
					return fmt.Errorf("invalid argument, expected format: {NS Team Name};{EW Team Name}")
				}

				nsTeamByNameQuery := application.GetTeamByNameQuery{
					TournamentId: boardProtocolsTournamentId,
					Name:         teamNames[0],
				}
				nsTeam, err := nsTeamByNameQuery.Execute(*TeamReadRepository)
				if err != nil {
					return err
				}

				ewTeamByNameQuery := application.GetTeamByNameQuery{
					TournamentId: boardProtocolsTournamentId,
					Name:         teamNames[1],
				}
				ewTeam, err := ewTeamByNameQuery.Execute(*TeamReadRepository)
				if err != nil {
					return err
				}

				teamPairs = append(teamPairs, struct {
					NS string
					EW string
				}{
					NS: nsTeam.Id,
					EW: ewTeam.Id,
				})
			}

			command := &application.CreateBoardProtocol{TournamentId: boardProtocolsTournamentId, BoardNo: boardNo, Vulnerable: int(vulnerability), TeamPairs: teamPairs}

			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			return nil
		},
	}

	command.Flags().IntVarP(&boardNo, "boardNo", "n", 0, "unique Board number")
	command.MarkFlagRequired("number")
	command.Flags().StringVarP(&vulnerable, "vulnerability", "v", "None",
		"the state of vulnerability of the board. Accepted values: None, NS, EW, Both. Default is None")

	return command
}
