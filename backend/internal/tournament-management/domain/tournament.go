package tournament_management

import (
	"bridge-tab/internal/idutil"
	"errors"
	"slices"
	"strings"
	"time"
)

type TournamentId string

type TournamentState struct {
	Id             TournamentId
	Name           string
	StartedAt      *time.Time
	FinishedAt     *time.Time
	Teams          []*Team
	Tables         []*Table
	Contestants    []*Contestant
	BoardProtocols []*BoardProtocol
	Sets           []*Set
	removed        bool
}

type Tournament struct {
	State  TournamentState
	events []any
}

// events
type TournamentCreated struct {
	TournamentId TournamentId
	Name         string
}

type TournamentRemoved struct {
	TournamentId TournamentId
	Name         string
}

type TournamentStarted struct {
	TournamentId TournamentId
	StartedAt    time.Time
}

type TournamentFinished struct {
	TournamentId TournamentId
	FinishedAt   time.Time
}

type SetCreated struct {
	TournamentId TournamentId
	SetId        SetId
	Label        string
	BoardNos     []int
	TeamPairs    []TeamPairs
}

type SetRemoved struct {
	TournamentId TournamentId
	SetId        SetId
}

type ContestantJoinedTournament struct {
	TournamentId TournamentId
	ContestantId ContestantId
}

type ContestantLeftTournament struct {
	TournamentId TournamentId
	ContestantId ContestantId
}

type TeamCreated struct {
	TournamentId TournamentId
	TeamId       TeamId
	Name         string
	Number       int
}

type TeamRemoved struct {
	TournamentId TournamentId
	TeamId       TeamId
}

type TableCreated struct {
	TournamentId TournamentId
	TableId      TableId
	Number       int
}

type TableRemoved struct {
	TournamentId TournamentId
	TableId      TableId
}

type BoardProtocolCreated struct {
	TournamentId TournamentId
	BoardNo      int
	Vulnerable   Vulnerable
	TeamPairs    []TeamPairs
}

type BoardProtocolRemoved struct {
	TournamentId TournamentId
	BoardNo      int
}

// errors
var ErrTournamentRemoved = errors.New("Tournament removed")
var ErrTournamentNotStarted = errors.New("Tournament not started")
var ErrTournamentAlreadyStarted = errors.New("Tournament already started")
var ErrSomeTeamHasNoMembers = errors.New("one of the teams has no members")
var ErrContestantNotJoinedTournament = errors.New("contestant not joined Tournament")
var ErrNoSuchTeamInTournament = errors.New("no such team in Tournament")
var ErrTeamAlreadyExists = errors.New("team already exists")
var ErrTeamNumberAlreadyExists = errors.New("team number already exists")
var ErrTeamNameAlreadyExists = errors.New("team name already exists")
var ErrInvalidTeamNumber = errors.New("team number must be positive")
var ErrInvalidTeamName = errors.New("team name must not be empty")
var ErrContestantAlreadyInOtherTeam = errors.New("contestant already in other team")
var ErrTableAlreadyExists = errors.New("table already exists")
var ErrTableNumberAlreadyExists = errors.New("table number already exists")
var ErrInvalidTableNumber = errors.New("table number must be positive")
var ErrNoSuchTableInTournament = errors.New("no such table in Tournament")
var ErrTableReferencedByBoardProtocol = errors.New("table is referenced by a board protocol")
var ErrBoardProtocolAlreadyExists = errors.New("board protocol already exists")
var ErrNoSuchBoardProtocol = errors.New("no such board protocol")
var ErrBoardProtocolHasTheSameTeamMultipleTimes = errors.New("board protocol has the same team multiple times")
var ErrBoardProtocolHasTheSameTableMultipleTimes = errors.New("board protocol has the same table multiple times")
var ErrTournamentAlreadyFinished = errors.New("tournament already finished")
var ErrSetAlreadyExists = errors.New("set with this label already exists")
var ErrSetIdAlreadyExists = errors.New("set already exists")
var ErrNoSuchSet = errors.New("no such set")
var ErrSetBoardsEmpty = errors.New("set must contain at least one board")
var ErrSetBoardMissingProtocol = errors.New("set board has no board protocol")
var ErrSetBoardAlreadyInSet = errors.New("board already belongs to a set")
var ErrSetTeamPairsMismatch = errors.New("board protocols in a set must share the same team pairs")
var ErrBoardProtocolBelongsToSet = errors.New("board protocol belongs to a set")
var ErrInvalidSetLabel = errors.New("set label must not be empty")

