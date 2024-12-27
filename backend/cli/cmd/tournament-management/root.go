package tournament_management

import (
	board_protocols "bridge-tab/cli/cmd/tournament-management/board-protocols"
	teams "bridge-tab/cli/cmd/tournament-management/teams"
	tournament "bridge-tab/internal/tournament-management/domain"

	"github.com/spf13/cobra"
)

var TournamentManagementCmd = func(
	tournamentRepository *tournament.TournamentRepository,
	tournamentReadRepository *tournament.TournamentReadRepository,
	teamReadRepository *tournament.TeamReadRepository,
	boardProtocolReadRepository *tournament.BoardProtocolReadRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:   "tournament",
		Short: "Responsible for managing Tournaments",
		Long:  "Tournament Management allows organizers or umpires to manage Tournaments like: create, delete, add deal protocols, etc.",
	}

	command.AddCommand(
		createTournamentCmd(tournamentRepository),
		removeTournamentCmd(tournamentRepository),
		listTournamentsCmd(tournamentReadRepository),
		startTournamentCmd(tournamentRepository),
		joinTournamentCmd(tournamentRepository),
		leaveTournamentCmd(tournamentRepository),
		listContestantsCmd(tournamentReadRepository),
		teams.TeamsCmd(tournamentRepository, teamReadRepository),
		board_protocols.BoardProtocolsCmd(tournamentRepository, teamReadRepository, boardProtocolReadRepository),
	)

	return command
}
