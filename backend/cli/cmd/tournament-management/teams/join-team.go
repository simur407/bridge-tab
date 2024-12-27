package teams

import (
	tournament_management "bridge-tab/internal/tournament-management/application"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"fmt"

	"github.com/spf13/cobra"
)

var joinTeamId string
var joinTeamContestantId string

var joinTeamCmd = func(TournamentRepository *tournament_domain.TournamentRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "join",
		Short:        "Joins contestant to a team in the tournament",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			command := &tournament_management.JoinTeamCommand{
				TournamentId: teamsTournamentId,
				TeamId:       joinTeamId,
				ContestantId: joinTeamContestantId,
			}

			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			fmt.Println("Contestant { Id:", command.ContestantId, " } joined Team { Id:", command.TeamId, " } in Tournament { Id:", command.TournamentId, "}")
			return nil
		},
	}

	command.Flags().StringVarP(&joinTeamId, "id", "i", "", "team id")
	command.MarkFlagRequired("id")
	command.Flags().StringVarP(&joinTeamContestantId, "contestantId", "c", "", "contestant id")
	command.MarkFlagRequired("contestantId")

	return command
}