func CreateTournament(id TournamentId, name string) *Tournament {
	return &Tournament{
		State: TournamentState{
			Id:          id,
			Name:        name,
			removed:     false,
			Teams:       []*Team{},
			Tables:      []*Table{},
			Contestants: []*Contestant{},
			Sets:        []*Set{},
		},
		events: []any{TournamentCreated{TournamentId: id, Name: name}},
	}
}

func (t *Tournament) Remove() error {
	if t.State.removed {
		return nil
	}

	if t.State.StartedAt != nil {
		return ErrTournamentAlreadyStarted
	}

	for _, set := range slices.Clone(t.State.Sets) {
		if err := t.RemoveSet(set.Id); err != nil {
			return err
		}
	}
	for _, boardProtocol := range slices.Clone(t.State.BoardProtocols) {
		if err := t.RemoveBoardProtocol(boardProtocol.BoardNo); err != nil {
			return err
		}
	}
	for _, team := range slices.Clone(t.State.Teams) {
		if err := t.DeleteTeam(&team.State.Id); err != nil {
			return err
		}
	}
	for _, table := range slices.Clone(t.State.Tables) {
		if err := t.DeleteTable(&table.State.Id); err != nil {
			return err
		}
	}
	for _, contestant := range slices.Clone(t.State.Contestants) {
		if err := t.LeaveTournament(&contestant.Id); err != nil {
			return err
		}
	}
	t.State.removed = true
	t.events = append(t.events, TournamentRemoved{TournamentId: t.State.Id})
	return nil
}

func (t *Tournament) Start() error {
	if t.State.removed {
		return ErrTournamentRemoved
	}

	if t.State.StartedAt != nil {
		return nil
	}

	// TODO: check if all teams have at least one contestant
	for _, team := range t.State.Teams {
		if len(team.State.Members) == 0 {
			return ErrSomeTeamHasNoMembers
		}
	}

	now := time.Now()
	t.State.StartedAt = &now
	t.events = append(t.events, TournamentStarted{TournamentId: t.State.Id, StartedAt: now})
	return nil
}

func (t *Tournament) Finish() error {
	if t.State.removed {
		return ErrTournamentRemoved
	}

	if t.State.StartedAt == nil {
		return ErrTournamentNotStarted
	}

	if t.State.FinishedAt != nil {
		return ErrTournamentAlreadyFinished
	}

	now := time.Now()
	t.State.FinishedAt = &now
	t.events = append(t.events, TournamentFinished{TournamentId: t.State.Id, FinishedAt: now})
	return nil
}

func (t *Tournament) JoinTournament(contestantId *ContestantId) error {
	if t.State.removed {
		return ErrTournamentRemoved
	}

	if t.State.StartedAt != nil {
		return ErrTournamentAlreadyStarted
	}

	if slices.ContainsFunc(t.State.Contestants, func(c *Contestant) bool {
		return idutil.SameId(c.Id, *contestantId)
	}) {
		return nil
	}

	t.State.Contestants = append(t.State.Contestants, &Contestant{Id: *contestantId})
	t.events = append(t.events, ContestantJoinedTournament{TournamentId: t.State.Id, ContestantId: *contestantId})
	return nil
}

