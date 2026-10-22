package main

import (
	"testing"

	play_query "bridge-tab/internal/play/application/query"
)

func TestSeatForTeam(t *testing.T) {
	matches := []play_query.ScheduledMatchDto{
		{DealNo: 1, RoundNo: 2, TableNumber: 5, NsTeamNumber: 1, EwTeamNumber: 2, Contract: ""},
		{DealNo: 2, RoundNo: 2, TableNumber: 5, NsTeamNumber: 1, EwTeamNumber: 2, Contract: "1S"},
		{DealNo: 3, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 4, Contract: "3NT"},
		{DealNo: 4, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 4, Contract: "2H"},
	}

	seat, ok := seatForTeam(matches, 1)
	if !ok || !seat.Current || seat.RoundNo != 2 || seat.TableNumber != 5 || seat.NsTeamNumber != 1 || seat.EwTeamNumber != 2 {
		t.Fatalf("current seat: %+v ok=%v", seat, ok)
	}

	done := []play_query.ScheduledMatchDto{
		{DealNo: 3, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 4, Contract: "3NT"},
		{DealNo: 1, RoundNo: 2, TableNumber: 5, NsTeamNumber: 1, EwTeamNumber: 2, Contract: "1S"},
	}
	seat, ok = seatForTeam(done, 1)
	if !ok || seat.Current || seat.RoundNo != 2 || seat.TableNumber != 5 {
		t.Fatalf("finished seat: %+v ok=%v", seat, ok)
	}

	if _, ok := seatForTeam(matches, 9); ok {
		t.Fatal("unknown team should have no seat")
	}
}
