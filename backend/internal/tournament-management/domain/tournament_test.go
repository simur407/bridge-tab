package tournament_management_test

import (
	"fmt"
	"testing"

	. "bridge-tab/internal/tournament-management/domain"
)

const id TournamentId = "id"

func TestCreateTournament(t *testing.T) {
	// when
	Tournament := CreateTournament(id, "name")
	// then
	assertEvent(t, Tournament.GetEvents(), TournamentCreated{TournamentId: id, Name: "name"})
}

// ------ Remove ------
func TestRemoveTournament(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	Tournament.Commit() // flush events
	// when
	err := Tournament.Remove()
	// then
	assertNoError(t, err)
	assertEvent(t, Tournament.GetEvents(), TournamentRemoved{TournamentId: id})
}

func TestRemoveRemovedTournament(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	Tournament.Remove()
	Tournament.Commit() // flush events
	// when
	err := Tournament.Remove()
	// then
	assertNoError(t, err)
	assertNoEvents(t, Tournament.GetEvents())
}

func TestRemoveStartedTournament(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	Tournament.Start()
	Tournament.Commit() // flush events
	// when
	err := Tournament.Remove()
	// then
	assertError(t, err, ErrTournamentAlreadyStarted)
}

func TestRemoveTournamentWithTeamsAndContestants(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var contestantId ContestantId = "id"
	Tournament.JoinTournament(&contestantId)
	var teamId TeamId = "id"
	Tournament.CreateTeam(&teamId, "name")
	Tournament.JoinTeam(&teamId, &contestantId)
	Tournament.Commit() // flush events
	// when
	err := Tournament.Remove()
	// then
	assertNoError(t, err)
	assertEvents(t, Tournament.GetEvents(), []any{
		ContestantLeftTeam{TeamId: teamId, ContestantId: contestantId},
		TeamRemoved{TournamentId: id, TeamId: teamId},
		ContestantLeftTournament{TournamentId: id, ContestantId: contestantId},
		TournamentRemoved{TournamentId: id},
	})
}


// ------ Start ------
func TestStartTournament(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "id"
	Tournament.CreateTeam(&teamId, "name")
	var contestantId ContestantId = "id"
	Tournament.JoinTournament(&contestantId)
	Tournament.JoinTeam(&teamId, &contestantId)
	Tournament.Commit() // flush events
	// when
	err := Tournament.Start()
	// then
	assertNoError(t, err)
	assertEvent(t, Tournament.GetEvents(), TournamentStarted{TournamentId: id, StartedAt: *Tournament.State.StartedAt})
}

func TestStartStartedTournament(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	Tournament.Start()
	Tournament.Commit() // flush events
	// when
	err := Tournament.Start()
	// then
	assertNoError(t, err)
	assertNoEvents(t, Tournament.GetEvents())
}

func TestStartRemovedTournament(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	Tournament.Remove()
	Tournament.Commit() // flush events
	// when
	err := Tournament.Start()
	// then
	assertError(t, err, ErrTournamentRemoved)
}

func TestStartTournamentWithTeamsThatHasNoMembers(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "id"
	Tournament.CreateTeam(&teamId, "name")
	Tournament.Commit() // flush events
	// when
	err := Tournament.Start()
	// then
	assertError(t, err, ErrSomeTeamHasNoMembers)
}

// ------ Join Tournament ------
func TestJoinTournament(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var contestantId ContestantId = "id"
	Tournament.Commit() // flush events
	// when
	err := Tournament.JoinTournament(&contestantId)
	// then
	assertNoError(t, err)
	assertEvent(t, Tournament.GetEvents(), ContestantJoinedTournament{TournamentId: id, ContestantId: contestantId})
}

func TestJoinAlreadyJoinedTournament(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var contestantId ContestantId = "id"
	Tournament.JoinTournament(&contestantId)
	Tournament.Commit() // flush events
	// when
	err := Tournament.JoinTournament(&contestantId)
	// then
	assertNoError(t, err)
	assertNoEvents(t, Tournament.GetEvents())
}

func TestJoinRemovedTournament(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	Tournament.Remove()
	Tournament.Commit() // flush events
	// when
	err := Tournament.JoinTournament(nil)
	// then
	assertError(t, err, ErrTournamentRemoved)
}

// ------ Leave Tournament ------
func TestLeaveTournament(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var contestantId ContestantId = "id"
	Tournament.JoinTournament(&contestantId)
	Tournament.Commit() // flush events
	// when
	err := Tournament.LeaveTournament(&contestantId)
	// then
	assertNoError(t, err)
	assertEvent(t, Tournament.GetEvents(), ContestantLeftTournament{TournamentId: Tournament.State.Id, ContestantId: contestantId})
}

func TestLeaveNotJoinedTournament(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var contestantId ContestantId = "id"
	Tournament.Commit() // flush events
	// when
	err := Tournament.LeaveTournament(&contestantId)
	// then
	assertNoError(t, err)
	assertNoEvents(t, Tournament.GetEvents())
}

func TestLeaveRemovedTournament(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var contestantId ContestantId = "id"
	Tournament.JoinTournament(&contestantId)
	Tournament.Remove()
	Tournament.Commit() // flush events
	// when
	err := Tournament.LeaveTournament(&contestantId)
	// then
	assertError(t, err, ErrTournamentRemoved)
}

// ------ Join Team ------
func TestJoinTeam(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "id"
	Tournament.CreateTeam(&teamId, "name")
	var contestantId ContestantId = "id"
	Tournament.JoinTournament(&contestantId)
	Tournament.Commit() // flush events
	// when
	err := Tournament.JoinTeam(&teamId, &contestantId)
	// then
	assertNoError(t, err)
	assertEvent(t, Tournament.GetEvents(), ContestantJoinedTeam{ContestantId: contestantId, TeamId: teamId})
}

