package schedule

import (
	domain "bridge-tab/internal/tournament-management/domain"

	"github.com/spf13/cobra"
)

var scheduleTournamentId string

var ScheduleCmd = func(
	TournamentRepository *domain.TournamentRepository,
	TeamReadRepository *domain.TeamReadRepository,
	TableReadRepository *domain.TableReadRepository,
	SetReadRepository *domain.SetReadRepository,
	ScheduleReadRepository *domain.ScheduleReadRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:   "schedule",
		Short: "Responsible for managing the tournament schedule",
		Long:  "Schedule says which teams play which set at which table in a round.",
	}

	command.PersistentFlags().StringVarP(&scheduleTournamentId, "tournamentId", "i", "", "tournament id")
	command.MarkPersistentFlagRequired("tournamentId")
	command.AddCommand(
		createScheduleCmd(TournamentRepository, TeamReadRepository, TableReadRepository, SetReadRepository),
		removeScheduleCmd(TournamentRepository, TableReadRepository),
		listScheduleCmd(ScheduleReadRepository),
	)

	return command
}
