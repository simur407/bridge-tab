package tables

import (
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/spf13/cobra"
)

var tablesTournamentId string

var TablesCmd = func(
	TournamentRepository *tournament_domain.TournamentRepository,
	TableReadRepository *tournament_domain.TableReadRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:   "table",
		Short: "Responsible for managing Tables in Tournaments",
		Long:  "Table Management allows organizers or umpires to manage Tables in Tournaments like: create, delete, list, etc.",
	}

	command.PersistentFlags().StringVarP(&tablesTournamentId, "tournamentId", "t", "", "tournament id")
	command.MarkPersistentFlagRequired("tournamentId")
	command.AddCommand(
		createTableCmd(TournamentRepository, TableReadRepository),
		removeTableCmd(TournamentRepository),
		listTablesCmd(TableReadRepository),
	)

	return command
}
