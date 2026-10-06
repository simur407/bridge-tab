package e2e_test

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
	tournament_command "bridge-tab/internal/tournament-management/application/command"
	tournament_query "bridge-tab/internal/tournament-management/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
	tournament_infra "bridge-tab/internal/tournament-management/infrastructure"
	user_app "bridge-tab/internal/user/application"
	user_infra "bridge-tab/internal/user/infrastructure"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func TestTournamentHappyPath(t *testing.T) {
	flow := newTournamentFlow(t)

	flow.createTournament()
	flow.createTeam(1)
	flow.createTeam(2)

	flow.createBoardProtocol(1, tournament_domain.None, 1, 2)
	flow.createBoardProtocol(2, tournament_domain.NS, 1, 2)
	flow.createBoardProtocol(3, tournament_domain.EW, 1, 2)
	flow.createBoardProtocol(4, tournament_domain.Both, 1, 2)

	alice := flow.registerUser("Alice")
	flow.joinTournament(alice)
	flow.joinTeam(alice, 1)

	bob := flow.registerUser("Bob")
	flow.joinTournament(bob)
	flow.joinTeam(bob, 2)

	flow.startTournament()
	flow.assertPendingRounds(4)

	flow.playRound(alice, 2, 1, "3NT", 9, "N", "KH")
	flow.playRound(bob, 1, 2, "4H", 10, "E", "AS")
	flow.playRound(alice, 2, 3, "1Cx", 7, "S", "QD")
	flow.playRound(bob, 1, 4, "Pass", 0, "", "")

	rounds := flow.listRounds()
	assertDeal(t, rounds, 1, 1, 2, "3N", 9, "N", "KH")
	assertDeal(t, rounds, 2, 1, 2, "4H", 10, "E", "AS")
	assertDeal(t, rounds, 3, 1, 2, "1Cx", 7, "S", "QD")
	assertDeal(t, rounds, 4, 1, 2, "Pass", 0, "", "")
}

type tournamentFlow struct {
	t                     *testing.T
	tournamentId          string
	teamIds               map[int]string
	tournamentRepo        tournament_domain.TournamentRepository
	tournamentReadRepo    tournament_domain.TournamentReadRepository
	teamReadRepo          tournament_domain.TeamReadRepository
	boardProtocolReadRepo tournament_domain.BoardProtocolReadRepository
	gameSessionRepo       rounds_registration_domain.GameSessionRepository
	gameSessionReadRepo   rounds_registration_domain.GameSessionReadRepository
	userRepo              *user_infra.PostgresUserRepository
}

func newTournamentFlow(t *testing.T) *tournamentFlow {
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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback() })

	return &tournamentFlow{
		t:            t,
		tournamentId: uuid.New().String(),
		teamIds:      map[int]string{},
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
		userRepo: &user_infra.PostgresUserRepository{
			Ctx: ctx,
			Tx:  tx,
		},
	}
}

func (f *tournamentFlow) createTournament() {
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

func (f *tournamentFlow) createTeam(number int) {
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

func (f *tournamentFlow) createBoardProtocol(boardNo int, vulnerable tournament_domain.Vulnerable, nsTeamNumber, ewTeamNumber int) {
	f.t.Helper()

	nsTeamId, ok := f.teamIds[nsTeamNumber]
	if !ok {
		f.t.Fatalf("unknown NS team number %d", nsTeamNumber)
	}
	ewTeamId, ok := f.teamIds[ewTeamNumber]
	if !ok {
		f.t.Fatalf("unknown EW team number %d", ewTeamNumber)
	}

	cmd := &tournament_command.CreateBoardProtocol{
		TournamentId: f.tournamentId,
		BoardNo:      boardNo,
		Vulnerable:   int(vulnerable),
		TeamPairs: []struct {
			NS string
			EW string
		}{
			{NS: nsTeamId, EW: ewTeamId},
		},
	}
	if err := cmd.Execute(f.tournamentRepo); err != nil {
		f.t.Fatalf("create board protocol %d: %v", boardNo, err)
	}
}

func (f *tournamentFlow) registerUser(name string) string {
	f.t.Helper()

	userId := uuid.New().String()
	if err := (&user_app.RegisterUserCommand{Id: userId, Name: name}).Execute(f.userRepo); err != nil {
		f.t.Fatalf("register user %s: %v", name, err)
	}
	return userId
}

func (f *tournamentFlow) joinTournament(userId string) {
	f.t.Helper()

	cmd := &tournament_command.JoinTournamentCommand{
		TournamentId: f.tournamentId,
		ContestantId: userId,
	}
	if err := cmd.Execute(f.tournamentRepo); err != nil {
		f.t.Fatalf("join tournament: %v", err)
	}
}

func (f *tournamentFlow) joinTeam(userId string, teamNumber int) {
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

func (f *tournamentFlow) startTournament() {
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

func (f *tournamentFlow) assertPendingRounds(expected int) {
	f.t.Helper()

	rounds := f.listRounds()
	if len(rounds) != expected {
		f.t.Fatalf("expected %d pending rounds after start, got %d", expected, len(rounds))
	}
	for _, round := range rounds {
		if round.Contract != "" {
			f.t.Fatalf("expected empty contract before play, got %#v", round)
		}
	}
}

func (f *tournamentFlow) playRound(playerId string, versusTeamNumber, dealNo int, contract string, tricks int, declarer, openingLead string) {
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

func (f *tournamentFlow) listRounds() []rounds_registration_domain.PlayedRoundDto {
	f.t.Helper()

	rounds, err := (&rounds_registration_query.ListRoundsQuery{GameSessionId: f.tournamentId}).Execute(f.gameSessionReadRepo)
	if err != nil {
		f.t.Fatalf("list rounds: %v", err)
	}
	return rounds
}

func assertDeal(
	t *testing.T,
	rounds []rounds_registration_domain.PlayedRoundDto,
	dealNo, nsTeamNumber, ewTeamNumber int,
	contract string,
	tricks int,
	declarer, openingLead string,
) {
	t.Helper()

	var deal *rounds_registration_domain.PlayedRoundDto
	for i := range rounds {
		if rounds[i].DealNo == dealNo {
			deal = &rounds[i]
			break
		}
	}
	if deal == nil {
		t.Fatalf("missing deal %d in final rounds", dealNo)
	}
	if deal.NsTeamNumber != nsTeamNumber || deal.EwTeamNumber != ewTeamNumber {
		t.Fatalf("deal %d teams: got NS=%d EW=%d, want NS=%d EW=%d",
			dealNo, deal.NsTeamNumber, deal.EwTeamNumber, nsTeamNumber, ewTeamNumber)
	}
	if deal.Contract != contract || deal.Tricks != tricks || deal.Declarer != declarer || deal.OpeningLead != openingLead {
		t.Fatalf("deal %d score mismatch: %#v", dealNo, deal)
	}
}
