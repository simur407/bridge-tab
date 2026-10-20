package main

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	tournament_cmd "bridge-tab/internal/tournament-management/application/command"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/google/uuid"
)

type presetSet struct {
	Label  string
	Boards []int
}

type presetTable struct {
	Table int
	NS    int
	EW    int
	Set   string
}

type presetRound struct {
	Number int
	Tables []presetTable
}

type tournamentPreset struct {
	ID     string
	Label  string
	Sets   []presetSet
	Rounds []presetRound
}

type presetView struct {
	ID      string
	Label   string
	Hint    string
	Summary string
}

var tournamentPresets = []tournamentPreset{
	{
		ID:    "6-teams-4-rounds",
		Label: "6 drużyn, 4 rundy",
		Sets: []presetSet{
			{Label: "A", Boards: []int{1, 2}},
			{Label: "B", Boards: []int{3, 4}},
			{Label: "C", Boards: []int{5, 6}},
			{Label: "D", Boards: []int{7, 8}},
		},
		Rounds: []presetRound{
			{Number: 1, Tables: []presetTable{
				{Table: 1, NS: 1, EW: 4, Set: "A"},
				{Table: 2, NS: 2, EW: 5, Set: "C"},
				{Table: 3, NS: 3, EW: 6, Set: "D"},
			}},
			{Number: 2, Tables: []presetTable{
				{Table: 1, NS: 1, EW: 3, Set: "B"},
				{Table: 2, NS: 2, EW: 4, Set: "D"},
				{Table: 3, NS: 6, EW: 5, Set: "A"},
			}},
			{Number: 3, Tables: []presetTable{
				{Table: 1, NS: 1, EW: 6, Set: "C"},
				{Table: 2, NS: 2, EW: 3, Set: "A"},
				{Table: 3, NS: 5, EW: 4, Set: "B"},
			}},
			{Number: 4, Tables: []presetTable{
				{Table: 1, NS: 1, EW: 5, Set: "D"},
				{Table: 2, NS: 2, EW: 6, Set: "B"},
				{Table: 3, NS: 4, EW: 3, Set: "C"},
			}},
		},
	},
	{
		ID:    "10-teams-9-rounds",
		Label: "10 drużyn, 9 rund",
		Sets: []presetSet{
			{Label: "A", Boards: []int{1, 2}},
			{Label: "B", Boards: []int{3, 4}},
			{Label: "C", Boards: []int{5, 6}},
			{Label: "D", Boards: []int{7, 8}},
			{Label: "E", Boards: []int{9, 10}},
			{Label: "F", Boards: []int{11, 12}},
			{Label: "G", Boards: []int{13, 14}},
			{Label: "H", Boards: []int{15, 16}},
			{Label: "I", Boards: []int{17, 18}},
		},
		Rounds: []presetRound{
			{Number: 1, Tables: []presetTable{
				{Table: 1, NS: 10, EW: 1, Set: "A"},
				{Table: 2, NS: 9, EW: 8, Set: "B"},
				{Table: 3, NS: 6, EW: 4, Set: "C"},
				{Table: 4, NS: 3, EW: 7, Set: "H"},
				{Table: 5, NS: 5, EW: 2, Set: "I"},
			}},
			{Number: 2, Tables: []presetTable{
				{Table: 1, NS: 10, EW: 2, Set: "B"},
				{Table: 2, NS: 1, EW: 9, Set: "C"},
				{Table: 3, NS: 7, EW: 5, Set: "D"},
				{Table: 4, NS: 4, EW: 8, Set: "I"},
				{Table: 5, NS: 6, EW: 3, Set: "A"},
			}},
			{Number: 3, Tables: []presetTable{
				{Table: 1, NS: 10, EW: 3, Set: "C"},
				{Table: 2, NS: 2, EW: 1, Set: "D"},
				{Table: 3, NS: 8, EW: 6, Set: "E"},
				{Table: 4, NS: 5, EW: 9, Set: "A"},
				{Table: 5, NS: 7, EW: 4, Set: "B"},
			}},
			{Number: 4, Tables: []presetTable{
				{Table: 1, NS: 10, EW: 4, Set: "D"},
				{Table: 2, NS: 3, EW: 2, Set: "E"},
				{Table: 3, NS: 9, EW: 7, Set: "F"},
				{Table: 4, NS: 6, EW: 1, Set: "B"},
				{Table: 5, NS: 8, EW: 5, Set: "C"},
			}},
			{Number: 5, Tables: []presetTable{
				{Table: 1, NS: 10, EW: 5, Set: "E"},
				{Table: 2, NS: 4, EW: 3, Set: "F"},
				{Table: 3, NS: 1, EW: 8, Set: "G"},
				{Table: 4, NS: 7, EW: 2, Set: "C"},
				{Table: 5, NS: 9, EW: 6, Set: "D"},
			}},
			{Number: 6, Tables: []presetTable{
				{Table: 1, NS: 10, EW: 6, Set: "F"},
				{Table: 2, NS: 5, EW: 4, Set: "G"},
				{Table: 3, NS: 2, EW: 9, Set: "H"},
				{Table: 4, NS: 8, EW: 3, Set: "D"},
				{Table: 5, NS: 1, EW: 7, Set: "E"},
			}},
			{Number: 7, Tables: []presetTable{
				{Table: 1, NS: 10, EW: 7, Set: "G"},
				{Table: 2, NS: 6, EW: 5, Set: "H"},
				{Table: 3, NS: 3, EW: 1, Set: "I"},
				{Table: 4, NS: 9, EW: 4, Set: "E"},
				{Table: 5, NS: 2, EW: 8, Set: "F"},
			}},
			{Number: 8, Tables: []presetTable{
				{Table: 1, NS: 10, EW: 8, Set: "H"},
				{Table: 2, NS: 7, EW: 6, Set: "I"},
				{Table: 3, NS: 4, EW: 2, Set: "A"},
				{Table: 4, NS: 1, EW: 5, Set: "F"},
				{Table: 5, NS: 3, EW: 9, Set: "G"},
			}},
			{Number: 9, Tables: []presetTable{
				{Table: 1, NS: 10, EW: 9, Set: "I"},
				{Table: 2, NS: 8, EW: 7, Set: "A"},
				{Table: 3, NS: 5, EW: 3, Set: "B"},
				{Table: 4, NS: 2, EW: 6, Set: "G"},
				{Table: 5, NS: 4, EW: 1, Set: "H"},
			}},
		},
	},
	{
		ID:    "12-teams-11-round-3-boards",
		Label: "12 drużyn, 11 rund, 3 rozdania",
		Sets: []presetSet{
			{Label: "A", Boards: []int{1, 2, 3}},
			{Label: "B", Boards: []int{4, 5, 6}},
			{Label: "C", Boards: []int{7, 8, 9}},
			{Label: "D", Boards: []int{10, 11, 12}},
			{Label: "E", Boards: []int{13, 14, 15}},
			{Label: "F", Boards: []int{16, 17, 18}},
			{Label: "G", Boards: []int{19, 20, 21}},
			{Label: "H", Boards: []int{22, 23, 24}},
			{Label: "I", Boards: []int{25, 26, 27}},
			{Label: "J", Boards: []int{28, 29, 30}},
			{Label: "K", Boards: []int{31, 32, 33}},
		},
		Rounds: []presetRound{
			{Number: 1, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 1, Set: "A"},
				{Table: 2, NS: 6, EW: 7, Set: "D"},
				{Table: 3, NS: 11, EW: 4, Set: "F"},
				{Table: 4, NS: 3, EW: 9, Set: "H"},
				{Table: 5, NS: 2, EW: 5, Set: "I"},
				{Table: 6, NS: 10, EW: 8, Set: "K"},
			}},
			{Number: 2, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 2, Set: "B"},
				{Table: 2, NS: 7, EW: 8, Set: "E"},
				{Table: 3, NS: 1, EW: 5, Set: "G"},
				{Table: 4, NS: 4, EW: 10, Set: "I"},
				{Table: 5, NS: 3, EW: 6, Set: "J"},
				{Table: 6, NS: 11, EW: 9, Set: "A"},
			}},
			{Number: 3, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 3, Set: "C"},
				{Table: 2, NS: 8, EW: 9, Set: "F"},
				{Table: 3, NS: 2, EW: 6, Set: "H"},
				{Table: 4, NS: 5, EW: 11, Set: "J"},
				{Table: 5, NS: 4, EW: 7, Set: "K"},
				{Table: 6, NS: 1, EW: 10, Set: "B"},
			}},
			{Number: 4, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 4, Set: "D"},
				{Table: 2, NS: 9, EW: 10, Set: "G"},
				{Table: 3, NS: 3, EW: 7, Set: "I"},
				{Table: 4, NS: 6, EW: 1, Set: "K"},
				{Table: 5, NS: 5, EW: 8, Set: "A"},
				{Table: 6, NS: 2, EW: 11, Set: "C"},
			}},
			{Number: 5, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 5, Set: "E"},
				{Table: 2, NS: 10, EW: 11, Set: "H"},
				{Table: 3, NS: 4, EW: 8, Set: "J"},
				{Table: 4, NS: 7, EW: 2, Set: "A"},
				{Table: 5, NS: 6, EW: 9, Set: "B"},
				{Table: 6, NS: 3, EW: 1, Set: "D"},
			}},
			{Number: 6, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 6, Set: "F"},
				{Table: 2, NS: 11, EW: 1, Set: "I"},
				{Table: 3, NS: 5, EW: 9, Set: "K"},
				{Table: 4, NS: 8, EW: 3, Set: "B"},
				{Table: 5, NS: 7, EW: 10, Set: "C"},
				{Table: 6, NS: 4, EW: 2, Set: "E"},
			}},
			{Number: 7, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 7, Set: "G"},
				{Table: 2, NS: 1, EW: 2, Set: "J"},
				{Table: 3, NS: 6, EW: 10, Set: "A"},
				{Table: 4, NS: 9, EW: 4, Set: "C"},
				{Table: 5, NS: 8, EW: 11, Set: "D"},
				{Table: 6, NS: 5, EW: 3, Set: "F"},
			}},
			{Number: 8, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 8, Set: "H"},
				{Table: 2, NS: 2, EW: 3, Set: "K"},
				{Table: 3, NS: 7, EW: 11, Set: "B"},
				{Table: 4, NS: 10, EW: 5, Set: "D"},
				{Table: 5, NS: 9, EW: 1, Set: "E"},
				{Table: 6, NS: 6, EW: 4, Set: "G"},
			}},
			{Number: 9, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 9, Set: "I"},
				{Table: 2, NS: 3, EW: 4, Set: "A"},
				{Table: 3, NS: 8, EW: 1, Set: "C"},
				{Table: 4, NS: 11, EW: 6, Set: "E"},
				{Table: 5, NS: 10, EW: 2, Set: "F"},
				{Table: 6, NS: 7, EW: 5, Set: "H"},
			}},
			{Number: 10, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 10, Set: "J"},
				{Table: 2, NS: 4, EW: 5, Set: "B"},
				{Table: 3, NS: 9, EW: 2, Set: "D"},
				{Table: 4, NS: 1, EW: 7, Set: "F"},
				{Table: 5, NS: 11, EW: 3, Set: "G"},
				{Table: 6, NS: 8, EW: 6, Set: "I"},
			}},
			{Number: 11, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 11, Set: "K"},
				{Table: 2, NS: 5, EW: 6, Set: "C"},
				{Table: 3, NS: 10, EW: 3, Set: "E"},
				{Table: 4, NS: 2, EW: 8, Set: "G"},
				{Table: 5, NS: 1, EW: 4, Set: "H"},
				{Table: 6, NS: 9, EW: 7, Set: "J"},
			}},
		},
	},
	{
		ID:    "12-teams-11-round-2-boards",
		Label: "12 drużyn, 11 rund, 2 rozdania",
		Sets: []presetSet{
			{Label: "A", Boards: []int{1, 2}},
			{Label: "B", Boards: []int{3, 4}},
			{Label: "C", Boards: []int{5, 6}},
			{Label: "D", Boards: []int{7, 8}},
			{Label: "E", Boards: []int{9, 10}},
			{Label: "F", Boards: []int{11, 12}},
			{Label: "G", Boards: []int{13, 14}},
			{Label: "H", Boards: []int{15, 16}},
			{Label: "I", Boards: []int{17, 18}},
			{Label: "J", Boards: []int{19, 20}},
			{Label: "K", Boards: []int{21, 22}},
		},
		Rounds: []presetRound{
			{Number: 1, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 1, Set: "A"},
				{Table: 2, NS: 6, EW: 7, Set: "D"},
				{Table: 3, NS: 11, EW: 4, Set: "F"},
				{Table: 4, NS: 3, EW: 9, Set: "H"},
				{Table: 5, NS: 2, EW: 5, Set: "I"},
				{Table: 6, NS: 10, EW: 8, Set: "K"},
			}},
			{Number: 2, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 2, Set: "B"},
				{Table: 2, NS: 7, EW: 8, Set: "E"},
				{Table: 3, NS: 1, EW: 5, Set: "G"},
				{Table: 4, NS: 4, EW: 10, Set: "I"},
				{Table: 5, NS: 3, EW: 6, Set: "J"},
				{Table: 6, NS: 11, EW: 9, Set: "A"},
			}},
			{Number: 3, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 3, Set: "C"},
				{Table: 2, NS: 8, EW: 9, Set: "F"},
				{Table: 3, NS: 2, EW: 6, Set: "H"},
				{Table: 4, NS: 5, EW: 11, Set: "J"},
				{Table: 5, NS: 4, EW: 7, Set: "K"},
				{Table: 6, NS: 1, EW: 10, Set: "B"},
			}},
			{Number: 4, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 4, Set: "D"},
				{Table: 2, NS: 9, EW: 10, Set: "G"},
				{Table: 3, NS: 3, EW: 7, Set: "I"},
				{Table: 4, NS: 6, EW: 1, Set: "K"},
				{Table: 5, NS: 5, EW: 8, Set: "A"},
				{Table: 6, NS: 2, EW: 11, Set: "C"},
			}},
			{Number: 5, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 5, Set: "E"},
				{Table: 2, NS: 10, EW: 11, Set: "H"},
				{Table: 3, NS: 4, EW: 8, Set: "J"},
				{Table: 4, NS: 7, EW: 2, Set: "A"},
				{Table: 5, NS: 6, EW: 9, Set: "B"},
				{Table: 6, NS: 3, EW: 1, Set: "D"},
			}},
			{Number: 6, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 6, Set: "F"},
				{Table: 2, NS: 11, EW: 1, Set: "I"},
				{Table: 3, NS: 5, EW: 9, Set: "K"},
				{Table: 4, NS: 8, EW: 3, Set: "B"},
				{Table: 5, NS: 7, EW: 10, Set: "C"},
				{Table: 6, NS: 4, EW: 2, Set: "E"},
			}},
			{Number: 7, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 7, Set: "G"},
				{Table: 2, NS: 1, EW: 2, Set: "J"},
				{Table: 3, NS: 6, EW: 10, Set: "A"},
				{Table: 4, NS: 9, EW: 4, Set: "C"},
				{Table: 5, NS: 8, EW: 11, Set: "D"},
				{Table: 6, NS: 5, EW: 3, Set: "F"},
			}},
			{Number: 8, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 8, Set: "H"},
				{Table: 2, NS: 2, EW: 3, Set: "K"},
				{Table: 3, NS: 7, EW: 11, Set: "B"},
				{Table: 4, NS: 10, EW: 5, Set: "D"},
				{Table: 5, NS: 9, EW: 1, Set: "E"},
				{Table: 6, NS: 6, EW: 4, Set: "G"},
			}},
			{Number: 9, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 9, Set: "I"},
				{Table: 2, NS: 3, EW: 4, Set: "A"},
				{Table: 3, NS: 8, EW: 1, Set: "C"},
				{Table: 4, NS: 11, EW: 6, Set: "E"},
				{Table: 5, NS: 10, EW: 2, Set: "F"},
				{Table: 6, NS: 7, EW: 5, Set: "H"},
			}},
			{Number: 10, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 10, Set: "J"},
				{Table: 2, NS: 4, EW: 5, Set: "B"},
				{Table: 3, NS: 9, EW: 2, Set: "D"},
				{Table: 4, NS: 1, EW: 7, Set: "F"},
				{Table: 5, NS: 11, EW: 3, Set: "G"},
				{Table: 6, NS: 8, EW: 6, Set: "I"},
			}},
			{Number: 11, Tables: []presetTable{
				{Table: 1, NS: 12, EW: 11, Set: "K"},
				{Table: 2, NS: 5, EW: 6, Set: "C"},
				{Table: 3, NS: 10, EW: 3, Set: "E"},
				{Table: 4, NS: 2, EW: 8, Set: "G"},
				{Table: 5, NS: 1, EW: 4, Set: "H"},
				{Table: 6, NS: 9, EW: 7, Set: "J"},
			}},
		},
	},
}

