package teams

import (
	tournament_management "bridge-tab/internal/tournament-management/application/command"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
	"fmt"

	"github.com/spf13/cobra"
)

var leaveTeamId string
var leaveTeamContestantId string

var leaveTeamCmd = func(TournamentRepository *tournament_domain.TournamentRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "leave",
		Short:        "Leaves team from the tournament",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			command := &tournament_management.LeaveTeamCommand{
				TournamentId: teamsTournamentId,
				TeamId:       leaveTeamId,
				ContestantId: leaveTeamContestantId,
			}

			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			fmt.Println("Contestant { Id:", command.ContestantId, " } leave Team { Id:", command.TeamId, " } in Tournament { Id:", command.TournamentId, "}")
			return nil
		},
	}

	command.Flags().StringVarP(&leaveTeamId, "id", "i", "", "team id")
	command.MarkFlagRequired("id")
	command.Flags().StringVarP(&leaveTeamContestantId, "contestantId", "c", "", "contestant id")
	command.MarkFlagRequired("contestantId")

	return command
}
