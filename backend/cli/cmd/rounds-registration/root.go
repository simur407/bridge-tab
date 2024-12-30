package rounds_registration

import (
	rounds_registration "bridge-tab/internal/rounds-registration/domain"
	tournament_management "bridge-tab/internal/tournament-management/domain"

	"github.com/spf13/cobra"
)

var RoundsRegistrationCmd = func(
	gameSessionRepository *rounds_registration.GameSessionRepository,
	teamRepository *tournament_management.TeamReadRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:   "round",
		Short: "Responsible for managing Rounds",
		Long:  "Rounds registration allows organizers or umpires to manage Rounds registration like: start session, add played round",
	}

	command.AddCommand(
		playRoundCmd(gameSessionRepository, teamRepository),
	)

	return command
}
