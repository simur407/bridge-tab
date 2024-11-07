package cmd

import (
	"fmt"

	turnament "bridge-tab/internal/turnament-management/application"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var name string

var createTurnamentCmd = &cobra.Command{
	Use: "create",
	Short: "Creates a new turnament",
	Long: `This command creates a new turnament with given name. The name should be unique.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		command := &turnament.CreateTurnamentCommand{ Id: uuid.New().String(), Name: name }

		if err := command.Execute(turnamentRepository); err != nil {
			return err
		}

		fmt.Println("Created turnament { Id:", command.Id, "Name:", name, "}")
		return nil
	},
}

func init () {
	createTurnamentCmd.Flags().StringVarP(&name, "name", "n", "", "unique turnament name")
	createTurnamentCmd.MarkFlagRequired("name")
}
