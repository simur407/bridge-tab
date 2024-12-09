package cmd

import (
	"context"
	"database/sql"
	"time"

	tournament "bridge-tab/internal/tournament-management/domain"
	tournament_infra "bridge-tab/internal/tournament-management/infrastructure"
	user "bridge-tab/internal/user/domain"
	user_infra "bridge-tab/internal/user/infrastructure"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bridge-tab",
	Short: "Bridge Tab CLI to manage duplicate bridge tournaments",
	Long: `Bridge Tab CLI is a tool to manage duplicate bridge tournaments. 
It allows organizers or umpires to prepare and manage tournaments, check scores, and more.`,
}

// Tournament Management
var TournamentRepository tournament.TournamentRepository
var TournamentReadRepository tournament.TournamentReadRepository
var TeamReadRepository tournament.TeamReadRepository
var UserReadRepository user.UserReadRepository

func Execute() error {
	dbString := "postgres://bridge-tab:bridge-tab@localhost/bridge-tab?sslmode=disable"
	db, err := sql.Open("postgres", dbString)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic("failed to connect to database")
	}

	user_infra.Migrate(db)
	tournament_infra.Migrate(db)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	TournamentRepository = &tournament_infra.PostgresTournamentRepository{
		Ctx: ctx,
		Tx:  tx,
	}
	TournamentReadRepository = &tournament_infra.PostgresTournamentReadRepository{
		Ctx: ctx,
		Tx:  tx,
	}
	TeamReadRepository = &tournament_infra.PostgresTeamReadRepository{
		Ctx: ctx,
		Tx:  tx,
	}

	UserReadRepository = &user_infra.PostgresUserRepository{
		Ctx: ctx,
		Tx:  tx,
	}

	err = rootCmd.Execute()

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func init() {
	cobra.OnInitialize()

	rootCmd.AddCommand(tournamentManagementCmd)
	rootCmd.AddCommand(userCmd)
}
