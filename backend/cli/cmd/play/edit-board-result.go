package play

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	play_app "bridge-tab/internal/play/application"
	play_query "bridge-tab/internal/play/application/query"
	play_domain "bridge-tab/internal/play/domain"
	score_domain "bridge-tab/internal/score-registration/domain"

	"github.com/spf13/cobra"
)

var editTournamentId string
var editFirstTeamNumber int
var editSecondTeamNumber int
var editDealNo int
var editContract string
var editTricks int
var editDeclarer string
var editOpeningLead string

var editBoardResultCmd = func(
	boardResultRepository *play_domain.BoardResultRepository,
	boardResultReadRepository *play_domain.BoardResultReadRepository,
	playSessionRepository *play_domain.PlaySessionRepository,
	playedResultRepository *score_domain.PlayedResultRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:          "edit",
		Short:        "Edits a board result",
		Long:         "Edits contract, tricks, declarer and opening lead of a board result. It prints the current and new score and asks for confirmation before saving. The scoring record is updated as well.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			command := &play_app.EditBoardResultCommand{
				TournamentId:     editTournamentId,
				FirstTeamNumber:  editFirstTeamNumber,
				SecondTeamNumber: editSecondTeamNumber,
				DealNo:           editDealNo,
				Contract:         editContract,
				Tricks:           editTricks,
				Declarer:         editDeclarer,
				OpeningLead:      editOpeningLead,
			}
			if err := command.Validate(); err != nil {
				return err
			}

			result, err := (&play_query.GetBoardResultByTeamsQuery{
				TournamentId:     command.TournamentId,
				FirstTeamNumber:  command.FirstTeamNumber,
				SecondTeamNumber: command.SecondTeamNumber,
				DealNo:           command.DealNo,
			}).Execute(*boardResultReadRepository, *playSessionRepository)
			if err != nil {
				return err
			}

			fmt.Printf("Deal %d (NS %d vs EW %d)\n", result.DealNo, result.NsTeamNumber, result.EwTeamNumber)
			fmt.Printf("%-14s %s -> %s\n", "Contract:", displayScore(result.Contract), displayScore(command.Contract))
			fmt.Printf("%-14s %s -> %s\n", "Declarer:", displayScore(result.Declarer), displayScore(command.Declarer))
			fmt.Printf("%-14s %d -> %d\n", "Tricks:", result.Tricks, command.Tricks)
			fmt.Printf("%-14s %s -> %s\n", "Opening lead:", displayScore(result.OpeningLead), displayScore(command.OpeningLead))
			fmt.Print("\nConfirm edit? [y/N]: ")

			answer, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
			if err != nil && err != io.EOF {
				return err
			}
			if !strings.EqualFold(strings.TrimSpace(answer), "y") && !strings.EqualFold(strings.TrimSpace(answer), "yes") {
				fmt.Println("Edit cancelled")
				return nil
			}

			if err := command.Execute(*boardResultRepository, *playSessionRepository, *playedResultRepository); err != nil {
				return err
			}

			fmt.Println("Board result updated successfully")
			return nil
		},
	}

	command.Flags().StringVarP(&editTournamentId, "id", "i", "", "tournament id")
	command.MarkFlagRequired("id")
	command.Flags().IntVarP(&editFirstTeamNumber, "firstTeamNumber", "f", 0, "1st team number")
	command.MarkFlagRequired("firstTeamNumber")
	command.Flags().IntVarP(&editSecondTeamNumber, "secondTeamNumber", "s", 0, "2nd team number")
	command.MarkFlagRequired("secondTeamNumber")
	command.Flags().IntVarP(&editDealNo, "dealNo", "n", 0, "deal no")
	command.MarkFlagRequired("dealNo")
	command.Flags().StringVarP(&editContract, "contract", "c", "", "contract")
	command.MarkFlagRequired("contract")
	command.Flags().IntVarP(&editTricks, "tricks", "r", 0, "tricks")
	command.MarkFlagRequired("tricks")
	command.Flags().StringVarP(&editDeclarer, "declarer", "d", "", "declarer")
	command.MarkFlagRequired("declarer")
	command.Flags().StringVarP(&editOpeningLead, "openingLead", "o", "", "opening lead")
	command.MarkFlagRequired("openingLead")

	return command
}

func displayScore(value string) string {
	if value == "" {
		return "-"
	}
	return value
}