func presetViews() []presetView {
	views := make([]presetView, len(tournamentPresets))
	for i, preset := range tournamentPresets {
		views[i] = presetView{
			ID:      preset.ID,
			Label:   preset.Label,
			Hint:    preset.hint(),
			Summary: preset.summary(),
		}
	}
	return views
}

func findPreset(id string) (tournamentPreset, bool) {
	for _, preset := range tournamentPresets {
		if preset.ID == id {
			return preset, true
		}
	}
	return tournamentPreset{}, false
}

func applyPreset(repo tournament_domain.TournamentRepository, tournamentId string, preset tournamentPreset) error {
	teamIDs := map[int]string{}
	for _, number := range preset.teamNumbers() {
		id := uuid.NewString()
		cmd := tournament_cmd.CreateTeamCommand{
			TournamentId: tournamentId,
			TeamId:       id,
			Number:       number,
		}
		if err := cmd.Execute(repo); err != nil {
			return err
		}
		teamIDs[number] = id
	}

	tableIDs := map[int]string{}
	for _, number := range preset.tableNumbers() {
		id := uuid.NewString()
		cmd := tournament_cmd.CreateTableCommand{
			TournamentId: tournamentId,
			TableId:      id,
			Number:       number,
		}
		if err := cmd.Execute(repo); err != nil {
			return err
		}
		tableIDs[number] = id
	}

	for _, boardNo := range preset.boardNumbers() {
		cmd := tournament_cmd.CreateBoardProtocol{
			TournamentId: tournamentId,
			BoardNo:      boardNo,
			Vulnerable:   int(standardVulnerable(boardNo)),
		}
		if err := cmd.Execute(repo); err != nil {
			return err
		}
	}

	setIDs := map[string]string{}
	for _, set := range preset.Sets {
		id := uuid.NewString()
		cmd := tournament_cmd.CreateSetCommand{
			TournamentId: tournamentId,
			SetId:        id,
			Label:        set.Label,
			BoardNos:     append([]int(nil), set.Boards...),
		}
		if err := cmd.Execute(repo); err != nil {
			return err
		}
		setIDs[set.Label] = id
	}

	for _, round := range preset.Rounds {
		for _, table := range round.Tables {
			cmd := tournament_cmd.AddScheduleCommand{
				TournamentId: tournamentId,
				RoundNo:      round.Number,
				TableId:      tableIDs[table.Table],
				NsTeamId:     teamIDs[table.NS],
				EwTeamId:     teamIDs[table.EW],
				SetId:        setIDs[table.Set],
			}
			if err := cmd.Execute(repo); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p tournamentPreset) teamNumbers() []int {
	var numbers []int
	for _, round := range p.Rounds {
		for _, table := range round.Tables {
			numbers = append(numbers, table.NS, table.EW)
		}
	}
	return sortedUnique(numbers)
}

func (p tournamentPreset) tableNumbers() []int {
	var numbers []int
	for _, round := range p.Rounds {
		for _, table := range round.Tables {
			numbers = append(numbers, table.Table)
		}
	}
	return sortedUnique(numbers)
}

func (p tournamentPreset) boardNumbers() []int {
	var numbers []int
	for _, set := range p.Sets {
		numbers = append(numbers, set.Boards...)
	}
	return sortedUnique(numbers)
}

func (p tournamentPreset) hint() string {
	sets := make([]string, len(p.Sets))
	for i, set := range p.Sets {
		sets[i] = fmt.Sprintf("%s (%s)", set.Label, joinInts(set.Boards))
	}
	return fmt.Sprintf(
		"Tworzy puste drużyny %s, stoły %s, protokoły %s (partia według standardu) i komplety %s wraz z harmonogramem.",
		formatNumbers(p.teamNumbers()),
		formatNumbers(p.tableNumbers()),
		formatNumbers(p.boardNumbers()),
		strings.Join(sets, ", "),
	)
}

func (p tournamentPreset) summary() string {
	lines := make([]string, 0, len(p.Rounds)+len(p.Sets))
	for _, round := range p.Rounds {
		parts := make([]string, len(round.Tables))
		for i, table := range round.Tables {
			parts[i] = fmt.Sprintf("%d:%d-%d (%s)", table.Table, table.NS, table.EW, table.Set)
		}
		lines = append(lines, fmt.Sprintf("%d: %s", round.Number, strings.Join(parts, " ")))
	}
	for _, set := range p.Sets {
		lines = append(lines, fmt.Sprintf("%s: %s", set.Label, joinInts(set.Boards)))
	}
	return strings.Join(lines, "\n")
}

func formatNumbers(nums []int) string {
	if len(nums) == 0 {
		return ""
	}
	consecutive := true
	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[i-1]+1 {
			consecutive = false
			break
		}
	}
	if consecutive && len(nums) > 1 {
		return fmt.Sprintf("%d–%d", nums[0], nums[len(nums)-1])
	}
	return joinInts(nums)
}

func sortedUnique(values []int) []int {
	if len(values) == 0 {
		return nil
	}
	out := append([]int(nil), values...)
	sort.Ints(out)
	return slices.Compact(out)
}

// standardVulnerable is the 16-board vulnerability cycle used in duplicate bridge.
func standardVulnerable(boardNo int) tournament_domain.Vulnerable {
	cycle := []tournament_domain.Vulnerable{
		tournament_domain.None,
		tournament_domain.NS,
		tournament_domain.EW,
		tournament_domain.Both,
		tournament_domain.NS,
		tournament_domain.EW,
		tournament_domain.Both,
		tournament_domain.None,
		tournament_domain.EW,
		tournament_domain.Both,
		tournament_domain.None,
		tournament_domain.NS,
		tournament_domain.Both,
		tournament_domain.None,
		tournament_domain.NS,
		tournament_domain.EW,
	}
	index := (boardNo - 1) % len(cycle)
	if index < 0 {
		index += len(cycle)
	}
	return cycle[index]
}
