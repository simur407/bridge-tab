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

			fmt.Printf("%-10v | %-15v\n", "Board No", "Vulnerability")
			fmt.Println(strings.Repeat("-", 28))
			for _, BoardProtocol := range results {
				fmt.Printf("%-10v | %-15v\n", BoardProtocol.BoardNo, BoardProtocol.Vulnerable)
				fmt.Println(strings.Repeat("-", 28))
			}

			return nil
		},
	}

	return command
}
