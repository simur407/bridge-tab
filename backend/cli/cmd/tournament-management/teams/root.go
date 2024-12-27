package teams

import (
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/spf13/cobra"
)

var teamsTournamentId string

var TeamsCmd = func(TournamentRepository *tournament_domain.TournamentRepository, TeamRepository *tournament_domain.TeamReadRepository) *cobra.Command {
	command := &cobra.Command{
		Use:   "team",
		Short: "Resposible for managing Teams in Tournaments",
		Long:  "Team Management allows organizers or umpires to manage Teams in Tournaments like: create, delete, join, leave, etc.",
	}

	command.PersistentFlags().StringVarP(&teamsTournamentId, "tournamentId", "t", "", "tournament id")
	command.MarkPersistentFlagRequired("tournamentId")
	command.AddCommand(
		createTeamCmd(TournamentRepository),
		removeTeamCmd(TournamentRepository),
		joinTeamCmd(TournamentRepository),
		leaveTeamCmd(TournamentRepository),
		listTeamsCmd(TeamRepository),
	)

	return command
}
