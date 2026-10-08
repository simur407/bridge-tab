package flow

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	rounds_registration "bridge-tab/internal/rounds-registration/application"
	rounds_registration_query "bridge-tab/internal/rounds-registration/application/query"
	rounds_registration_domain "bridge-tab/internal/rounds-registration/domain"
	rounds_registration_infra "bridge-tab/internal/rounds-registration/infrastructure"
	score_app "bridge-tab/internal/score-registration/application"
	score_query "bridge-tab/internal/score-registration/application/query"
	score_registration_domain "bridge-tab/internal/score-registration/domain"
	score_registration_infra "bridge-tab/internal/score-registration/infrastructure"
	tournament_command "bridge-tab/internal/tournament-management/application/command"
	tournament_query "bridge-tab/internal/tournament-management/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
	tournament_infra "bridge-tab/internal/tournament-management/infrastructure"
	user_app "bridge-tab/internal/user/application"
	user_infra "bridge-tab/internal/user/infrastructure"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Flow struct {
	t                     *testing.T
	tournamentId          string
	teamIds               map[int]string
	tableIds              map[int]string
	tournamentRepo        tournament_domain.TournamentRepository
	tournamentReadRepo    tournament_domain.TournamentReadRepository
	teamReadRepo          tournament_domain.TeamReadRepository
	boardProtocolReadRepo tournament_domain.BoardProtocolReadRepository
	gameSessionRepo       rounds_registration_domain.GameSessionRepository
	gameSessionReadRepo   rounds_registration_domain.GameSessionReadRepository
	playedResultRepo      score_registration_domain.PlayedResultRepository
	scoreRepo             score_registration_domain.TournamentScoreRepository
	scoreReadRepo         score_registration_domain.TournamentScoreReadRepository
	userRepo              *user_infra.PostgresUserRepository
}

func New(t *testing.T) *Flow {
	t.Helper()

	dbURL := os.Getenv("TEST_DATABASE_STRING")
	if dbURL == "" {
		t.Skip("skipping integration test: set TEST_DATABASE_STRING to a Postgres URL")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := db.Ping(); err != nil {
		t.Fatalf("ping database: %v", err)
	}

	user_infra.Migrate(db)
	tournament_infra.Migrate(db)
	rounds_registration_infra.Migrate(db)
	score_registration_infra.Migrate(db)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback() })

	return &Flow{
		t:            t,
		tournamentId: uuid.New().String(),
		teamIds:      map[int]string{},
		tableIds:     map[int]string{},
		tournamentRepo: &tournament_infra.PostgresTournamentRepository{
			Ctx: ctx,
			Tx:  tx,
		},
		tournamentReadRepo: &tournament_infra.PostgresTournamentReadRepository{
			Ctx: ctx,
			Tx:  tx,
		},
		teamReadRepo: &tournament_infra.PostgresTeamReadRepository{
			Ctx: ctx,
			Tx:  tx,
		},
		boardProtocolReadRepo: &tournament_infra.PostgresBoardProtocolReadRepository{
			Ctx: ctx,
			Tx:  tx,
		},
		gameSessionRepo: &rounds_registration_infra.PostgresGameSessionRepository{
			Ctx: ctx,
			Tx:  tx,
		},
		gameSessionReadRepo: &rounds_registration_infra.PostgresGameSessionReadRepository{
			Ctx: ctx,
			Tx:  tx,
		},
		playedResultRepo: &score_registration_infra.PostgresPlayedResultRepository{
			Ctx: ctx,
			Tx:  tx,
		},
		scoreRepo: &score_registration_infra.PostgresTournamentScoreRepository{
			Ctx: ctx,
			Tx:  tx,
		},
		scoreReadRepo: &score_registration_infra.PostgresTournamentScoreReadRepository{
			Ctx: ctx,
			Tx:  tx,
		},
		userRepo: &user_infra.PostgresUserRepository{
			Ctx: ctx,
			Tx:  tx,
		},
	}
}

func (f *Flow) CreateTournament() {
	f.t.Helper()

	name := "E2E Happy Path " + f.tournamentId[:8]
	cmd := &tournament_command.CreateTournamentCommand{
		TournamentId: f.tournamentId,
		Name:         name,
	}
	if err := cmd.Execute(f.tournamentRepo); err != nil {
		f.t.Fatalf("create tournament: %v", err)
	}

	tournament, err := (&tournament_query.GetTournamentById{Id: f.tournamentId}).Execute(f.tournamentReadRepo)
	if err != nil {
		f.t.Fatalf("get tournament: %v", err)
	}
	if tournament == nil || tournament.Name != name {
		f.t.Fatalf("expected tournament %q, got %#v", name, tournament)
	}
}

func (f *Flow) CreateTeam(number int) {
	f.t.Helper()

	teamId := uuid.New().String()
	cmd := &tournament_command.CreateTeamCommand{
		TournamentId: f.tournamentId,
		TeamId:       teamId,
		Number:       number,
	}
	if err := cmd.Execute(f.tournamentRepo); err != nil {
		f.t.Fatalf("create team %d: %v", number, err)
	}
	f.teamIds[number] = teamId
}

