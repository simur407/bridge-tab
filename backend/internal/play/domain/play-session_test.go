package play_test

import (
	"errors"
	"testing"

	. "bridge-tab/internal/play/domain"
)

func TestProvisionPlaySession(t *testing.T) {
	ns := &Team{Id: "ns", Number: 1, Players: []PlayerId{"n"}}
	ew := &Team{Id: "ew", Number: 2, Players: []PlayerId{"e"}}
	match := &Match{DealNo: 1, NsTeamId: ns.Id, EwTeamId: ew.Id, Vulnerable: "NS"}

	session, err := ProvisionPlaySession("tournament", []*Team{ns, ew}, []*Match{match})
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
	if provisioned.TournamentId != "tournament" || len(provisioned.Teams) != 2 || len(provisioned.Matches) != 1 {
		t.Fatalf("unexpected event: %+v", provisioned)
	}
}

func TestProvisionPlaySessionRequiresTeamsAndMatches(t *testing.T) {
	_, err := ProvisionPlaySession("tournament", nil, []*Match{{DealNo: 1}})
	if !errors.Is(err, ErrTeamsAreEmpty) {
		t.Fatalf("expected empty teams, got %v", err)
	}

	_, err = ProvisionPlaySession("tournament", []*Team{{Id: "ns", Number: 1}}, nil)
	if !errors.Is(err, ErrMatchesAreEmpty) {
		t.Fatalf("expected empty matches, got %v", err)
	}
}

func TestMatchKeepsProtocolOrientation(t *testing.T) {
	session, err := ProvisionPlaySession("tournament", []*Team{
		{Id: "ns", Number: 1, Players: []PlayerId{"north"}},
		{Id: "ew", Number: 2, Players: []PlayerId{"east"}},
	}, []*Match{
		{DealNo: 1, NsTeamId: "ns", EwTeamId: "ew", Vulnerable: "NS"},
	})
	if err != nil {
		t.Fatal(err)
	}
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
	if match.NsTeamId != "ns" || match.EwTeamId != "ew" || match.Vulnerable != "NS" {
		t.Fatalf("unexpected match: %+v", match)
	}
}

func TestMatchNotFound(t *testing.T) {
	session, err := ProvisionPlaySession("tournament", []*Team{
		{Id: "ns", Number: 1},
		{Id: "ew", Number: 2},
	}, []*Match{
		{DealNo: 1, NsTeamId: "ns", EwTeamId: "ew", Vulnerable: "None"},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = session.Match(1, "ns", "other")
	if !errors.Is(err, ErrMatchNotFound) {
		t.Fatalf("expected missing match, got %v", err)
	}
}
