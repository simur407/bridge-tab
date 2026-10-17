package schedule

import (
	"fmt"
	"strings"

	application_query "bridge-tab/internal/tournament-management/application/query"
	domain "bridge-tab/internal/tournament-management/domain"

	"github.com/spf13/cobra"
)

var listScheduleCmd = func(ScheduleReadRepository *domain.ScheduleReadRepository) *cobra.Command {
	return &cobra.Command{
		Use:          "list",
		Short:        "Lists the tournament schedule",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := (&application_query.ListScheduleQuery{TournamentId: scheduleTournamentId}).Execute(*ScheduleReadRepository)
			if err != nil {
				return err
			}

			fmt.Printf("%-8v | %-8v | %-8v | %-8v | %-10v\n", "Round", "Table", "NS", "EW", "Set")
			fmt.Println(strings.Repeat("-", 52))
			for _, entry := range entries {
				fmt.Printf("%-8v | %-8v | %-8v | %-8v | %-10v\n", entry.RoundNo, entry.TableNumber, entry.NsNumber, entry.EwNumber, entry.SetLabel)
				fmt.Println(strings.Repeat("-", 52))
			}
			return nil
		},
	}
}
