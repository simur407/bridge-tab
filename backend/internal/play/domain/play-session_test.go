package play_test

import (
	"errors"
	"testing"

	. "bridge-tab/internal/play/domain"
)

func TestProvisionPlaySession(t *testing.T) {
	ns := &Team{Id: "ns", Number: 1, Players: []PlayerId{"n"}}
	ew := &Team{Id: "ew", Number: 2, Players: []PlayerId{"e"}}
	schedule := &TableSchedule{
		RoundNo:     1,
		TableNumber: 1,
		NsTeamId:    ns.Id,
		EwTeamId:    ew.Id,
		SetId:       "set-a",
		SetLabel:    "A",
		Boards:      []TableScheduleBoard{{DealNo: 1, Vulnerable: "NS"}},
	}

	session, err := ProvisionPlaySession("session", []*Team{ns, ew}, []*TableSchedule{schedule})
	if err != nil {
		t.Fatal(err)
	}

	events := session.GetEvents()
	if len(events) != 1 {
		t.Fatalf("expected one event, got %d", len(events))
	}
	provisioned, ok := events[0].(PlaySessionProvisioned)
	if !ok {
		t.Fatalf("expected PlaySessionProvisioned, got %T", events[0])
	}
	if provisioned.PlaySessionId != "session" || len(provisioned.Teams) != 2 || len(provisioned.TableSchedules) != 1 {
		t.Fatalf("unexpected event: %+v", provisioned)
	}
	if len(session.Matches()) != 1 {
		t.Fatalf("expected one match derived from the table schedule, got %d", len(session.Matches()))
	}
}

func TestProvisionPlaySessionRequiresTeamsAndTableSchedules(t *testing.T) {
	_, err := ProvisionPlaySession("session", nil, []*TableSchedule{{RoundNo: 1}})
	if !errors.Is(err, ErrTeamsAreEmpty) {
		t.Fatalf("expected empty teams, got %v", err)
	}

	_, err = ProvisionPlaySession("session", []*Team{{Id: "ns", Number: 1}}, nil)
	if !errors.Is(err, ErrMatchesAreEmpty) {
		t.Fatalf("expected empty table schedules, got %v", err)
	}
}

func TestMatchKeepsProtocolOrientation(t *testing.T) {
	session := mustSession(t)
	session.Commit()

	team, err := session.TeamByMember("east")
	if err != nil {
		t.Fatal(err)
	}
	if team.Number != 2 {
		t.Fatalf("expected team 2, got %d", team.Number)
	}

	match, err := session.Match(1, "ew", "ns")
	if err != nil {
		t.Fatal(err)
	}
	if match.NsTeamId != "ns" || match.EwTeamId != "ew" || match.Vulnerable != "NS" || match.RoundNo != 1 || match.TableNumber != 1 {
		t.Fatalf("unexpected match: %+v", match)
	}
}

func TestMatchNotFoundWhenDealIsNotOnTableSchedule(t *testing.T) {
	session := mustSession(t)

	_, err := session.Match(2, "ns", "ew")
	if !errors.Is(err, ErrMatchNotFound) {
		t.Fatalf("expected missing deal, got %v", err)
	}
}

func TestMatchNotFound(t *testing.T) {
	session := mustSession(t)

	_, err := session.Match(1, "ns", "other")
	if !errors.Is(err, ErrMatchNotFound) {
		t.Fatalf("expected missing match, got %v", err)
	}
}

func mustSession(t *testing.T) *PlaySession {
	t.Helper()
	session, err := ProvisionPlaySession("session", []*Team{
		{Id: "ns", Number: 1, Players: []PlayerId{"north"}},
		{Id: "ew", Number: 2, Players: []PlayerId{"east"}},
	}, []*TableSchedule{{
		RoundNo:     1,
		TableNumber: 1,
		NsTeamId:    "ns",
		EwTeamId:    "ew",
		SetLabel:    "A",
		Boards:      []TableScheduleBoard{{DealNo: 1, Vulnerable: "NS"}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	return session
}
