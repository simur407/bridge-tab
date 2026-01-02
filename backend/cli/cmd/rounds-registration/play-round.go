package rounds_registration

import (
	rounds_registration "bridge-tab/internal/rounds-registration/application"
	rounds_registration_domain "bridge-tab/internal/rounds-registration/domain"
	tournament_management_domain "bridge-tab/internal/tournament-management/domain"
	"fmt"

	"github.com/spf13/cobra"
)

var gameSessionId string
var playerId string
var teamName string
var dealNo int
var contract string
var tricks int
var declarer string
var openingLead string

var playRoundCmd = func(GameSessionRepository *rounds_registration_domain.GameSessionRepository, TeamRepository *tournament_management_domain.TeamReadRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "add-round",
		Short:        "Saves score of a round",
		Long:         `This command saves score of a round.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			command := &rounds_registration.PlayRoundCommand{
				GameSessionId:  gameSessionId,
				PlayerId:       playerId,
				VersusTeamName: teamName,
				DealNo:         dealNo,
				Contract:       contract,
				Tricks:         tricks,
				Declarer:       declarer,
				OpeningLead:    openingLead,
			}

			if err := command.Execute(*GameSessionRepository, *TeamRepository); err != nil {
				return err
			}

			fmt.Println("Round saved successfully")
			return nil
		},
	}

	command.Flags().StringVarP(&gameSessionId, "id", "i", "", "game session id")
	command.MarkFlagRequired("id")
	command.Flags().StringVarP(&playerId, "playerId", "p", "", "player id")
	command.MarkFlagRequired("playerId")
	command.Flags().StringVarP(&teamName, "teamName", "t", "", "team name")
	command.MarkFlagRequired("teamName")
	command.Flags().IntVarP(&dealNo, "dealNo", "n", 0, "deal no")
	command.MarkFlagRequired("dealNo")
	command.Flags().StringVarP(&contract, "contract", "c", "", "contract")
	command.MarkFlagRequired("contract")
	command.Flags().IntVarP(&tricks, "tricks", "r", 0, "tricks")
	command.MarkFlagRequired("tricks")
	command.Flags().StringVarP(&declarer, "declarer", "d", "", "declarer")
	command.MarkFlagRequired("declarer")
	command.Flags().StringVarP(&openingLead, "openingLead", "o", "", "opening lead")
	command.MarkFlagRequired("openingLead")

	return command
}
