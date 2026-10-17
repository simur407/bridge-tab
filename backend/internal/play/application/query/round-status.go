package play

import "sort"

type RoundStatusBoard struct {
	Tables []int
	Rows   []RoundStatusRow
}

type RoundStatusRow struct {
	RoundNo int
	Cells   []RoundStatusCell
}

type RoundStatusCell struct {
	Seated       bool
	NsTeamNumber int
	EwTeamNumber int
	Status       string
}

func BuildRoundStatus(matches []ScheduledMatchDto) RoundStatusBoard {
	type seat struct {
		ns, ew int
		total  int
		done   int
	}
	type seatKey struct {
		round int
		table int
	}

	seats := map[seatKey]*seat{}
	tableSet := map[int]struct{}{}
	roundSet := map[int]struct{}{}
	for _, match := range matches {
		key := seatKey{round: match.RoundNo, table: match.TableNumber}
		current := seats[key]
		if current == nil {
			current = &seat{ns: match.NsTeamNumber, ew: match.EwTeamNumber}
			seats[key] = current
		}
		current.total++
		if match.Contract != "" {
			current.done++
		}
		tableSet[match.TableNumber] = struct{}{}
		roundSet[match.RoundNo] = struct{}{}
	}

	tables := sortedKeys(tableSet)
	rounds := sortedKeys(roundSet)
	rows := make([]RoundStatusRow, 0, len(rounds))
	for _, roundNo := range rounds {
		cells := make([]RoundStatusCell, 0, len(tables))
		for _, tableNo := range tables {
			current := seats[seatKey{round: roundNo, table: tableNo}]
			if current == nil {
				cells = append(cells, RoundStatusCell{})
				continue
			}
			status := "pending"
			if current.total > 0 && current.done == current.total {
				status = "completed"
			} else if current.done > 0 {
				status = "in-progress"
			}
			cells = append(cells, RoundStatusCell{
				Seated:       true,
				NsTeamNumber: current.ns,
				EwTeamNumber: current.ew,
				Status:       status,
			})
		}
		rows = append(rows, RoundStatusRow{RoundNo: roundNo, Cells: cells})
	}
	return RoundStatusBoard{Tables: tables, Rows: rows}
}

func sortedKeys(values map[int]struct{}) []int {
	keys := make([]int, 0, len(values))
	for value := range values {
		keys = append(keys, value)
	}
	sort.Ints(keys)
	return keys
}