func (t *Tournament) LeaveTournament(contestantId *ContestantId) error {
	if t.State.removed {
		return ErrTournamentRemoved
	}

	if t.State.StartedAt != nil {
		return ErrTournamentAlreadyStarted
	}

	contestantIndex := slices.IndexFunc(t.State.Contestants, func(c *Contestant) bool {
		return idutil.SameId(c.Id, *contestantId)
	})
	if contestantIndex == -1 {
		return nil
	}

	contestant := t.State.Contestants[contestantIndex]

	if contestant.Team != nil {
		team := contestant.Team
		if err := team.Leave(contestantId); err != nil {
			return err
		}
		t.events = append(t.events, team.GetEvents()...)
		team.Commit()
	}

	// Leave Tournament
	t.State.Contestants = slices.Delete(t.State.Contestants, contestantIndex, contestantIndex+1)
	t.events = append(t.events, ContestantLeftTournament{TournamentId: t.State.Id, ContestantId: *contestantId})
	return nil
}

func (t *Tournament) CreateTeam(teamId *TeamId, name string, number int) error {
	if t.State.removed {
		return ErrTournamentRemoved
	}

	if t.State.StartedAt != nil {
		return ErrTournamentAlreadyStarted
	}

	if slices.ContainsFunc(t.State.Teams, func(tt *Team) bool {
		return idutil.SameId(tt.State.Id, *teamId)
	}) {
		return ErrTeamAlreadyExists
	}

	if number == 0 {
		number = t.nextTeamNumber()
	} else if number < 1 {
		return ErrInvalidTeamNumber
	}

	if slices.ContainsFunc(t.State.Teams, func(tt *Team) bool {
		return tt.State.Number == number
	}) {
		return ErrTeamNumberAlreadyExists
	}

	team := CreateTeam(*teamId, t.State.Id, name, number)

	t.State.Teams = append(t.State.Teams, team)
	t.events = append(t.events, team.GetEvents()...)
	team.Commit()
	return nil
}

func (t *Tournament) RenameTeam(teamId *TeamId, name string) error {
	if t.State.removed {
		return ErrTournamentRemoved
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return ErrInvalidTeamName
	}

	teamIndex := slices.IndexFunc(t.State.Teams, func(tt *Team) bool {
		return idutil.SameId(tt.State.Id, *teamId)
	})
	if teamIndex == -1 {
		return ErrNoSuchTeamInTournament
	}

	for _, other := range t.State.Teams {
		if idutil.SameId(other.State.Id, *teamId) {
			continue
		}
		if other.State.Name == name {
			return ErrTeamNameAlreadyExists
		}
	}

	team := t.State.Teams[teamIndex]
	if err := team.Rename(name); err != nil {
		return err
	}
	t.events = append(t.events, team.GetEvents()...)
	team.Commit()
	return nil
}

func (t *Tournament) nextTeamNumber() int {
	next := 1
	for _, team := range t.State.Teams {
		if team.State.Number >= next {
			next = team.State.Number + 1
		}
	}
	return next
}

func (t *Tournament) CreateTable(tableId *TableId, number int) error {
	if t.State.removed {
		return ErrTournamentRemoved
	}

	if t.State.StartedAt != nil {
		return ErrTournamentAlreadyStarted
	}

	if slices.ContainsFunc(t.State.Tables, func(tt *Table) bool {
		return idutil.SameId(tt.State.Id, *tableId)
	}) {
		return ErrTableAlreadyExists
	}

	if number == 0 {
		number = t.nextTableNumber()
	} else if number < 1 {
		return ErrInvalidTableNumber
	}

	if slices.ContainsFunc(t.State.Tables, func(tt *Table) bool {
		return tt.State.Number == number
	}) {
		return ErrTableNumberAlreadyExists
	}

	table := CreateTable(*tableId, t.State.Id, number)

	t.State.Tables = append(t.State.Tables, table)
	t.events = append(t.events, TableCreated{TournamentId: t.State.Id, TableId: *tableId, Number: number})
	return nil
}

