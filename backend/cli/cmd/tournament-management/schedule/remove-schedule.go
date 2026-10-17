package schedule

import (
	"fmt"

	application "bridge-tab/internal/tournament-management/application/command"
	application_query "bridge-tab/internal/tournament-management/application/query"
	domain "bridge-tab/internal/tournament-management/domain"

	"github.com/spf13/cobra"
)

var removeRoundNo int
var removeTableNumber int

var removeScheduleCmd = func(
	TournamentRepository *domain.TournamentRepository,
	TableReadRepository *domain.TableReadRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:          "remove",
		Short:        "Removes a schedule entry for a round and table",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			table, err := (&application_query.GetTableByNumberQuery{
				TournamentId: scheduleTournamentId,
				Number:       removeTableNumber,
			}).Execute(*TableReadRepository)
			if err != nil {
				return fmt.Errorf("table %d: %w", removeTableNumber, err)
			}

			command := &application.RemoveScheduleCommand{
				TournamentId: scheduleTournamentId,
				RoundNo:      removeRoundNo,
				TableId:      table.Id,
			}
			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			fmt.Println("Schedule { Round:", removeRoundNo, ", Table:", removeTableNumber, "} removed from Tournament { Id:", scheduleTournamentId, "}")
			return nil
		},
	}

	command.Flags().IntVarP(&removeRoundNo, "round", "r", 0, "round number")
	command.MarkFlagRequired("round")
	command.Flags().IntVarP(&removeTableNumber, "table", "t", 0, "table number")
	command.MarkFlagRequired("table")

	return command
}
