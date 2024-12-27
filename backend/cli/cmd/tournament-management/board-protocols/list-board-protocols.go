package board_protocols

import (
	application "bridge-tab/internal/tournament-management/application"
	domain "bridge-tab/internal/tournament-management/domain"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var listBoardProtocolsCmd = func(boardProtocolReadRepository *domain.BoardProtocolReadRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "list",
		Short:        "Lists existing board protocols",
		Long:         `This command lists all existing board protocols.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			query := application.ListBoardProtocolsQuery{
				TournamentId: boardProtocolsTournamentId,
			}

			results, err := query.Execute(*boardProtocolReadRepository)
			if err != nil {
				return err
			}

			fmt.Printf("%-10v | %-15v | %-36v | %-36v\n", "Board No", "Vulnerability", "NS", "EW")
			fmt.Println(strings.Repeat("-", 106))
			for _, BoardProtocol := range results {
				row := fmt.Sprintf("%-10v | %-15v | ", BoardProtocol.BoardNo, BoardProtocol.Vulnerable)

				for i, TeamPair := range BoardProtocol.TeamPairs {
					if i == 0 {
						row += fmt.Sprintf("%-36v | %-36v\n", TeamPair.NS, TeamPair.EW)
					} else {
						row += fmt.Sprintf("%-10v   %-15v | %-36v | %-36v\n", "", "", TeamPair.NS, TeamPair.EW)
					}
				}

				lastRune := []rune(row)[len(row)-1]
				if lastRune != '\n' {
					row += "\n"
				}
				fmt.Print(row)
				fmt.Println(strings.Repeat("-", 106))

			}

			return nil
		},
	}

	return command
}
