package tournament_management

import (
	tournament_management "bridge-tab/internal/tournament-management/application/command"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var joinTournamentId string
var joinContestantId string

var joinTournamentCmd = func(TournamentRepository *tournament_domain.TournamentRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "join",
		Short:        "Joins contestant to a tournament",
		Long:         `This command joins contestant to a tournament.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if joinContestantId == "" {
				joinContestantId = uuid.New().String()
			}
			command := &tournament_management.JoinTournamentCommand{TournamentId: joinTournamentId, ContestantId: joinContestantId}

			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			fmt.Println("Contestant { Id:", command.ContestantId, " } joined Tournament { Id:", command.TournamentId, "}")
			return nil
		},
	}

	command.Flags().StringVarP(&joinTournamentId, "id", "i", "", "tournament id")
	command.MarkFlagRequired("id")
	command.Flags().StringVarP(&joinContestantId, "contestantId", "c", "", "contestant id")

	return command
}
