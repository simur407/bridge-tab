package board_protocols

import (
	application "bridge-tab/internal/tournament-management/application/command"
	application_query "bridge-tab/internal/tournament-management/application/query"
	domain "bridge-tab/internal/tournament-management/domain"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var boardNo int
var vulnerable string

var createBoardProtocolCmd = func(
	TournamentRepository *domain.TournamentRepository,
	TeamReadRepository *domain.TeamReadRepository,
	TableReadRepository *domain.TableReadRepository,
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
				Table *string
				NS    string
				EW    string
			}
			for _, arg := range args {
				pairArg := arg
				var tableId *string

				if tablePart, rest, ok := strings.Cut(arg, ":"); ok {
					tableNumber, err := strconv.Atoi(tablePart)
					if err != nil {
						return fmt.Errorf("invalid argument, expected format: [{table number}:]{NS team number};{EW team number}")
					}

					tableByNumberQuery := application_query.GetTableByNumberQuery{
						TournamentId: boardProtocolsTournamentId,
						Number:       tableNumber,
					}
					table, err := tableByNumberQuery.Execute(*TableReadRepository)
					if err != nil {
						return fmt.Errorf("table %d: %w", tableNumber, err)
					}
					tableId = &table.Id
					pairArg = rest
				}

				teamNumbers := strings.SplitN(pairArg, ";", 2)
				if len(teamNumbers) != 2 {
					return fmt.Errorf("invalid argument, expected format: [{table number}:]{NS team number};{EW team number}")
				}

				nsNumber, err := strconv.Atoi(teamNumbers[0])
				if err != nil {
					return fmt.Errorf("invalid NS team number %q", teamNumbers[0])
				}
				ewNumber, err := strconv.Atoi(teamNumbers[1])
				if err != nil {
					return fmt.Errorf("invalid EW team number %q", teamNumbers[1])
				}

				nsTeamByNumberQuery := application_query.GetTeamByNumberQuery{
					TournamentId: boardProtocolsTournamentId,
					Number:       nsNumber,
				}
				nsTeam, err := nsTeamByNumberQuery.Execute(*TeamReadRepository)
				if err != nil {
					return fmt.Errorf("NS team %d: %w", nsNumber, err)
				}

				ewTeamByNumberQuery := application_query.GetTeamByNumberQuery{
					TournamentId: boardProtocolsTournamentId,
					Number:       ewNumber,
				}
				ewTeam, err := ewTeamByNumberQuery.Execute(*TeamReadRepository)
				if err != nil {
					return fmt.Errorf("EW team %d: %w", ewNumber, err)
				}

				teamPairs = append(teamPairs, struct {
					Table *string
					NS    string
					EW    string
				}{
					Table: tableId,
					NS:    nsTeam.Id,
					EW:    ewTeam.Id,
				})
			}

			command := &application.CreateBoardProtocol{TournamentId: boardProtocolsTournamentId, BoardNo: boardNo, Vulnerable: int(vulnerability), TeamPairs: teamPairs}

			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			fmt.Println("Board Protocol { BoardNo:", command.BoardNo, ", Vulnerable:", vulnerable, ", args:", args, "} created in Tournament { Id:", command.TournamentId, "}")

			return nil
		},
	}

	command.Flags().IntVarP(&boardNo, "boardNo", "n", 0, "unique Board number")
	command.MarkFlagRequired("number")
	command.Flags().StringVarP(&vulnerable, "vulnerability", "v", "None",
		"the state of vulnerability of the board. Accepted values: None, NS, EW, Both. Default is None")

	return command
}
