package tournament_management

import (
	board_protocols "bridge-tab/cli/cmd/tournament-management/board-protocols"
	sets "bridge-tab/cli/cmd/tournament-management/sets"
	tables "bridge-tab/cli/cmd/tournament-management/tables"
	teams "bridge-tab/cli/cmd/tournament-management/teams"
	rounds_registration "bridge-tab/internal/rounds-registration/domain"
	score_registration "bridge-tab/internal/score-registration/domain"
	tournament "bridge-tab/internal/tournament-management/domain"

	"github.com/spf13/cobra"
)

var TournamentManagementCmd = func(
	tournamentRepository *tournament.TournamentRepository,
	tournamentReadRepository *tournament.TournamentReadRepository,
	teamReadRepository *tournament.TeamReadRepository,
	tableReadRepository *tournament.TableReadRepository,
	boardProtocolReadRepository *tournament.BoardProtocolReadRepository,
	setReadRepository *tournament.SetReadRepository,
	gameSessionRepository *rounds_registration.GameSessionRepository,
	playedResultRepository *score_registration.PlayedResultRepository,
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
		startTournamentCmd(tournamentRepository, teamReadRepository, boardProtocolReadRepository, gameSessionRepository),
		finishTournamentCmd(tournamentRepository, gameSessionRepository, boardProtocolReadRepository, playedResultRepository),
		joinTournamentCmd(tournamentRepository),
		leaveTournamentCmd(tournamentRepository),
		listContestantsCmd(tournamentReadRepository),
		teams.TeamsCmd(tournamentRepository, teamReadRepository),
		tables.TablesCmd(tournamentRepository, tableReadRepository),
		board_protocols.BoardProtocolsCmd(tournamentRepository, teamReadRepository, tableReadRepository, boardProtocolReadRepository),
		sets.SetsCmd(tournamentRepository, setReadRepository),
	)

	return command
}
