package cmd

import (
	"database/sql"

	tournament "bridge-tab/internal/tournament-management/domain"
	tournament_infra "bridge-tab/internal/tournament-management/infrastructure"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "bridge-tab",
	Short: "Bridge Tab CLI to manage duplicate bridge tournaments",
	Long: `Bridge Tab CLI is a tool to manage duplicate bridge tournaments. 
It allows organizers or umpires to prepare and manage tournaments, check scores, and more.`,
}

func Execute() error {
	return rootCmd.Execute()
}

// Tournament Management
var TournamentRepository tournament.TournamentRepository
var TournamentReadRepository tournament.TournamentReadRepository
var TeamReadRepository tournament.TeamReadRepository


func init() {
	cobra.OnInitialize()

	dbString := "postgres://bridge-tab:bridge-tab@localhost/bridge-tab?sslmode=disable"
	db, err := sql.Open("postgres", dbString)

	if err != nil {
		panic(err)
	}

	tournament_infra.Migrate(db)

	tournamentRepo := tournament_infra.PostgresTournamentRepository{Db: db}
	tournamentReadRepo := tournament_infra.PostgresTournamentReadRepository{Db: db}
	teamReadRepo := tournament_infra.PostgresTeamReadRepository{Db: db}
	TournamentRepository = &tournamentRepo
	TournamentReadRepository = &tournamentReadRepo
	TeamReadRepository = &teamReadRepo

	rootCmd.AddCommand(tournamentManagementCmd)
}
