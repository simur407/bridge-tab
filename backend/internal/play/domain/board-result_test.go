package play_test

import (
	"fmt"
	"reflect"
	"testing"

	. "bridge-tab/internal/play/domain"
)

func TestRecordBoardResult(t *testing.T) {
	result := RecordBoardResult("session", 3, "ns", "ew", "3N", 9, "N", "KH")

	if result.State.DealNo != 3 || result.State.Contract != "3N" || result.State.Tricks != 9 {
		t.Fatalf("unexpected state: %+v", result.State)
	}
	assertEvent(t, result.GetEvents(), BoardResultRecorded{
		PlaySessionId: "session",
		DealNo:        3,
		NsTeamId:      "ns",
		EwTeamId:      "ew",
		Contract:      "3N",
		Tricks:        9,
		Declarer:      "N",
		OpeningLead:   "KH",
	})
}

func TestEditBoardResult(t *testing.T) {
	result := RecordBoardResult("session", 3, "ns", "ew", "3N", 9, "N", "KH")
	result.Commit()

	result.Edit("4H", 10, "S", "AS")

	if result.State.Contract != "4H" || result.State.Tricks != 10 || result.State.Declarer != "S" || result.State.OpeningLead != "AS" {
		t.Fatalf("unexpected state: %+v", result.State)
	}
	assertEvent(t, result.GetEvents(), BoardResultEdited{
		PlaySessionId: "session",
		DealNo:        3,
		NsTeamId:      "ns",
		EwTeamId:      "ew",
		Contract:      "4H",
		Tricks:        10,
		Declarer:      "S",
		OpeningLead:   "AS",
	})
}

func assertEvent(t *testing.T, events []any, expectedEvent any) {
	t.Helper()
	var eventsString string
	for _, event := range events {
		if reflect.DeepEqual(event, expectedEvent) {
			return
		}
		eventsString += fmt.Sprintf("%[1]T: %+[1]v\n", event)
	}
	t.Errorf("expected event %[1]T: %+[1]v\nGot:\n%v", expectedEvent, eventsString)
}
