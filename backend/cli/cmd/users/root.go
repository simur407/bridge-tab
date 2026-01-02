package users

import (
	user_domain "bridge-tab/internal/user/domain"

	"github.com/spf13/cobra"
)

var UserCmd = func(userReadRepository *user_domain.UserReadRepository) *cobra.Command {
	command := &cobra.Command{
		Use:   "user",
		Short: "Responsible for managing Users",
		Long:  "User allows admins to manage Users like: list",
	}

	command.AddCommand(
		listUsersCmd(userReadRepository),
	)

	return command
}
