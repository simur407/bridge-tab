package tournament_management

import (
	board_protocols "bridge-tab/cli/cmd/tournament-management/board-protocols"
	sets "bridge-tab/cli/cmd/tournament-management/sets"
	tables "bridge-tab/cli/cmd/tournament-management/tables"
	teams "bridge-tab/cli/cmd/tournament-management/teams"
	play "bridge-tab/internal/play/domain"
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
	playedResultRepository *score_registration.PlayedResultRepository,
	boardResultReadRepository *play.BoardResultReadRepository,
	playSessionRepository *play.PlaySessionRepository,
	tournamentScoreRepository *score_registration.TournamentScoreRepository,
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
		startTournamentCmd(tournamentRepository, teamReadRepository, boardProtocolReadRepository, playSessionRepository),
		finishTournamentCmd(tournamentRepository, playSessionRepository, boardResultReadRepository, playedResultRepository, tournamentScoreRepository),
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
