package board_protocols

import (
	domain "bridge-tab/internal/tournament-management/domain"

	"github.com/spf13/cobra"
)

var boardProtocolsTournamentId string

var BoardProtocolsCmd = func(
	TournamentRepository *domain.TournamentRepository,
	TeamReadRepository *domain.TeamReadRepository,
	BoardProtocolReadRepository *domain.BoardProtocolReadRepository,
) *cobra.Command {
	command := &cobra.Command{
		Use:   "board-protocol",
		Short: "Resposible for managing Board Protocols in Tournaments",
		Long:  "Board Protocol Management allows organizers or umpires to manage Board Protocols in Tournaments like: create, delete, etc.",
	}

	command.PersistentFlags().StringVarP(&boardProtocolsTournamentId, "tournamentId", "i", "", "tournament id")
	command.MarkPersistentFlagRequired("tournamentId")
	command.AddCommand(
		createBoardProtocolCmd(TournamentRepository, TeamReadRepository),
		listBoardProtocolsCmd(BoardProtocolReadRepository),
	)

	return command
}
