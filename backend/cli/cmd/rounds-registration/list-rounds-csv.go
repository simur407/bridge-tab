package rounds_registration

import (
	rounds_registration_query "bridge-tab/internal/rounds-registration/application/query"
	rounds_registration_domain "bridge-tab/internal/rounds-registration/domain"
	"fmt"

	"github.com/spf13/cobra"
)

var listRoundsGameSessionId string

var listRoundsCsvCmd = func(gameSessionReadRepository *rounds_registration_domain.GameSessionReadRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "list",
		Short:        "Lists all rounds in csv format",
		Long:         `This command lists all rounds in csv format`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			query := rounds_registration_query.ListRoundsQuery{GameSessionId: listRoundsGameSessionId}

			results, err := query.Execute(*gameSessionReadRepository)
			if err != nil {
				return err
			}

			fmt.Println("Deal No;NS;EW;Contract;Declarer;Tricks;Opening Lead")
			for _, Round := range results {
				fmt.Printf("%d;%s;%s;%s;%s;%d;%s\n", Round.DealNo, Round.NsTeamName, Round.EwTeamName, Round.Contract, Round.Declarer, Round.Tricks, Round.OpeningLead)
			}

			return nil

		},
	}

	command.Flags().StringVarP(&listRoundsGameSessionId, "id", "i", "", "tournament id")
	command.MarkFlagRequired("id")

	return command
}