func (t *Tournament) nextTableNumber() int {
	next := 1
	for _, table := range t.State.Tables {
		if table.State.Number >= next {
			next = table.State.Number + 1
		}
	}
	return next
}

func (t *Tournament) DeleteTable(tableId *TableId) error {
	if t.State.removed {
		return ErrTournamentRemoved
	}

	if t.State.StartedAt != nil {
		return ErrTournamentAlreadyStarted
	}

	tableIndex := slices.IndexFunc(t.State.Tables, func(tt *Table) bool {
		return idutil.SameId(tt.State.Id, *tableId)
	})
	if tableIndex == -1 {
		return ErrNoSuchTableInTournament
	}

	for _, boardProtocol := range t.State.BoardProtocols {
		for _, pair := range boardProtocol.TeamPairs {
			if pair.Table != nil && idutil.SameId(*pair.Table, *tableId) {
				return ErrTableReferencedByBoardProtocol
			}
		}
	}

	t.State.Tables = slices.Delete(t.State.Tables, tableIndex, tableIndex+1)
	t.events = append(t.events, TableRemoved{TournamentId: t.State.Id, TableId: *tableId})
	return nil
}

func (t *Tournament) DeleteTeam(teamId *TeamId) error {
	if t.State.removed {
		return ErrTournamentRemoved
	}

	if t.State.StartedAt != nil {
		return ErrTournamentAlreadyStarted
	}

	teamIndex := slices.IndexFunc(t.State.Teams, func(tt *Team) bool {
		return idutil.SameId(tt.State.Id, *teamId)
	})
	teamToRemove := t.State.Teams[teamIndex]

	if teamToRemove != nil {
		if err := teamToRemove.Remove(); err != nil {
			return err
		}

		t.State.Teams = slices.Delete(t.State.Teams, teamIndex, teamIndex+1)
		t.events = append(t.events, teamToRemove.GetEvents()...)
		teamToRemove.Commit()
		t.events = append(t.events, TeamRemoved{TournamentId: t.State.Id, TeamId: teamToRemove.State.Id})
	}

	return nil
}

func (t *Tournament) JoinTeam(teamId *TeamId, contestantId *ContestantId) error {
	if t.State.removed {
		return ErrTournamentRemoved
	}

	if t.State.StartedAt != nil {
		return ErrTournamentAlreadyStarted
	}

	contestantIndex := slices.IndexFunc(t.State.Contestants, func(c *Contestant) bool {
		return idutil.SameId(c.Id, *contestantId)
	})
	if contestantIndex == -1 {
		return ErrContestantNotJoinedTournament
	}
	contestant := t.State.Contestants[contestantIndex]

	contestantHasTeam := contestant.Team != nil
	itsDifferentTeam := contestant.Team != nil && !idutil.SameId(contestant.Team.State.Id, *teamId)

	if contestantHasTeam && itsDifferentTeam {
		return ErrContestantAlreadyInOtherTeam
	}

	teamIndex := slices.IndexFunc(t.State.Teams, func(tt *Team) bool {
		return idutil.SameId(tt.State.Id, *teamId)
	})
	if teamIndex == -1 {
		return ErrNoSuchTeamInTournament
	}
	team := t.State.Teams[teamIndex]

	if err := team.Join(contestant); err != nil {
		return err
	}

	t.events = append(t.events, team.GetEvents()...)
	team.Commit()
	return nil
}

