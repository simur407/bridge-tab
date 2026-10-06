package teams

import (
	tournament_management "bridge-tab/internal/tournament-management/application/command"
	tournament_query "bridge-tab/internal/tournament-management/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var teamName string
var teamNumber int

var createTeamCmd = func(
	TournamentRepository *tournament_domain.TournamentRepository,
	TeamReadRepository *tournament_domain.TeamReadRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:          "create",
		Short:        "Creates team in a tournament",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			command := &tournament_management.CreateTeamCommand{
				TournamentId: teamsTournamentId,
				TeamId:       uuid.New().String(),
				Name:         teamName,
				Number:       teamNumber,
			}

			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			team, err := (&tournament_query.GetTeamByIdQuery{
				TournamentId: command.TournamentId,
				TeamId:       command.TeamId,
			}).Execute(*TeamReadRepository)
			if err != nil {
				return err
			}

			fmt.Println("Team { Id:", team.Id, ", Name:", team.Name, ", Number:", team.Number, " } created in Tournament { Id:", command.TournamentId, "}")
			return nil
		},
	}

	command.Flags().StringVarP(&teamName, "name", "n", "", "optional Team name")
	command.Flags().IntVar(&teamNumber, "number", 0, "team number; omitted values continue from 1")

	return command
}