func (f *Flow) CreateTable(number int) {
	f.t.Helper()

	tableId := uuid.New().String()
	cmd := &tournament_command.CreateTableCommand{
		TournamentId: f.tournamentId,
		TableId:      tableId,
		Number:       number,
	}
	if err := cmd.Execute(f.tournamentRepo); err != nil {
		f.t.Fatalf("create table %d: %v", number, err)
	}
	f.tableIds[number] = tableId
}

func (f *Flow) CreateBoardProtocol(boardNo int, vulnerable tournament_domain.Vulnerable, nsTeamNumber, ewTeamNumber int) {
	f.t.Helper()
	f.createBoardProtocolWithOptionalTable(boardNo, vulnerable, 0, nsTeamNumber, ewTeamNumber)
}

func (f *Flow) CreateBoardProtocolWithPairs(boardNo int, vulnerable tournament_domain.Vulnerable, pairs [][2]int) {
	f.t.Helper()

	teamPairs := make([]struct {
		Table *string
		NS    string
		EW    string
	}, 0, len(pairs))
	for _, pair := range pairs {
		nsTeamId, ok := f.teamIds[pair[0]]
		if !ok {
			f.t.Fatalf("unknown NS team number %d", pair[0])
		}
		ewTeamId, ok := f.teamIds[pair[1]]
		if !ok {
			f.t.Fatalf("unknown EW team number %d", pair[1])
		}
		teamPairs = append(teamPairs, struct {
			Table *string
			NS    string
			EW    string
		}{
			NS: nsTeamId,
			EW: ewTeamId,
		})
	}

	cmd := &tournament_command.CreateBoardProtocol{
		TournamentId: f.tournamentId,
		BoardNo:      boardNo,
		Vulnerable:   int(vulnerable),
		TeamPairs:    teamPairs,
	}
	if err := cmd.Execute(f.tournamentRepo); err != nil {
		f.t.Fatalf("create board protocol %d: %v", boardNo, err)
	}
}

func (f *Flow) CreateSet(label string, boardNos []int) {
	f.t.Helper()

	cmd := &tournament_command.CreateSetCommand{
		TournamentId: f.tournamentId,
		SetId:        uuid.New().String(),
		Label:        label,
		BoardNos:     boardNos,
	}
	if err := cmd.Execute(f.tournamentRepo); err != nil {
		f.t.Fatalf("create set %s: %v", label, err)
	}
}

func (f *Flow) CreateBoardProtocolWithTable(boardNo int, vulnerable tournament_domain.Vulnerable, tableNumber, nsTeamNumber, ewTeamNumber int) {
	f.t.Helper()
	f.createBoardProtocolWithOptionalTable(boardNo, vulnerable, tableNumber, nsTeamNumber, ewTeamNumber)
}

func (f *Flow) createBoardProtocolWithOptionalTable(boardNo int, vulnerable tournament_domain.Vulnerable, tableNumber, nsTeamNumber, ewTeamNumber int) {
	f.t.Helper()

	nsTeamId, ok := f.teamIds[nsTeamNumber]
	if !ok {
		f.t.Fatalf("unknown NS team number %d", nsTeamNumber)
	}
	ewTeamId, ok := f.teamIds[ewTeamNumber]
	if !ok {
		f.t.Fatalf("unknown EW team number %d", ewTeamNumber)
	}

	pair := struct {
		Table *string
		NS    string
		EW    string
	}{
		NS: nsTeamId,
		EW: ewTeamId,
	}
	if tableNumber != 0 {
		tableId, ok := f.tableIds[tableNumber]
		if !ok {
			f.t.Fatalf("unknown table number %d", tableNumber)
		}
		pair.Table = &tableId
	}

	cmd := &tournament_command.CreateBoardProtocol{
		TournamentId: f.tournamentId,
		BoardNo:      boardNo,
		Vulnerable:   int(vulnerable),
		TeamPairs: []struct {
			Table *string
			NS    string
			EW    string
		}{pair},
	}
	if err := cmd.Execute(f.tournamentRepo); err != nil {
		f.t.Fatalf("create board protocol %d: %v", boardNo, err)
	}
}

func (f *Flow) RegisterUser(name string) string {
	f.t.Helper()

	userId := uuid.New().String()
	if err := (&user_app.RegisterUserCommand{Id: userId, Name: name}).Execute(f.userRepo); err != nil {
		f.t.Fatalf("register user %s: %v", name, err)
	}
	return userId
}

func (f *Flow) JoinTournament(userId string) {
	f.t.Helper()

	cmd := &tournament_command.JoinTournamentCommand{
		TournamentId: f.tournamentId,
		ContestantId: userId,
	}
	if err := cmd.Execute(f.tournamentRepo); err != nil {
		f.t.Fatalf("join tournament: %v", err)
	}
}