func (t *Tournament) LeaveTeam(teamId *TeamId, contestantId *ContestantId) error {
	if t.State.removed {
		return ErrTournamentRemoved
	}

	if t.State.StartedAt != nil {
		return ErrTournamentAlreadyStarted
	}

	contestantIndex := slices.IndexFunc(t.State.Contestants, func(c *Contestant) bool {
		return idutil.SameId(c.Id, *contestantId)
	})
	if contestantIndex == -1 {
		return ErrContestantNotJoinedTournament
	}

	contestant := t.State.Contestants[contestantIndex]

	if contestant.Team != nil {
		team := contestant.Team

		if err := team.Leave(contestantId); err != nil {
			return err
		}

		t.events = append(t.events, team.GetEvents()...)
		team.Commit()
	}

	return nil
}

func (t *Tournament) CreateBoardProtocol(boardNo int, vulnerable Vulnerable, teamPairs []TeamPairs) error {
	if t.State.removed {
		return ErrTournamentRemoved
	}

	if t.State.StartedAt != nil {
		return ErrTournamentAlreadyStarted
	}

	if slices.ContainsFunc(t.State.BoardProtocols, func(bp *BoardProtocol) bool {
		return bp.BoardNo == boardNo
	}) {
		return ErrBoardProtocolAlreadyExists
	}

	// teams cannot play against themselves and cannot play twice;
	// optional tables must exist and cannot appear twice on one board
	for i := 0; i < len(teamPairs); i++ {
		if teamPairs[i].Table != nil {
			if !slices.ContainsFunc(t.State.Tables, func(tt *Table) bool {
				return idutil.SameId(tt.State.Id, *teamPairs[i].Table)
			}) {
				return ErrNoSuchTableInTournament
			}
		}
		if idutil.SameId(teamPairs[i].NS, teamPairs[i].EW) {
			return ErrBoardProtocolHasTheSameTeamMultipleTimes
		}
		for j := i + 1; j < len(teamPairs); j++ {
			if idutil.SameId(teamPairs[i].NS, teamPairs[j].NS) || idutil.SameId(teamPairs[i].NS, teamPairs[j].EW) ||
				idutil.SameId(teamPairs[i].EW, teamPairs[j].EW) || idutil.SameId(teamPairs[i].EW, teamPairs[j].NS) {
				return ErrBoardProtocolHasTheSameTeamMultipleTimes
			}
			if teamPairs[i].Table != nil && teamPairs[j].Table != nil && idutil.SameId(*teamPairs[i].Table, *teamPairs[j].Table) {
				return ErrBoardProtocolHasTheSameTableMultipleTimes
			}
		}
	}

	boardProtocol := CreateBoardProtocol(t.State.Id, boardNo, vulnerable, teamPairs)
	t.State.BoardProtocols = append(t.State.BoardProtocols, boardProtocol)
	t.events = append(t.events, BoardProtocolCreated{TournamentId: t.State.Id, BoardNo: boardNo, Vulnerable: vulnerable, TeamPairs: teamPairs})
	return nil
}

func (t *Tournament) RemoveBoardProtocol(boardNo int) error {
	if t.State.removed {
		return ErrTournamentRemoved
	}

	if t.State.StartedAt != nil {
		return ErrTournamentAlreadyStarted
	}

	protocolIndex := slices.IndexFunc(t.State.BoardProtocols, func(bp *BoardProtocol) bool {
		return bp.BoardNo == boardNo
	})
	if protocolIndex == -1 {
		return ErrNoSuchBoardProtocol
	}

	if t.State.BoardProtocols[protocolIndex].SetId != nil {
		return ErrBoardProtocolBelongsToSet
	}

	t.State.BoardProtocols = slices.Delete(t.State.BoardProtocols, protocolIndex, protocolIndex+1)
	t.events = append(t.events, BoardProtocolRemoved{TournamentId: t.State.Id, BoardNo: boardNo})
	return nil
}

