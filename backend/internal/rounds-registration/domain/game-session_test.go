package rounds_registration_test

import (
	"fmt"
	"reflect"
	"testing"

	. "bridge-tab/internal/rounds-registration/domain"
)

const id GameSessionId = "id"

// ------ Edit Round ------
func TestEditRound(t *testing.T) {
	// given
	var ns TeamId = "ns"
	var ew TeamId = "ew"
	round := &Round{DealNo: 1, NsTeam: &ns, EwTeam: &ew}
	GameSession, err := StartGameSession(id, []*Team{
		{Id: ns, Players: []PlayerId{"n"}},
		{Id: ew, Players: []PlayerId{"e"}},
	}, []*Round{round})
	assertNoError(t, err)
	assertNoError(t, GameSession.AddRoundScore(round.DealNo, ns, ew, "3N", 9, "N", "KH"))
	GameSession.Commit() // flush events
	// when
	err = GameSession.EditRoundScore(round.DealNo, ns, ew, "4H", 10, "S", "AS")
	// then
	assertNoError(t, err)
	assertEvent(t, GameSession.GetEvents(), RoundEdited{
		GameSessionId: id,
		DealNo:        round.DealNo,
		NsTeamId:      ns,
		EwTeamId:      ew,
		Contract:      "4H",
		Tricks:        10,
		Declarer:      "S",
		OpeningLead:   "AS",
	})
}

func TestEditRoundNotPlayed(t *testing.T) {
	// given
	var ns TeamId = "ns"
	var ew TeamId = "ew"
	round := &Round{DealNo: 1, NsTeam: &ns, EwTeam: &ew}
	GameSession, err := StartGameSession(id, []*Team{
		{Id: ns, Players: []PlayerId{"n"}},
		{Id: ew, Players: []PlayerId{"e"}},
	}, []*Round{round})
	assertNoError(t, err)
	GameSession.Commit() // flush events
	// when
	err = GameSession.EditRoundScore(round.DealNo, ns, ew, "4H", 10, "S", "AS")
	// then
	assertError(t, err, ErrRoundNotPlayed)
	assertNoEvents(t, GameSession.GetEvents())
}

func TestEditMissingRound(t *testing.T) {
	// given
	var ns TeamId = "ns"
	var ew TeamId = "ew"
	round := &Round{DealNo: 1, NsTeam: &ns, EwTeam: &ew}
	GameSession, err := StartGameSession(id, []*Team{
		{Id: ns, Players: []PlayerId{"n"}},
		{Id: ew, Players: []PlayerId{"e"}},
	}, []*Round{round})
	assertNoError(t, err)
	assertNoError(t, GameSession.AddRoundScore(round.DealNo, ns, ew, "3N", 9, "N", "KH"))
	GameSession.Commit() // flush events
	// when
	err = GameSession.EditRoundScore(2, ns, ew, "4H", 10, "S", "AS")
	// then
	assertError(t, err, ErrRoundNotFound)
	assertNoEvents(t, GameSession.GetEvents())
}

// ------ Helpers ------

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
