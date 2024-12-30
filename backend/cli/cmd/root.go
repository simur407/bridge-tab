package cmd

import (
	"context"
	"database/sql"
	"time"

	rounds "bridge-tab/cli/cmd/rounds-registration"
	tournament_management "bridge-tab/cli/cmd/tournament-management"
	users "bridge-tab/cli/cmd/users"

	rounds_registration "bridge-tab/internal/rounds-registration/domain"
	rounds_registration_infra "bridge-tab/internal/rounds-registration/infrastructure"
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

// Round Registration
var GameSessionRepository rounds_registration.GameSessionRepository

// Tournament Management
var TournamentRepository tournament.TournamentRepository
var TournamentReadRepository tournament.TournamentReadRepository
var TeamReadRepository tournament.TeamReadRepository
var BoardProtocolReadRepository tournament.BoardProtocolReadRepository

// Users
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
	rounds_registration_infra.Migrate(db)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Round Registration
	GameSessionRepository = &rounds_registration_infra.PostgresGameSessionRepository{
		Ctx: ctx,
		Tx:  tx,
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
	BoardProtocolReadRepository = &tournament_infra.PostgresBoardProtocolReadRepository{
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

	rootCmd.AddCommand(tournament_management.TournamentManagementCmd(
		&TournamentRepository,
		&TournamentReadRepository,
		&TeamReadRepository,
		&BoardProtocolReadRepository,
		&GameSessionRepository,
	))
	rootCmd.AddCommand(users.UserCmd(&UserReadRepository))
	rootCmd.AddCommand(rounds.RoundsRegistrationCmd(&GameSessionRepository, &TeamReadRepository))
}
