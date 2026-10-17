package play

import "testing"

func TestBuildRoundStatusGroupsPairsByRoundAndTable(t *testing.T) {
	board := BuildRoundStatus([]ScheduledMatchDto{
		{RoundNo: 2, TableNumber: 1, DealNo: 5, NsTeamNumber: 1, EwTeamNumber: 4},
		{RoundNo: 1, TableNumber: 3, DealNo: 4, NsTeamNumber: 5, EwTeamNumber: 6},
		{RoundNo: 1, TableNumber: 2, DealNo: 3, NsTeamNumber: 3, EwTeamNumber: 4},
		{RoundNo: 1, TableNumber: 1, DealNo: 2, NsTeamNumber: 1, EwTeamNumber: 2},
		{RoundNo: 1, TableNumber: 1, DealNo: 1, NsTeamNumber: 1, EwTeamNumber: 2, Contract: "3NT"},
		{RoundNo: 2, TableNumber: 2, DealNo: 7, NsTeamNumber: 3, EwTeamNumber: 2, Contract: "Pass"},
		{RoundNo: 2, TableNumber: 2, DealNo: 6, NsTeamNumber: 3, EwTeamNumber: 2, Contract: "4H"},
	})

	if len(board.Tables) != 3 || board.Tables[0] != 1 || board.Tables[1] != 2 || board.Tables[2] != 3 {
		t.Fatalf("tables: %+v", board.Tables)
	}
	if len(board.Rows) != 2 || board.Rows[0].RoundNo != 1 || board.Rows[1].RoundNo != 2 {
		t.Fatalf("rows: %+v", board.Rows)
	}

	round1 := board.Rows[0].Cells
	if !round1[0].Seated || round1[0].NsTeamNumber != 1 || round1[0].EwTeamNumber != 2 || round1[0].Status != "in-progress" {
		t.Fatalf("round 1 table 1: %+v", round1[0])
	}
	if !round1[1].Seated || round1[1].NsTeamNumber != 3 || round1[1].EwTeamNumber != 4 || round1[1].Status != "pending" {
		t.Fatalf("round 1 table 2: %+v", round1[1])
	}
	if !round1[2].Seated || round1[2].NsTeamNumber != 5 || round1[2].EwTeamNumber != 6 || round1[2].Status != "pending" {
		t.Fatalf("round 1 table 3: %+v", round1[2])
	}

	round2 := board.Rows[1].Cells
	if !round2[0].Seated || round2[0].NsTeamNumber != 1 || round2[0].EwTeamNumber != 4 || round2[0].Status != "pending" {
		t.Fatalf("round 2 table 1: %+v", round2[0])
	}
	if !round2[1].Seated || round2[1].NsTeamNumber != 3 || round2[1].EwTeamNumber != 2 || round2[1].Status != "completed" {
		t.Fatalf("round 2 table 2: %+v", round2[1])
	}
	if round2[2].Seated {
		t.Fatalf("round 2 table 3 should be empty: %+v", round2[2])
	}
}

func TestBuildRoundStatusEmpty(t *testing.T) {
	board := BuildRoundStatus(nil)
	if len(board.Tables) != 0 || len(board.Rows) != 0 {
		t.Fatalf("expected empty board, got %+v", board)
	}
}
