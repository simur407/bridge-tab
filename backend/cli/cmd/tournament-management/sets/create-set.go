package sets

import (
	application "bridge-tab/internal/tournament-management/application/command"
	domain "bridge-tab/internal/tournament-management/domain"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var setLabel string
var setBoardNos string

var createSetCmd = func(TournamentRepository *domain.TournamentRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "create",
		Short:        "Creates a set grouping existing board protocols",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if setLabel == "" {
				return fmt.Errorf("label is required")
			}
			if setBoardNos == "" {
				return fmt.Errorf("boards numbers are required")
			}

			parts := strings.Split(setBoardNos, ",")
			boardNos := make([]int, 0, len(parts))
			for _, part := range parts {
				part = strings.TrimSpace(part)
				boardNo, err := strconv.Atoi(part)
				if err != nil {
					return fmt.Errorf("invalid board number %q", part)
				}
				boardNos = append(boardNos, boardNo)
			}

			setId := uuid.New().String()
			command := &application.CreateSetCommand{
				TournamentId: setsTournamentId,
				SetId:        setId,
				Label:        setLabel,
				BoardNos:     boardNos,
			}
			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}

			fmt.Println("Set { Id:", setId, ", Label:", setLabel, ", Boards:", boardNos, "} created in Tournament { Id:", setsTournamentId, "}")
			return nil
		},
	}

	command.Flags().StringVarP(&setLabel, "label", "l", "", "set label (e.g. A)")
	command.MarkFlagRequired("label")
	command.Flags().StringVarP(&setBoardNos, "boards", "b", "", "comma-separated board numbers that share team pairings")
	command.MarkFlagRequired("boards")

	return command
}
