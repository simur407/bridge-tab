package tournament_management_test

import (
	"fmt"
	"reflect"
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
	Tournament := startableTournament(t)
	assertNoError(t, Tournament.Start())
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
	Tournament.CreateTeam(&teamId, "name", 1)
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
	Tournament := startableTournament(t)
	Tournament.Commit()

	err := Tournament.Start()

	assertNoError(t, err)
	assertEvent(t, Tournament.GetEvents(), TournamentStarted{TournamentId: id, StartedAt: *Tournament.State.StartedAt})
}

func TestStartRequiresSchedule(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "id"
	var contestantId ContestantId = "id"
	assertNoError(t, Tournament.CreateTeam(&teamId, "name", 1))
	assertNoError(t, Tournament.JoinTournament(&contestantId))
	assertNoError(t, Tournament.JoinTeam(&teamId, &contestantId))

	err := Tournament.Start()
	assertError(t, err, ErrScheduleEmpty)
}

func TestStartRequiresEverySetOnSchedule(t *testing.T) {
	Tournament := startableTournament(t)
	assertNoError(t, Tournament.CreateBoardProtocol(2, None))
	assertNoError(t, Tournament.CreateSet("set-b", "B", []int{2}))

	err := Tournament.Start()
	assertError(t, err, ErrSetNotScheduled)
}

func TestStartStartedTournament(t *testing.T) {
	// given
	Tournament := startableTournament(t)
	assertNoError(t, Tournament.Start())
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
	Tournament.CreateTeam(&teamId, "name", 1)
	Tournament.Commit() // flush events
	// when
	err := Tournament.Start()
	// then
	assertError(t, err, ErrSomeTeamHasNoMembers)
}

func TestStartTournamentWithoutRequiringTeamMembers(t *testing.T) {
	Tournament := scheduledTournament(t)
	assertNoError(t, Tournament.SetRequireTeamMembers(false))
	Tournament.Commit()

	err := Tournament.Start()

	assertNoError(t, err)
	assertEvent(t, Tournament.GetEvents(), TournamentStarted{TournamentId: id, StartedAt: *Tournament.State.StartedAt})
}

func TestSetRequireTeamMembers(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	Tournament.Commit()

	err := Tournament.SetRequireTeamMembers(false)

	assertNoError(t, err)
	if Tournament.State.RequireTeamMembers {
		t.Errorf("expected RequireTeamMembers to be false")
	}
	assertEvent(t, Tournament.GetEvents(), TournamentRequireTeamMembersChanged{
		TournamentId:       id,
		RequireTeamMembers: false,
	})
}

func TestSetRequireTeamMembersUnchangedEmitsNothing(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	Tournament.Commit()

	err := Tournament.SetRequireTeamMembers(true)

	assertNoError(t, err)
	assertNoEvents(t, Tournament.GetEvents())
}

func TestSetRequireTeamMembersAfterStart(t *testing.T) {
	Tournament := startableTournament(t)
	assertNoError(t, Tournament.Start())
	Tournament.Commit()

	err := Tournament.SetRequireTeamMembers(false)

	assertError(t, err, ErrTournamentAlreadyStarted)
	assertNoEvents(t, Tournament.GetEvents())
}

