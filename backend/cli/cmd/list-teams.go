package cmd

import (
	"fmt"
	"strings"

	tournament "bridge-tab/internal/tournament-management/application"

	"github.com/spf13/cobra"
)

var listTeamsCmd = &cobra.Command{
	Use: "list",
	Short: "Lists existing teams in a tournament",
	Long: `This command lists all existing teams in a tournaments.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		query := tournament.ListTeamsQuery{TournamentId: teamsTournamentId}

		results, err := query.Execute(TeamReadRepository)
		if err != nil {
			return err
		}

		fmt.Printf("%-36v | %-15v | %-36v\n", "Id", "Name", "Members")
		fmt.Println(strings.Repeat("-", 95))
		for _, Team := range results {
			row := fmt.Sprintf("%-36v | %-15v | ", Team.Id, Team.Name)
			for i, Member := range Team.Members {
				if i == 0 {
					row += fmt.Sprintf("%-36v\n", Member.Id)
				} else {
					row += fmt.Sprintf("%-36v   %-15v   %-36v\n", "", "", Member.Id)
				}
			}
			
			lastRune := []rune(row)[len(row)-1]
			if lastRune != '\n' {
				row += "\n"
			} 
			fmt.Print(row)
			fmt.Println(strings.Repeat("-", 95))
		}
		return nil
	},
}
