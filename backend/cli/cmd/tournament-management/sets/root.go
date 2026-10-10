package sets

import (
	domain "bridge-tab/internal/tournament-management/domain"

	"github.com/spf13/cobra"
)

var setsTournamentId string

var SetsCmd = func(
	TournamentRepository *domain.TournamentRepository,
	SetReadRepository *domain.SetReadRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:   "set",
		Short: "Responsible for managing Sets (komplety) in Tournaments",
		Long:  "Set Management groups board protocols that share the same team pairings for reports.",
	}

	command.PersistentFlags().StringVarP(&setsTournamentId, "tournamentId", "i", "", "tournament id")
	command.MarkPersistentFlagRequired("tournamentId")
	command.AddCommand(
		createSetCmd(TournamentRepository),
		removeSetCmd(TournamentRepository),
		listSetsCmd(SetReadRepository),
	)

	return command
}
