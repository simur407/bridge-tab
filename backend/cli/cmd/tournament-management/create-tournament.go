package tournament_management

import (
	"fmt"

	tournament "bridge-tab/internal/tournament-management/application"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var name string

var createTournamentCmd = func(TournamentRepository *tournament_domain.TournamentRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "create",
		Short:        "Creates a new Tournament",
		Long:         `This command creates a new Tournament with given name. The name should be unique.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			command := &tournament.CreateTournamentCommand{TournamentId: uuid.New().String(), Name: name}

			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			fmt.Println("Created Tournament { Id:", command.TournamentId, "Name:", name, "}")
			return nil
		},
	}

	command.Flags().StringVarP(&name, "name", "n", "", "unique Tournament name")
	command.MarkFlagRequired("name")

	return command
}