func (t *Tournament) CreateSet(setId SetId, label string, boardNos []int) error {
	if t.State.removed {
		return ErrTournamentRemoved
	}

	if t.State.StartedAt != nil {
		return ErrTournamentAlreadyStarted
	}

	if label == "" {
		return ErrInvalidSetLabel
	}

	if len(boardNos) == 0 {
		return ErrSetBoardsEmpty
	}

	if slices.ContainsFunc(t.State.Sets, func(s *Set) bool {
		return idutil.SameId(s.Id, setId)
	}) {
		return ErrSetIdAlreadyExists
	}

	if slices.ContainsFunc(t.State.Sets, func(s *Set) bool {
		return s.Label == label
	}) {
		return ErrSetAlreadyExists
	}

	protocols := make([]*BoardProtocol, 0, len(boardNos))
	for _, boardNo := range boardNos {
		protocolIndex := slices.IndexFunc(t.State.BoardProtocols, func(bp *BoardProtocol) bool {
			return bp.BoardNo == boardNo
		})
		if protocolIndex == -1 {
			return ErrSetBoardMissingProtocol
		}
		protocol := t.State.BoardProtocols[protocolIndex]
		if protocol.SetId != nil {
			return ErrSetBoardAlreadyInSet
		}
		protocols = append(protocols, protocol)
	}

	referencePairs := protocols[0].TeamPairs
	for _, protocol := range protocols[1:] {
		if !teamPairsEqual(referencePairs, protocol.TeamPairs) {
			return ErrSetTeamPairsMismatch
		}
	}

	set := &Set{
		Id:        setId,
		Label:     label,
		BoardNos:  append([]int(nil), boardNos...),
		TeamPairs: append([]TeamPairs(nil), referencePairs...),
	}
	t.State.Sets = append(t.State.Sets, set)
	for _, protocol := range protocols {
		id := setId
		protocol.SetId = &id
	}
	t.events = append(t.events, SetCreated{
		TournamentId: t.State.Id,
		SetId:        setId,
		Label:        label,
		BoardNos:     append([]int(nil), boardNos...),
		TeamPairs:    append([]TeamPairs(nil), referencePairs...),
	})
	return nil
}

func (t *Tournament) RemoveSet(setId SetId) error {
	if t.State.removed {
		return ErrTournamentRemoved
	}

	if t.State.StartedAt != nil {
		return ErrTournamentAlreadyStarted
	}

	setIndex := slices.IndexFunc(t.State.Sets, func(s *Set) bool {
		return idutil.SameId(s.Id, setId)
	})
	if setIndex == -1 {
		return ErrNoSuchSet
	}

	for _, protocol := range t.State.BoardProtocols {
		if protocol.SetId != nil && idutil.SameId(*protocol.SetId, setId) {
			protocol.SetId = nil
		}
	}

	t.State.Sets = slices.Delete(t.State.Sets, setIndex, setIndex+1)
	t.events = append(t.events, SetRemoved{TournamentId: t.State.Id, SetId: setId})
	return nil
}

func teamPairsEqual(a, b []TeamPairs) bool {
	if len(a) != len(b) {
		return false
	}
	used := make([]bool, len(b))
	for _, left := range a {
		matched := false
		for j, right := range b {
			if used[j] {
				continue
			}
			if idutil.SameId(left.NS, right.NS) && idutil.SameId(left.EW, right.EW) && tableIdEqual(left.Table, right.Table) {
				used[j] = true
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return true
}

func tableIdEqual(a, b *TableId) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return idutil.SameId(*a, *b)
}

// GetEvents returns the events of a Tournament
func (t *Tournament) GetEvents() []any {
	return t.events
}

func (t *Tournament) Commit() {
	t.events = slices.Delete(t.events, 0, len(t.events))
}

type TournamentRepository interface {
	Load(Id *TournamentId) (*Tournament, error)
	Save(t *Tournament) error
}

type TournamentDto struct {
	Id         string
	Name       string
	StartedAt  string
	FinishedAt string
}

type TournamentReadRepository interface {
	FindById(id string) (*TournamentDto, error)
	FindAll() ([]TournamentDto, error)
	FindAllContestants(id *TournamentId) ([]ContestantDto, error)
}