func TestJoinAlreadyJoinedTeam(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "id"
	Tournament.CreateTeam(&teamId, "name")
	var contestantId ContestantId = "id"
	Tournament.JoinTournament(&contestantId)
	Tournament.JoinTeam(&teamId, &contestantId)
	Tournament.Commit() // flush events
	// when
	err := Tournament.JoinTeam(&teamId, &contestantId)
	// then
	assertNoError(t, err)
	assertNoEvents(t, Tournament.GetEvents())
}

func TestJoinNotCreatedTeam(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "id"
	var contestantId ContestantId = "id"
	Tournament.JoinTournament(&contestantId)
	Tournament.Commit() // flush events
	// when
	err := Tournament.JoinTeam(&teamId, &contestantId)
	// then
	assertError(t, err, ErrNoSuchTeamInTournament)
}

func TestJoinTeamInRemovedTournament(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "id"
	Tournament.CreateTeam(&teamId, "name")
	var contestantId ContestantId = "id"
	Tournament.JoinTournament(&contestantId)
	Tournament.Remove()
	Tournament.Commit() // flush events
	// when
	err := Tournament.JoinTeam(&teamId, &contestantId)
	// then
	assertError(t, err, ErrTournamentRemoved)
}

func TestJoinFullTeam(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "id"
	Tournament.CreateTeam(&teamId, "name")
	var contestantId1 ContestantId = "id1"
	var contestantId2 ContestantId = "id2"
	var contestantId3 ContestantId = "id3"
	Tournament.JoinTournament(&contestantId1)
	Tournament.JoinTournament(&contestantId2)
	Tournament.JoinTournament(&contestantId3)
	Tournament.JoinTeam(&teamId, &contestantId1)
	Tournament.JoinTeam(&teamId, &contestantId2)
	Tournament.Commit() // flush events
	// when
	err := Tournament.JoinTeam(&teamId, &contestantId3)
	// then
	assertError(t, err, ErrTeamFull)
}

func TestJoinTeamInStartedTournament(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "id"
	Tournament.CreateTeam(&teamId, "name")
	var contestantId ContestantId = "id"
	Tournament.JoinTournament(&contestantId)
	Tournament.Start()
	Tournament.Commit() // flush events
	// when
	err := Tournament.JoinTeam(&teamId, &contestantId)
	// then
	assertError(t, err, ErrTournamentAlreadyStarted)
}

// ------ Leave Team ------
func TestLeaveTeam(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "id"
	Tournament.CreateTeam(&teamId, "name")
	var contestantId ContestantId = "id"
	Tournament.JoinTournament(&contestantId)
	Tournament.JoinTeam(&teamId, &contestantId)
	Tournament.Commit() // flush events
	// when
	err := Tournament.LeaveTeam(&teamId, &contestantId)
	// then
	assertNoError(t, err)
	assertEvent(t, Tournament.GetEvents(), ContestantLeftTeam{ContestantId: contestantId, TeamId: teamId})
}

func TestLeaveNotJoinedTeam(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "id"
	Tournament.CreateTeam(&teamId, "name")
	var contestantId ContestantId = "id"
	Tournament.JoinTournament(&contestantId)
	Tournament.Commit() // flush events
	// when
	err := Tournament.LeaveTeam(&teamId, &contestantId)
	// then
	assertNoError(t, err)
	assertNoEvents(t, Tournament.GetEvents())
}

func TestLeaveTeamInRemovedTournament(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "id"
	Tournament.CreateTeam(&teamId, "name")
	var contestantId ContestantId = "id"
	Tournament.JoinTournament(&contestantId)
	Tournament.Remove()
	Tournament.Commit() // flush events
	// when
	err := Tournament.LeaveTeam(&teamId, &contestantId)
	// then
	assertError(t, err, ErrTournamentRemoved)
}

func TestLeaveTeamInStartedTournament(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "id"
	Tournament.CreateTeam(&teamId, "name")
	var contestantId ContestantId = "id"
	Tournament.JoinTournament(&contestantId)
	Tournament.Start()
	Tournament.Commit() // flush events
	// when
	err := Tournament.LeaveTeam(&teamId, &contestantId)
	// then
	assertError(t, err, ErrTournamentAlreadyStarted)
}

// ------ Helpers ------

func assertEvents(t *testing.T, events []any, expectedEvents []any) {
	for _, e := range expectedEvents {
		assertEvent(t, events, e)
	}
}

func assertEvent(t *testing.T, events []any, expectedEvent any) {
	var eventsString string
	for _, e := range events {
		if e == expectedEvent {
			return
		}
		eventsString += fmt.Sprintf("%[1]T: %+[1]v\n", e)
	}
	t.Errorf("expected event %[1]T: %+[1]v\nGot:\n%v", expectedEvent, eventsString)
}

func assertNoEvents(t *testing.T, events []any) {
	var eventsString string
	if len(events) != 0 {
		for _, e := range events {
			eventsString += fmt.Sprintf("%[1]T: %+[1]v\n", e)
		}
		t.Errorf("expected no events\nGot: %v", eventsString)
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("expected no error\nGot: %v", err)
	}
}

func assertError(t *testing.T, err error, expected error) {
	if err == nil {
		t.Errorf("expected error %v\nGot: nil", expected)
	}
	if err != expected {
		t.Errorf("expected error %v\nGot: %v", expected, err)
	}
}
