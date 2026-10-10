package rounds_registration

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"strings"

	rounds_registration "bridge-tab/internal/rounds-registration/application"
	rounds_registration_query "bridge-tab/internal/rounds-registration/application/query"
	rounds_registration_domain "bridge-tab/internal/rounds-registration/domain"
	tournament_management_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/spf13/cobra"
)

var editGameSessionId string
var editFirstTeamNumber int
var editSecondTeamNumber int
var editDealNo int
var editContract string
var editTricks int
var editDeclarer string
var editOpeningLead string

var editRoundCmd = func(
	GameSessionRepository *rounds_registration_domain.GameSessionRepository,
	GameSessionReadRepository *rounds_registration_domain.GameSessionReadRepository,
	TeamRepository *tournament_management_domain.TeamReadRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:          "edit-round",
		Short:        "Edits score of a played round",
		Long:         `This command edits contract, tricks, declarer and opening lead of a played round. It prints the current and new score and asks for confirmation before saving.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			command := &rounds_registration.EditRoundCommand{
				GameSessionId:    editGameSessionId,
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

			query := rounds_registration_query.GetRoundByTeamsQuery{
				GameSessionId:    command.GameSessionId,
				FirstTeamNumber:  command.FirstTeamNumber,
				SecondTeamNumber: command.SecondTeamNumber,
				DealNo:           command.DealNo,
			}
			round, err := query.Execute(*GameSessionReadRepository, *TeamRepository)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return rounds_registration_domain.ErrRoundNotFound
				}
				return err
			}
			if round == nil {
				return rounds_registration_domain.ErrRoundNotFound
			}
			if round.Contract == "" {
				return rounds_registration_domain.ErrRoundNotPlayed
			}

			fmt.Printf("Deal %d (NS %d vs EW %d)\n", round.DealNo, round.NsTeamNumber, round.EwTeamNumber)
			fmt.Printf("%-14s %s -> %s\n", "Contract:", displayScore(round.Contract), displayScore(command.Contract))
			fmt.Printf("%-14s %s -> %s\n", "Declarer:", displayScore(round.Declarer), displayScore(command.Declarer))
			fmt.Printf("%-14s %d -> %d\n", "Tricks:", round.Tricks, command.Tricks)
			fmt.Printf("%-14s %s -> %s\n", "Opening lead:", displayScore(round.OpeningLead), displayScore(command.OpeningLead))
			fmt.Print("\nConfirm edit? [y/N]: ")

			answer, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
			if err != nil && err != io.EOF {
				return err
			}
			if !strings.EqualFold(strings.TrimSpace(answer), "y") && !strings.EqualFold(strings.TrimSpace(answer), "yes") {
				fmt.Println("Edit cancelled")
				return nil
			}

			if err := command.Execute(*GameSessionRepository, *TeamRepository); err != nil {
				return err
			}

			fmt.Println("Round updated successfully")
			return nil
		},
	}

	command.Flags().StringVarP(&editGameSessionId, "id", "i", "", "game session id")
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
