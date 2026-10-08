package cmd

import (
	"context"
	"database/sql"
	"os"
	"time"

	rounds "bridge-tab/cli/cmd/rounds-registration"
	score "bridge-tab/cli/cmd/score-registration"
	tournament_management "bridge-tab/cli/cmd/tournament-management"
	users "bridge-tab/cli/cmd/users"

	rounds_registration "bridge-tab/internal/rounds-registration/domain"
	rounds_registration_infra "bridge-tab/internal/rounds-registration/infrastructure"
	score_registration "bridge-tab/internal/score-registration/domain"
	score_registration_infra "bridge-tab/internal/score-registration/infrastructure"
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
var GameSessionReadRepository rounds_registration.GameSessionReadRepository

// Tournament Management
var TournamentRepository tournament.TournamentRepository
var TournamentReadRepository tournament.TournamentReadRepository
var TeamReadRepository tournament.TeamReadRepository
var TableReadRepository tournament.TableReadRepository
var BoardProtocolReadRepository tournament.BoardProtocolReadRepository
var SetReadRepository tournament.SetReadRepository

// Score Registration
var PlayedResultRepository score_registration.PlayedResultRepository
var TournamentScoreRepository score_registration.TournamentScoreRepository
var TournamentScoreReadRepository score_registration.TournamentScoreReadRepository

// Users
var UserReadRepository user.UserReadRepository

func Execute() error {
	dbString := os.Getenv("DATABASE_STRING")
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
	score_registration_infra.Migrate(db)

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
	GameSessionReadRepository = &rounds_registration_infra.PostgresGameSessionReadRepository{
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
	TableReadRepository = &tournament_infra.PostgresTableReadRepository{
		Ctx: ctx,
		Tx:  tx,
	}
	BoardProtocolReadRepository = &tournament_infra.PostgresBoardProtocolReadRepository{
		Ctx: ctx,
		Tx:  tx,
	}
	SetReadRepository = &tournament_infra.PostgresSetReadRepository{
		Ctx: ctx,
		Tx:  tx,
	}

	PlayedResultRepository = &score_registration_infra.PostgresPlayedResultRepository{
		Ctx: ctx,
		Tx:  tx,
	}
	TournamentScoreRepository = &score_registration_infra.PostgresTournamentScoreRepository{
		Ctx: ctx,
		Tx:  tx,
	}
	TournamentScoreReadRepository = &score_registration_infra.PostgresTournamentScoreReadRepository{
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
		&TableReadRepository,
		&BoardProtocolReadRepository,
		&SetReadRepository,
		&GameSessionRepository,
		&PlayedResultRepository,
	))
	rootCmd.AddCommand(users.UserCmd(&UserReadRepository))
	rootCmd.AddCommand(rounds.RoundsRegistrationCmd(&GameSessionRepository, &GameSessionReadRepository, &TeamReadRepository))
	rootCmd.AddCommand(score.ScoreRegistrationCmd(&PlayedResultRepository, &TournamentScoreRepository, &TournamentScoreReadRepository))
}
