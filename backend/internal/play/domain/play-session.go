package play

import (
	"errors"
	"slices"

	"bridge-tab/internal/idutil"
)

type PlayerId string

type Team struct {
	Id      TeamId
	Number  int
	Players []PlayerId
}

type TableScheduleBoard struct {
	DealNo     int
	Vulnerable string
}

type TableSchedule struct {
	RoundNo     int
	TableNumber int
	NsTeamId    TeamId
	EwTeamId    TeamId
	SetId       string
	SetLabel    string
	Boards      []TableScheduleBoard
}

type Match struct {
	DealNo     int
	NsTeamId   TeamId
	EwTeamId   TeamId
	Vulnerable string
}

type PlaySessionState struct {
	PlaySessionId  PlaySessionId
	Teams          []*Team
	TableSchedules []*TableSchedule
}

type PlaySessionProvisioned struct {
	PlaySessionId  PlaySessionId
	Teams          []Team
	TableSchedules []TableSchedule
}

var ErrTeamsAreEmpty = errors.New("no teams in play")
var ErrMatchesAreEmpty = errors.New("no table schedules in play")
var ErrPlaySessionNotProvisioned = errors.New("play session is not provisioned")
var ErrTeamNotFound = errors.New("team not found")
var ErrMatchNotFound = errors.New("match not found")

type PlaySession struct {
	State  PlaySessionState
	events []any
}

func ProvisionPlaySession(playSessionId PlaySessionId, teams []*Team, tableSchedules []*TableSchedule) (*PlaySession, error) {
	if len(teams) == 0 {
		return nil, ErrTeamsAreEmpty
	}
	if len(tableSchedules) == 0 {
		return nil, ErrMatchesAreEmpty
	}

	teamValues := make([]Team, len(teams))
	for i, team := range teams {
		teamValues[i] = *team
	}
	scheduleValues := make([]TableSchedule, len(tableSchedules))
	for i, schedule := range tableSchedules {
		scheduleValues[i] = *schedule
		scheduleValues[i].Boards = append([]TableScheduleBoard(nil), schedule.Boards...)
	}

	return &PlaySession{
		State: PlaySessionState{PlaySessionId: playSessionId, Teams: teams, TableSchedules: tableSchedules},
		events: []any{PlaySessionProvisioned{
			PlaySessionId:  playSessionId,
			Teams:          teamValues,
			TableSchedules: scheduleValues,
		}},
	}, nil
}

func (p *PlaySession) TeamByMember(playerId PlayerId) (*Team, error) {
	for _, team := range p.State.Teams {
		for _, player := range team.Players {
			if idutil.SameId(player, playerId) {
				return team, nil
			}
		}
	}
	return nil, ErrTeamNotFound
}

func (p *PlaySession) TeamByNumber(number int) (*Team, error) {
	for _, team := range p.State.Teams {
		if team.Number == number {
			return team, nil
		}
	}
	return nil, ErrTeamNotFound
}

func (p *PlaySession) Match(dealNo int, teamA TeamId, teamB TeamId) (*Match, error) {
	for _, schedule := range p.State.TableSchedules {
		sameOrder := idutil.SameId(schedule.NsTeamId, teamA) && idutil.SameId(schedule.EwTeamId, teamB)
		oppositeOrder := idutil.SameId(schedule.NsTeamId, teamB) && idutil.SameId(schedule.EwTeamId, teamA)
		if !sameOrder && !oppositeOrder {
			continue
		}
		for _, board := range schedule.Boards {
			if board.DealNo != dealNo {
				continue
			}
			return &Match{
				DealNo:     board.DealNo,
				NsTeamId:   schedule.NsTeamId,
				EwTeamId:   schedule.EwTeamId,
				Vulnerable: board.Vulnerable,
			}, nil
		}
	}
	return nil, ErrMatchNotFound
}

func (p *PlaySession) Matches() []*Match {
	matches := make([]*Match, 0)
	for _, schedule := range p.State.TableSchedules {
		for _, board := range schedule.Boards {
			matches = append(matches, &Match{
				DealNo:     board.DealNo,
				NsTeamId:   schedule.NsTeamId,
				EwTeamId:   schedule.EwTeamId,
				Vulnerable: board.Vulnerable,
			})
		}
	}
	return matches
}

func (p *PlaySession) GetEvents() []any {
	return p.events
}

func (p *PlaySession) Commit() {
	p.events = slices.Delete(p.events, 0, len(p.events))
}

type PlaySessionRepository interface {
	Save(session *PlaySession) error
	Load(playSessionId PlaySessionId) (*PlaySession, error)
}