func TestSetRequireTeamMembersOnRemovedTournament(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	assertNoError(t, Tournament.Remove())
	Tournament.Commit()

	err := Tournament.SetRequireTeamMembers(false)

	assertError(t, err, ErrTournamentRemoved)
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

// ------ Create Team ------
func TestCreateTeamAutoIncrementsNumberFromOne(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var first TeamId = "first"
	var second TeamId = "second"

	assertNoError(t, Tournament.CreateTeam(&first, "Alpha", 0))
	assertNoError(t, Tournament.CreateTeam(&second, "Beta", 0))

	if Tournament.State.Teams[0].State.Number != 1 || Tournament.State.Teams[1].State.Number != 2 {
		t.Errorf("expected numbers 1 and 2, got %d and %d", Tournament.State.Teams[0].State.Number, Tournament.State.Teams[1].State.Number)
	}
}

func TestCreateTeamUsesGivenNumber(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "team"
	var next TeamId = "next"

	assertNoError(t, Tournament.CreateTeam(&teamId, "Alpha", 4))
	assertNoError(t, Tournament.CreateTeam(&next, "Beta", 0))

	if Tournament.State.Teams[0].State.Number != 4 || Tournament.State.Teams[1].State.Number != 5 {
		t.Errorf("expected numbers 4 and 5, got %d and %d", Tournament.State.Teams[0].State.Number, Tournament.State.Teams[1].State.Number)
	}
}

func TestCreateTeamRejectsDuplicateNumber(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var first TeamId = "first"
	var second TeamId = "second"

	assertNoError(t, Tournament.CreateTeam(&first, "Alpha", 2))
	err := Tournament.CreateTeam(&second, "Beta", 2)

	assertError(t, err, ErrTeamNumberAlreadyExists)
}

func TestCreateTeamRejectsNegativeNumber(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "team"

	err := Tournament.CreateTeam(&teamId, "Alpha", -1)

	assertError(t, err, ErrInvalidTeamNumber)
}

func TestRenameTeam(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "team"
	assertNoError(t, Tournament.CreateTeam(&teamId, "Alpha", 1))
	Tournament.Commit()

	err := Tournament.RenameTeam(&teamId, "  Beta  ")

	assertNoError(t, err)
	assertEvent(t, Tournament.GetEvents(), TeamRenamed{TeamId: teamId, Name: "Beta"})
	if Tournament.State.Teams[0].State.Name != "Beta" {
		t.Errorf("expected name Beta, got %s", Tournament.State.Teams[0].State.Name)
	}
}

func TestRenameTeamAfterStart(t *testing.T) {
	Tournament := startableTournament(t)
	assertNoError(t, Tournament.Start())
	Tournament.Commit()

	var teamId TeamId = "ns"
	err := Tournament.RenameTeam(&teamId, "Beta")

	assertNoError(t, err)
	assertEvent(t, Tournament.GetEvents(), TeamRenamed{TeamId: teamId, Name: "Beta"})
	if Tournament.State.Teams[0].State.Name != "Beta" {
		t.Errorf("expected name Beta, got %s", Tournament.State.Teams[0].State.Name)
	}
}

func TestRenameTeamSameNameHasNoEvent(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "team"
	assertNoError(t, Tournament.CreateTeam(&teamId, "Alpha", 1))
	Tournament.Commit()

	err := Tournament.RenameTeam(&teamId, " Alpha ")

	assertNoError(t, err)
	assertNoEvents(t, Tournament.GetEvents())
}

func TestRenameTeamAllowsEmptyName(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "team"
	assertNoError(t, Tournament.CreateTeam(&teamId, "Alpha", 1))
	Tournament.Commit()

	err := Tournament.RenameTeam(&teamId, "   ")

	assertNoError(t, err)
	assertEvent(t, Tournament.GetEvents(), TeamRenamed{TeamId: teamId, Name: ""})
}

func TestRenameTeamRejectsDuplicateName(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var first TeamId = "first"
	var second TeamId = "second"
	assertNoError(t, Tournament.CreateTeam(&first, "Alpha", 1))
	assertNoError(t, Tournament.CreateTeam(&second, "Beta", 2))

	err := Tournament.RenameTeam(&second, "Alpha")

	assertError(t, err, ErrTeamNameAlreadyExists)
}

func TestRenameTeamRejectsMissingTeam(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "missing"

	err := Tournament.RenameTeam(&teamId, "Alpha")

	assertError(t, err, ErrNoSuchTeamInTournament)
}

func TestRenameTeamRejectsRemovedTournament(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "team"
	assertNoError(t, Tournament.CreateTeam(&teamId, "Alpha", 1))
	assertNoError(t, Tournament.Remove())

	err := Tournament.RenameTeam(&teamId, "Beta")

	assertError(t, err, ErrTournamentRemoved)
}

// ------ Join Team ------
func TestJoinTeamMatchesContestantIdRegardlessOfCase(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "498697ab-156d-47a7-bd65-020b31b73886"
	Tournament.CreateTeam(&teamId, "name", 1)
	// Postgres returns uuid values in lowercase. macOS uuidgen supplies uppercase.
	stored := ContestantId("d0ccd7c7-0323-4c9e-8ac7-40d4556577ef")
	Tournament.JoinTournament(&stored)
	Tournament.Commit()
	requested := ContestantId("D0CCD7C7-0323-4C9E-8AC7-40D4556577EF")

	// when
	err := Tournament.JoinTeam(&teamId, &requested)

	// then
	assertNoError(t, err)
	assertEvent(t, Tournament.GetEvents(), ContestantJoinedTeam{ContestantId: stored, TeamId: teamId})
}

func TestJoinTeam(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "id"
	Tournament.CreateTeam(&teamId, "name", 1)
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
	Tournament.CreateTeam(&teamId, "name", 1)
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
	Tournament.CreateTeam(&teamId, "name", 1)
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
	Tournament.CreateTeam(&teamId, "name", 1)
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
	Tournament := startableTournament(t)
	var teamId TeamId = "ns"
	var contestantId2 ContestantId = "id2"
	assertNoError(t, Tournament.JoinTournament(&contestantId2))
	assertNoError(t, Tournament.Start())
	Tournament.Commit()

	err := Tournament.JoinTeam(&teamId, &contestantId2)
	assertError(t, err, ErrTournamentAlreadyStarted)
}

// ------ Leave Team ------
func TestLeaveTeam(t *testing.T) {
	// given
	Tournament := CreateTournament(id, "name")
	var teamId TeamId = "id"
	Tournament.CreateTeam(&teamId, "name", 1)
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
	Tournament.CreateTeam(&teamId, "name", 1)
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
	Tournament.CreateTeam(&teamId, "name", 1)
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
	Tournament := startableTournament(t)
	assertNoError(t, Tournament.Start())
	Tournament.Commit()

	var teamId TeamId = "ns"
	var contestantId ContestantId = "ns-player"
	err := Tournament.LeaveTeam(&teamId, &contestantId)
	assertError(t, err, ErrTournamentAlreadyStarted)
}

// ------ Create Table ------
func TestCreateTableAutoIncrementsNumberFromOne(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var first TableId = "first"
	var second TableId = "second"

	assertNoError(t, Tournament.CreateTable(&first, 0))
	assertNoError(t, Tournament.CreateTable(&second, 0))

	if Tournament.State.Tables[0].State.Number != 1 || Tournament.State.Tables[1].State.Number != 2 {
		t.Errorf("expected numbers 1 and 2, got %d and %d", Tournament.State.Tables[0].State.Number, Tournament.State.Tables[1].State.Number)
	}
}

func TestCreateTableUsesGivenNumber(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var tableId TableId = "table"
	var next TableId = "next"

	assertNoError(t, Tournament.CreateTable(&tableId, 4))
	assertNoError(t, Tournament.CreateTable(&next, 0))

	if Tournament.State.Tables[0].State.Number != 4 || Tournament.State.Tables[1].State.Number != 5 {
		t.Errorf("expected numbers 4 and 5, got %d and %d", Tournament.State.Tables[0].State.Number, Tournament.State.Tables[1].State.Number)
	}
}

func TestCreateTableRejectsDuplicateNumber(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var first TableId = "first"
	var second TableId = "second"

	assertNoError(t, Tournament.CreateTable(&first, 2))
	err := Tournament.CreateTable(&second, 2)

	assertError(t, err, ErrTableNumberAlreadyExists)
}

func TestCreateTableRejectsNegativeNumber(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var tableId TableId = "table"

	err := Tournament.CreateTable(&tableId, -1)

	assertError(t, err, ErrInvalidTableNumber)
}

func TestDeleteTable(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var tableId TableId = "table"
	assertNoError(t, Tournament.CreateTable(&tableId, 1))
	Tournament.Commit()

	assertNoError(t, Tournament.DeleteTable(&tableId))
	assertEvent(t, Tournament.GetEvents(), TableRemoved{TournamentId: id, TableId: tableId})
}

func TestDeleteTableReferencedBySchedule(t *testing.T) {
	Tournament := scheduledTournament(t)
	var tableId TableId = "table"
	Tournament.Commit()

	err := Tournament.DeleteTable(&tableId)
	assertError(t, err, ErrTableReferencedBySchedule)
}

func TestCreateBoardProtocol(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	Tournament.Commit()

	assertNoError(t, Tournament.CreateBoardProtocol(1, None))
	assertEvent(t, Tournament.GetEvents(), BoardProtocolCreated{
		TournamentId: id,
		BoardNo:      1,
		Vulnerable:   None,
	})
}

func TestRemoveTournamentWithTables(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var tableId TableId = "table"
	assertNoError(t, Tournament.CreateTable(&tableId, 1))
	Tournament.Commit()

	assertNoError(t, Tournament.Remove())
	assertEvents(t, Tournament.GetEvents(), []any{
		TableRemoved{TournamentId: id, TableId: tableId},
		TournamentRemoved{TournamentId: id},
	})
}

// ------ Finish ------
func TestFinishTournament(t *testing.T) {
	Tournament := startableTournament(t)
	assertNoError(t, Tournament.Start())
	Tournament.Commit()

	err := Tournament.Finish()
	assertNoError(t, err)
	assertEvent(t, Tournament.GetEvents(), TournamentFinished{TournamentId: id, FinishedAt: *Tournament.State.FinishedAt})
}

func TestFinishNotStartedTournament(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	err := Tournament.Finish()
	assertError(t, err, ErrTournamentNotStarted)
}

func TestFinishAlreadyFinishedTournament(t *testing.T) {
	Tournament := startableTournament(t)
	assertNoError(t, Tournament.Start())
	assertNoError(t, Tournament.Finish())
	Tournament.Commit()

	err := Tournament.Finish()
	assertError(t, err, ErrTournamentAlreadyFinished)
}

// ------ Set ------
func TestCreateSet(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	assertNoError(t, Tournament.CreateBoardProtocol(1, None))
	assertNoError(t, Tournament.CreateBoardProtocol(2, NS))
	Tournament.Commit()

	var setId SetId = "set-a"
	assertNoError(t, Tournament.CreateSet(setId, "A", []int{1, 2}))
	assertEvent(t, Tournament.GetEvents(), SetCreated{
		TournamentId: id,
		SetId:        setId,
		Label:        "A",
		BoardNos:     []int{1, 2},
	})
}

func TestCreateSetRejectsBoardAlreadyInSet(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	assertNoError(t, Tournament.CreateBoardProtocol(1, None))
	assertNoError(t, Tournament.CreateBoardProtocol(2, None))
	assertNoError(t, Tournament.CreateSet("set-a", "A", []int{1, 2}))

	err := Tournament.CreateSet("set-b", "B", []int{1})
	assertError(t, err, ErrSetBoardAlreadyInSet)
}

func TestRemoveSet(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	assertNoError(t, Tournament.CreateBoardProtocol(1, None))
	var setId SetId = "set-a"
	assertNoError(t, Tournament.CreateSet(setId, "A", []int{1}))
	Tournament.Commit()

	assertNoError(t, Tournament.RemoveSet(setId))
	assertEvent(t, Tournament.GetEvents(), SetRemoved{TournamentId: id, SetId: setId})
}

func TestRemoveSetReferencedBySchedule(t *testing.T) {
	Tournament := scheduledTournament(t)

	err := Tournament.RemoveSet("set-a")
	assertError(t, err, ErrSetReferencedBySchedule)
}

func TestRemoveBoardProtocolBelongsToSet(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	assertNoError(t, Tournament.CreateBoardProtocol(1, None))
	assertNoError(t, Tournament.CreateSet("set-a", "A", []int{1}))

	err := Tournament.RemoveBoardProtocol(1)
	assertError(t, err, ErrBoardProtocolBelongsToSet)
}

func TestAddSchedule(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var tableId TableId = "table"
	var ns TeamId = "ns"
	var ew TeamId = "ew"
	var setId SetId = "set-a"
	assertNoError(t, Tournament.CreateTable(&tableId, 1))
	assertNoError(t, Tournament.CreateTeam(&ns, "ns", 1))
	assertNoError(t, Tournament.CreateTeam(&ew, "ew", 2))
	assertNoError(t, Tournament.CreateBoardProtocol(1, None))
	assertNoError(t, Tournament.CreateSet(setId, "A", []int{1}))
	Tournament.Commit()

	assertNoError(t, Tournament.AddSchedule(2, tableId, ns, ew, setId))
	assertEvent(t, Tournament.GetEvents(), ScheduleAdded{
		TournamentId: id,
		RoundNo:      2,
		TableId:      tableId,
		NsTeamId:     ns,
		EwTeamId:     ew,
		SetId:        setId,
	})
}

func TestAddScheduleRejectsSameSetTwiceInRound(t *testing.T) {
	Tournament := scheduledTournament(t)
	var other TableId = "other"
	var ns2 TeamId = "ns2"
	var ew2 TeamId = "ew2"
	assertNoError(t, Tournament.CreateTable(&other, 2))
	assertNoError(t, Tournament.CreateTeam(&ns2, "ns2", 3))
	assertNoError(t, Tournament.CreateTeam(&ew2, "ew2", 4))

	err := Tournament.AddSchedule(1, other, ns2, ew2, "set-a")
	assertError(t, err, ErrScheduleSetAlreadyUsed)
}

func TestAddScheduleAllowsSameSetInAnotherRound(t *testing.T) {
	Tournament := scheduledTournament(t)
	var other TableId = "other"
	var ns2 TeamId = "ns2"
	var ew2 TeamId = "ew2"
	assertNoError(t, Tournament.CreateTable(&other, 2))
	assertNoError(t, Tournament.CreateTeam(&ns2, "ns2", 3))
	assertNoError(t, Tournament.CreateTeam(&ew2, "ew2", 4))

	assertNoError(t, Tournament.AddSchedule(2, other, ns2, ew2, "set-a"))
}

func TestAddScheduleRejectsTeamTwiceInRound(t *testing.T) {
	Tournament := scheduledTournament(t)
	var other TableId = "other"
	var ew2 TeamId = "ew2"
	assertNoError(t, Tournament.CreateTable(&other, 2))
	assertNoError(t, Tournament.CreateTeam(&ew2, "ew2", 3))

	err := Tournament.AddSchedule(1, other, "ns", ew2, "set-a")
	assertError(t, err, ErrScheduleTeamAlreadyUsed)
}

func TestAddScheduleRejectsTableTwiceInRound(t *testing.T) {
	Tournament := scheduledTournament(t)
	var ns2 TeamId = "ns2"
	var ew2 TeamId = "ew2"
	assertNoError(t, Tournament.CreateTeam(&ns2, "ns2", 3))
	assertNoError(t, Tournament.CreateTeam(&ew2, "ew2", 4))

	err := Tournament.AddSchedule(1, "table", ns2, ew2, "set-a")
	assertError(t, err, ErrScheduleTableAlreadyUsed)
}

func TestAddScheduleRejectsSameTeam(t *testing.T) {
	Tournament := CreateTournament(id, "name")
	var tableId TableId = "table"
	var teamId TeamId = "team"
	assertNoError(t, Tournament.CreateTable(&tableId, 1))
	assertNoError(t, Tournament.CreateTeam(&teamId, "team", 1))
	assertNoError(t, Tournament.CreateBoardProtocol(1, None))
	assertNoError(t, Tournament.CreateSet("set-a", "A", []int{1}))

	err := Tournament.AddSchedule(1, tableId, teamId, teamId, "set-a")
	assertError(t, err, ErrScheduleSameTeam)
}

func TestDeleteTeamReferencedBySchedule(t *testing.T) {
	Tournament := scheduledTournament(t)
	var ns TeamId = "ns"

	err := Tournament.DeleteTeam(&ns)
	assertError(t, err, ErrTeamReferencedBySchedule)
}

func TestRemoveSchedule(t *testing.T) {
	Tournament := scheduledTournament(t)
	Tournament.Commit()

	assertNoError(t, Tournament.RemoveSchedule(1, "table"))
	assertEvent(t, Tournament.GetEvents(), ScheduleRemoved{TournamentId: id, RoundNo: 1, TableId: "table"})
}

func scheduledTournament(t *testing.T) *Tournament {
	t.Helper()
	Tournament := CreateTournament(id, "name")
	var tableId TableId = "table"
	var ns TeamId = "ns"
	var ew TeamId = "ew"
	assertNoError(t, Tournament.CreateTable(&tableId, 1))
	assertNoError(t, Tournament.CreateTeam(&ns, "ns", 1))
	assertNoError(t, Tournament.CreateTeam(&ew, "ew", 2))
	assertNoError(t, Tournament.CreateBoardProtocol(1, None))
	assertNoError(t, Tournament.CreateSet("set-a", "A", []int{1}))
	assertNoError(t, Tournament.AddSchedule(1, tableId, ns, ew, "set-a"))
	return Tournament
}

func startableTournament(t *testing.T) *Tournament {
	t.Helper()
	tournament := scheduledTournament(t)
	var ns TeamId = "ns"
	var ew TeamId = "ew"
	var nsPlayer ContestantId = "ns-player"
	var ewPlayer ContestantId = "ew-player"
	assertNoError(t, tournament.JoinTournament(&nsPlayer))
	assertNoError(t, tournament.JoinTournament(&ewPlayer))
	assertNoError(t, tournament.JoinTeam(&ns, &nsPlayer))
	assertNoError(t, tournament.JoinTeam(&ew, &ewPlayer))
	return tournament
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
		if reflect.DeepEqual(e, expectedEvent) {
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
