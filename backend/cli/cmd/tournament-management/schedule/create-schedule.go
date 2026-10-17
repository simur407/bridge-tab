package schedule

import (
	"fmt"
	"strings"

	application "bridge-tab/internal/tournament-management/application/command"
	application_query "bridge-tab/internal/tournament-management/application/query"
	domain "bridge-tab/internal/tournament-management/domain"

	"github.com/spf13/cobra"
)

var roundNo int
var tableNumber int
var nsNumber int
var ewNumber int
var setLabel string

var createScheduleCmd = func(
	TournamentRepository *domain.TournamentRepository,
	TeamReadRepository *domain.TeamReadRepository,
	TableReadRepository *domain.TableReadRepository,
	SetReadRepository *domain.SetReadRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:          "create",
		Short:        "Schedules a set at a table for a round",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if roundNo < 1 {
				return fmt.Errorf("round number is required")
			}

			table, err := (&application_query.GetTableByNumberQuery{
				TournamentId: scheduleTournamentId,
				Number:       tableNumber,
			}).Execute(*TableReadRepository)
			if err != nil {
				return fmt.Errorf("table %d: %w", tableNumber, err)
			}

			nsTeam, err := (&application_query.GetTeamByNumberQuery{
				TournamentId: scheduleTournamentId,
				Number:       nsNumber,
			}).Execute(*TeamReadRepository)
			if err != nil {
				return fmt.Errorf("NS team %d: %w", nsNumber, err)
			}

			ewTeam, err := (&application_query.GetTeamByNumberQuery{
				TournamentId: scheduleTournamentId,
				Number:       ewNumber,
			}).Execute(*TeamReadRepository)
			if err != nil {
				return fmt.Errorf("EW team %d: %w", ewNumber, err)
			}

			sets, err := (&application_query.ListSetsQuery{TournamentId: scheduleTournamentId}).Execute(*SetReadRepository)
			if err != nil {
				return err
			}
			var setId string
			for _, set := range sets {
				if strings.EqualFold(set.Label, setLabel) {
					setId = set.Id
					break
				}
			}
			if setId == "" {
				return fmt.Errorf("set %q not found", setLabel)
			}

			command := &application.AddScheduleCommand{
				TournamentId: scheduleTournamentId,
				RoundNo:      roundNo,
				TableId:      table.Id,
				NsTeamId:     nsTeam.Id,
				EwTeamId:     ewTeam.Id,
				SetId:        setId,
			}
			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			fmt.Println("Schedule { Round:", roundNo, ", Table:", tableNumber, ", NS:", nsNumber, ", EW:", ewNumber, ", Set:", setLabel, "} created in Tournament { Id:", scheduleTournamentId, "}")
			return nil
		},
	}

	command.Flags().IntVarP(&roundNo, "round", "r", 0, "round number")
	command.MarkFlagRequired("round")
	command.Flags().IntVarP(&tableNumber, "table", "t", 0, "table number")
	command.MarkFlagRequired("table")
	command.Flags().IntVarP(&nsNumber, "ns", "n", 0, "NS team number")
	command.MarkFlagRequired("ns")
	command.Flags().IntVarP(&ewNumber, "ew", "e", 0, "EW team number")
	command.MarkFlagRequired("ew")
	command.Flags().StringVarP(&setLabel, "set", "s", "", "set label")
	command.MarkFlagRequired("set")

	return command
}
