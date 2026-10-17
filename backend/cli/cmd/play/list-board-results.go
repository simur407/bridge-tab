package play

import (
	"fmt"

	play_query "bridge-tab/internal/play/application/query"
	play_domain "bridge-tab/internal/play/domain"

	"github.com/spf13/cobra"
)

var listPlaySessionId string

var listBoardResultsCmd = func(boardResultReadRepository *play_domain.BoardResultReadRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "list",
		Short:        "Lists board results in csv format",
		Long:         "Lists every recorded board result in csv format.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			playSessionId := listPlaySessionId
			results, err := (&play_query.ListBoardResultsQuery{PlaySessionId: playSessionId}).Execute(*boardResultReadRepository)
			if err != nil {
				return err
			}

			fmt.Println("Deal No;NS;EW;Contract;Declarer;Tricks;Opening Lead")
			for _, result := range results {
				fmt.Printf("%d;%d;%d;%s;%s;%d;%s\n", result.DealNo, result.NsTeamNumber, result.EwTeamNumber, result.Contract, result.Declarer, result.Tricks, result.OpeningLead)
			}
			return nil
		},
	}

	command.Flags().StringVarP(&listPlaySessionId, "id", "i", "", "play session id")
	command.MarkFlagRequired("id")

	return command
}
