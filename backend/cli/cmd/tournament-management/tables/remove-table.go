package tables

import (
	tournament_management "bridge-tab/internal/tournament-management/application/command"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
	"fmt"

	"github.com/spf13/cobra"
)

var removeTableId string

var removeTableCmd = func(TournamentRepository *tournament_domain.TournamentRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "remove",
		Short:        "Removes table from the tournament",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			command := &tournament_management.RemoveTableCommand{TournamentId: tablesTournamentId, TableId: removeTableId}

			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			fmt.Println("Table { Id:", command.TableId, " } removed from Tournament { Id:", command.TournamentId, "}")
			return nil
		},
	}

	command.Flags().StringVarP(&removeTableId, "id", "i", "", "table id")
	command.MarkFlagRequired("id")

	return command
}
