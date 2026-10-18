package tournament_management

import (
	"errors"
	"fmt"

	tournament_management "bridge-tab/internal/tournament-management/application/command"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/spf13/cobra"
)

var requireTeamMembersTournamentId string
var requireTeamMembersOn bool
var requireTeamMembersOff bool

var requireTeamMembersCmd = func(TournamentRepository *tournament_domain.TournamentRepository) *cobra.Command {
	command := &cobra.Command{
		Use:          "require-team-members",
		Short:        "Turns the start-time team member requirement on or off",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if requireTeamMembersOn == requireTeamMembersOff {
				return errors.New("specify exactly one of --on or --off")
			}
			command := &tournament_management.SetRequireTeamMembersCommand{
				TournamentId:       requireTeamMembersTournamentId,
				RequireTeamMembers: requireTeamMembersOn,
			}
			if err := command.Execute(*TournamentRepository); err != nil {
				return err
			}
			state := "off"
			if requireTeamMembersOn {
				state = "on"
			}
			fmt.Println("Tournament { Id:", command.TournamentId, "} require-team-members:", state)
			return nil
		},
	}

	command.Flags().StringVarP(&requireTeamMembersTournamentId, "id", "i", "", "tournament id")
	command.MarkFlagRequired("id")
	command.Flags().BoolVar(&requireTeamMembersOn, "on", false, "require at least one contestant in each team before start")
	command.Flags().BoolVar(&requireTeamMembersOff, "off", false, "allow starting with empty teams")

	return command
}
