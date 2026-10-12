package cmd

import (
	"context"
	"database/sql"
	"os"
	"time"

	play "bridge-tab/cli/cmd/play"
	score "bridge-tab/cli/cmd/score-registration"
	tournament_management "bridge-tab/cli/cmd/tournament-management"
	users "bridge-tab/cli/cmd/users"

	play_domain "bridge-tab/internal/play/domain"
	play_infra "bridge-tab/internal/play/infrastructure"
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

// Tournament Management
var TournamentRepository tournament.TournamentRepository
var TournamentReadRepository tournament.TournamentReadRepository
var TeamReadRepository tournament.TeamReadRepository
var TableReadRepository tournament.TableReadRepository
var BoardProtocolReadRepository tournament.BoardProtocolReadRepository
var SetReadRepository tournament.SetReadRepository

// Play
var BoardResultRepository play_domain.BoardResultRepository
var BoardResultReadRepository play_domain.BoardResultReadRepository
var PlaySessionRepository play_domain.PlaySessionRepository

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
	play_infra.Migrate(db)
	score_registration_infra.Migrate(db)

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

	BoardResultRepository = &play_infra.PostgresBoardResultRepository{
		Ctx: ctx,
		Tx:  tx,
	}
	BoardResultReadRepository = &play_infra.PostgresBoardResultReadRepository{
		Ctx: ctx,
		Tx:  tx,
	}
	PlaySessionRepository = &play_infra.PostgresPlaySessionRepository{
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
		&PlayedResultRepository,
		&BoardResultReadRepository,
		&PlaySessionRepository,
		&TournamentScoreRepository,
	))
	rootCmd.AddCommand(users.UserCmd(&UserReadRepository))
	rootCmd.AddCommand(play.PlayCmd(
		&BoardResultRepository,
		&BoardResultReadRepository,
		&PlaySessionRepository,
		&PlayedResultRepository,
	))
	rootCmd.AddCommand(score.ScoreRegistrationCmd(&PlayedResultRepository, &TournamentScoreRepository, &TournamentScoreReadRepository))
}