func (f *Flow) JoinTeam(userId string, teamNumber int) {
	f.t.Helper()

	teamId, ok := f.teamIds[teamNumber]
	if !ok {
		f.t.Fatalf("unknown team number %d", teamNumber)
	}

	cmd := &tournament_command.JoinTeamCommand{
		TournamentId: f.tournamentId,
		TeamId:       teamId,
		ContestantId: userId,
	}
	if err := cmd.Execute(f.tournamentRepo); err != nil {
		f.t.Fatalf("join team %d: %v", teamNumber, err)
	}
}

func (f *Flow) StartTournament() {
	f.t.Helper()

	cmd := &tournament_command.StartTurnamentCommand{TournamentId: f.tournamentId}
	if err := cmd.Execute(f.tournamentRepo, f.teamReadRepo, f.boardProtocolReadRepo, f.gameSessionRepo); err != nil {
		f.t.Fatalf("start tournament: %v", err)
	}

	started, err := (&tournament_query.GetTournamentById{Id: f.tournamentId}).Execute(f.tournamentReadRepo)
	if err != nil {
		f.t.Fatalf("get started tournament: %v", err)
	}
	if started == nil || started.StartedAt == "" {
		f.t.Fatal("expected tournament to be started")
	}
}

func (f *Flow) AssertPendingRounds(expected int) {
	f.t.Helper()

	rounds := f.ListRounds()
	if len(rounds) != expected {
		f.t.Fatalf("expected %d pending rounds after start, got %d", expected, len(rounds))
	}
	for _, round := range rounds {
		if round.Contract != "" {
			f.t.Fatalf("expected empty contract before play, got %#v", round)
		}
	}
}

func (f *Flow) PlayRound(playerId string, versusTeamNumber, dealNo int, contract string, tricks int, declarer, openingLead string) {
	f.t.Helper()

	cmd := &rounds_registration.PlayRoundCommand{
		GameSessionId:    f.tournamentId,
		PlayerId:         playerId,
		VersusTeamNumber: versusTeamNumber,
		DealNo:           dealNo,
		Contract:         contract,
		Tricks:           tricks,
		Declarer:         declarer,
		OpeningLead:      openingLead,
	}
	if err := cmd.Execute(f.gameSessionRepo, f.teamReadRepo); err != nil {
		f.t.Fatalf("play deal %d: %v", dealNo, err)
	}
}

func (f *Flow) ListRounds() []rounds_registration_domain.PlayedRoundDto {
	f.t.Helper()

	rounds, err := (&rounds_registration_query.ListRoundsQuery{GameSessionId: f.tournamentId}).Execute(f.gameSessionReadRepo)
	if err != nil {
		f.t.Fatalf("list rounds: %v", err)
	}
	return rounds
}

func (f *Flow) FinishTournament() {
	f.t.Helper()

	cmd := &tournament_command.FinishTournamentCommand{TournamentId: f.tournamentId}
	if err := cmd.Execute(f.tournamentRepo, f.gameSessionRepo, f.boardProtocolReadRepo, f.playedResultRepo); err != nil {
		f.t.Fatalf("finish tournament: %v", err)
	}

	finished, err := (&tournament_query.GetTournamentById{Id: f.tournamentId}).Execute(f.tournamentReadRepo)
	if err != nil {
		f.t.Fatalf("get finished tournament: %v", err)
	}
	if finished == nil || finished.FinishedAt == "" {
		f.t.Fatal("expected tournament to be finished")
	}
}

func (f *Flow) CalculateStandings(method score_registration_domain.ScoringMethodName) {
	f.t.Helper()

	cmd := &score_app.CalculateStandingsCommand{
		TournamentId: f.tournamentId,
		Method:       method,
	}
	if err := cmd.Execute(f.playedResultRepo, f.scoreRepo); err != nil {
		f.t.Fatalf("calculate standings: %v", err)
	}
}

func (f *Flow) GetStandings(method score_registration_domain.ScoringMethodName) []score_registration_domain.StandingDto {
	f.t.Helper()

	standings, err := (&score_query.GetStandingsQuery{
		TournamentId: f.tournamentId,
		Method:       method,
	}).Execute(f.scoreReadRepo)
	if err != nil {
		f.t.Fatalf("get standings: %v", err)
	}
	return standings
}

func (f *Flow) GetBoardScores(method score_registration_domain.ScoringMethodName, dealNo int) []score_registration_domain.BoardScoreDto {
	f.t.Helper()

	scores, err := (&score_query.ListBoardScoresQuery{
		TournamentId: f.tournamentId,
		Method:       method,
		DealNo:       &dealNo,
	}).Execute(f.scoreReadRepo)
	if err != nil {
		f.t.Fatalf("get board scores: %v", err)
	}
	return scores
}

func (f *Flow) AllBoardScores(method score_registration_domain.ScoringMethodName) []score_registration_domain.BoardScoreDto {
	f.t.Helper()

	scores, err := (&score_query.ListBoardScoresQuery{
		TournamentId: f.tournamentId,
		Method:       method,
	}).Execute(f.scoreReadRepo)
	if err != nil {
		f.t.Fatalf("get board scores: %v", err)
	}
	return scores
}
