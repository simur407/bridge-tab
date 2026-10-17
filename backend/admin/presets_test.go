package main

import (
	"strings"
	"testing"

	tournament_cmd "bridge-tab/internal/tournament-management/application/command"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
	tournament_infra "bridge-tab/internal/tournament-management/infrastructure"

	"github.com/google/uuid"
)

func TestSixTeamPresetMovement(t *testing.T) {
	preset, ok := findPreset("6-teams-4-rounds")
	if !ok {
		t.Fatal("missing preset")
	}
	summary := preset.summary()
	for _, line := range []string{
		"1: 1:1-4 (A) 2:2-5 (C) 3:3-6 (D)",
		"2: 1:1-3 (B) 2:2-4 (D) 3:6-5 (A)",
		"3: 1:1-6 (C) 2:2-3 (A) 3:5-4 (B)",
		"4: 1:1-5 (D) 2:2-6 (B) 3:4-3 (C)",
		"A: 1, 2",
		"B: 3, 4",
		"C: 5, 6",
		"D: 7, 8",
	} {
		if !strings.Contains(summary, line) {
			t.Fatalf("summary missing %q:\n%s", line, summary)
		}
	}
}

func TestApplySixTeamPreset(t *testing.T) {
	preset, ok := findPreset("6-teams-4-rounds")
	if !ok {
		t.Fatal("missing preset")
	}
	repo := &memoryTournaments{}
	tournamentId := uuid.NewString()
	create := tournament_cmd.CreateTournamentCommand{TournamentId: tournamentId, Name: "Preset"}
	if err := create.Execute(repo); err != nil {
		t.Fatal(err)
	}
	if err := applyPreset(repo, tournamentId, preset); err != nil {
		t.Fatal(err)
	}

	id := tournament_domain.TournamentId(tournamentId)
	tournament, err := repo.Load(&id)
	if err != nil {
		t.Fatal(err)
	}

	if len(tournament.State.Teams) != 6 {
		t.Fatalf("expected 6 teams, got %d", len(tournament.State.Teams))
	}
	teams := map[int]tournament_domain.TeamId{}
	for _, team := range tournament.State.Teams {
		if team.State.Name != "" || len(team.State.Members) != 0 {
			t.Fatalf("team %d is not empty: %+v", team.State.Number, team.State)
		}
		teams[team.State.Number] = team.State.Id
	}
	for number := 1; number <= 6; number++ {
		if _, ok := teams[number]; !ok {
			t.Fatalf("missing team %d", number)
		}
	}

	if len(tournament.State.Tables) != 3 {
		t.Fatalf("expected 3 tables, got %d", len(tournament.State.Tables))
	}
	tables := map[int]tournament_domain.TableId{}
	for _, table := range tournament.State.Tables {
		tables[table.State.Number] = table.State.Id
	}

	if len(tournament.State.BoardProtocols) != 8 {
		t.Fatalf("expected 8 protocols, got %d", len(tournament.State.BoardProtocols))
	}
	wantVulnerable := []tournament_domain.Vulnerable{
		tournament_domain.None,
		tournament_domain.NS,
		tournament_domain.EW,
		tournament_domain.Both,
		tournament_domain.NS,
		tournament_domain.EW,
		tournament_domain.Both,
		tournament_domain.None,
	}
	for _, protocol := range tournament.State.BoardProtocols {
		if protocol.BoardNo < 1 || protocol.BoardNo > 8 {
			t.Fatalf("unexpected board %d", protocol.BoardNo)
		}
		if protocol.Vulnerable != wantVulnerable[protocol.BoardNo-1] {
			t.Fatalf("board %d vulnerability %v", protocol.BoardNo, protocol.Vulnerable)
		}
		if protocol.SetId == nil {
			t.Fatalf("board %d is not in a set", protocol.BoardNo)
		}
	}

	sets := map[string]tournament_domain.SetId{}
	if len(tournament.State.Sets) != 4 {
		t.Fatalf("expected 4 sets, got %d", len(tournament.State.Sets))
	}
	for _, set := range tournament.State.Sets {
		sets[set.Label] = set.Id
		want := map[string][]int{"A": {1, 2}, "B": {3, 4}, "C": {5, 6}, "D": {7, 8}}[set.Label]
		if len(set.BoardNos) != 2 || set.BoardNos[0] != want[0] || set.BoardNos[1] != want[1] {
			t.Fatalf("set %s boards %v", set.Label, set.BoardNos)
		}
	}

	if len(tournament.State.Schedule) != 12 {
		t.Fatalf("expected 12 schedule rows, got %d", len(tournament.State.Schedule))
	}
	for _, round := range preset.Rounds {
		for _, table := range round.Tables {
			found := false
			for _, entry := range tournament.State.Schedule {
				if entry.RoundNo == round.Number &&
					entry.TableId == tables[table.Table] &&
					entry.NsTeamId == teams[table.NS] &&
					entry.EwTeamId == teams[table.EW] &&
					entry.SetId == sets[table.Set] {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("missing schedule round %d table %d", round.Number, table.Table)
			}
		}
	}
}

func TestStandardVulnerableCycle(t *testing.T) {
	if standardVulnerable(1) != tournament_domain.None || standardVulnerable(17) != tournament_domain.None {
		t.Fatal("board 1 and 17 should be none")
	}
	if standardVulnerable(8) != tournament_domain.None || standardVulnerable(4) != tournament_domain.Both {
		t.Fatal("unexpected vulnerability")
	}
}

type memoryTournaments struct {
	byID map[string]*tournament_domain.Tournament
}

func (m *memoryTournaments) Load(id *tournament_domain.TournamentId) (*tournament_domain.Tournament, error) {
	tournament, ok := m.byID[string(*id)]
	if !ok {
		return nil, tournament_infra.ErrTournamentNotFound
	}
	return tournament, nil
}

func (m *memoryTournaments) Save(tournament *tournament_domain.Tournament) error {
	if m.byID == nil {
		m.byID = map[string]*tournament_domain.Tournament{}
	}
	m.byID[string(tournament.State.Id)] = tournament
	tournament.Commit()
	return nil
}
