package tables

import (
	tournament_management "bridge-tab/internal/tournament-management/application/command"
	tournament_query "bridge-tab/internal/tournament-management/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var tableNumber int

var createTableCmd = func(
	TournamentRepository *tournament_domain.TournamentRepository,
	TableReadRepository *tournament_domain.TableReadRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:          "create",
		Short:        "Creates table in a tournament",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			command := &tournament_management.CreateTableCommand{
				TournamentId: tablesTournamentId,
				TableId:      uuid.New().String(),
				Number:       tableNumber,
			}

			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			table, err := (&tournament_query.GetTableByIdQuery{
				TournamentId: command.TournamentId,
				TableId:      command.TableId,
			}).Execute(*TableReadRepository)
			if err != nil {
				return err
			}

			fmt.Println("Table { Id:", table.Id, ", Number:", table.Number, " } created in Tournament { Id:", command.TournamentId, "}")
			return nil
		},
	}

	command.Flags().IntVar(&tableNumber, "number", 0, "table number; omitted values continue from 1")

	return command
}
