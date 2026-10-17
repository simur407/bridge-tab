package play

import (
	"fmt"

	play_app "bridge-tab/internal/play/application"
	play_domain "bridge-tab/internal/play/domain"
	score_domain "bridge-tab/internal/score-registration/domain"

	"github.com/spf13/cobra"
)

var recordPlaySessionId string
var recordPlayerId string
var recordVersusTeamNumber int
var recordDealNo int
var recordContract string
var recordTricks int
var recordDeclarer string
var recordOpeningLead string

var recordBoardResultCmd = func(
	boardResultRepository *play_domain.BoardResultRepository,
	playSessionRepository *play_domain.PlaySessionRepository,
	playedResultRepository *score_domain.PlayedResultRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:          "add",
		Short:        "Records a board result",
		Long:         "Records the result of one board at one table and registers it for scoring.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			playSessionId := recordPlaySessionId
			command := &play_app.RecordBoardResultCommand{
				PlaySessionId:    playSessionId,
				PlayerId:         recordPlayerId,
				VersusTeamNumber: recordVersusTeamNumber,
				DealNo:           recordDealNo,
				Contract:         recordContract,
				Tricks:           recordTricks,
				Declarer:         recordDeclarer,
				OpeningLead:      recordOpeningLead,
			}
			if err := command.Execute(*boardResultRepository, *playSessionRepository, *playedResultRepository); err != nil {
				return err
			}

			fmt.Println("Board result saved successfully")
			return nil
		},
	}

	command.Flags().StringVarP(&recordPlaySessionId, "id", "i", "", "play session id")
	command.MarkFlagRequired("id")
	command.Flags().StringVarP(&recordPlayerId, "playerId", "p", "", "player id")
	command.MarkFlagRequired("playerId")
	command.Flags().IntVarP(&recordVersusTeamNumber, "teamNumber", "t", 0, "versus team number")
	command.MarkFlagRequired("teamNumber")
	command.Flags().IntVarP(&recordDealNo, "dealNo", "n", 0, "deal no")
	command.MarkFlagRequired("dealNo")
	command.Flags().StringVarP(&recordContract, "contract", "c", "", "contract")
	command.MarkFlagRequired("contract")
	command.Flags().IntVarP(&recordTricks, "tricks", "r", 0, "tricks")
	command.MarkFlagRequired("tricks")
	command.Flags().StringVarP(&recordDeclarer, "declarer", "d", "", "declarer")
	command.MarkFlagRequired("declarer")
	command.Flags().StringVarP(&recordOpeningLead, "openingLead", "o", "", "opening lead")
	command.MarkFlagRequired("openingLead")

	return command
}
