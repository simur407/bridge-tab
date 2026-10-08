package sets

import (
	application "bridge-tab/internal/tournament-management/application/command"
	domain "bridge-tab/internal/tournament-management/domain"
	"fmt"

	"github.com/spf13/cobra"
)

var removeSetId string

var removeSetCmd = func(TournamentRepository *domain.TournamentRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "remove",
		Short:        "Removes a set grouping (does not delete board protocols)",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			command := &application.RemoveSetCommand{
				TournamentId: setsTournamentId,
				SetId:        removeSetId,
			}
			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			fmt.Println("Set { Id:", removeSetId, "} removed from Tournament { Id:", setsTournamentId, "}")
			return nil
		},
	}

	command.Flags().StringVarP(&removeSetId, "id", "s", "", "set id")
	command.MarkFlagRequired("id")

	return command
}
