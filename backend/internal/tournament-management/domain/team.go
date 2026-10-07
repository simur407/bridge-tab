package tournament_management

import (
	"errors"
	"slices"
)

type TeamId string

type TeamState struct {
	Id           TeamId
	TournamentId TournamentId
	Name         string
	Number       int
	Members      []*Contestant
	removed      bool
}

type Team struct {
	State  TeamState
	events []any
}

// events
type ContestantJoinedTeam struct {
	TeamId       TeamId
	ContestantId ContestantId
}

type ContestantLeftTeam struct {
	TeamId       TeamId
	ContestantId ContestantId
}

// errors
var ErrTeamFull = errors.New("team is full")
var ErrTeamRemoved = errors.New("team is removed")

func CreateTeam(id TeamId, tournamentId TournamentId, name string, number int) *Team {
	return &Team{
		State:  TeamState{Id: id, TournamentId: tournamentId, Name: name, Number: number, Members: []*Contestant{}},
		events: []any{TeamCreated{TeamId: id, TournamentId: tournamentId, Name: name, Number: number}},
	}
}

func (t *Team) Remove() error {
	if t.State.removed {
		return nil
	}

	// throw all members out
	for _, c := range slices.Clone(t.State.Members) {
		if err := t.Leave(&c.Id); err != nil {
			return err
		}
	}
	t.State.removed = true
	return nil
}

func (t *Team) Join(contestant *Contestant) error {
	if t.State.removed {
		return ErrTeamRemoved
	}

	if slices.ContainsFunc(t.State.Members, func(c *Contestant) bool {
		return c.Id == contestant.Id
	}) {
		return nil
	}

	if len(t.State.Members) == 2 {
		return ErrTeamFull
	}

	contestant.Team = t
	t.State.Members = append(t.State.Members, contestant)
	t.events = append(t.events, ContestantJoinedTeam{TeamId: t.State.Id, ContestantId: contestant.Id})

	return nil
}

func (t *Team) Leave(contenstantId *ContestantId) error {
	if t.State.removed {
		return ErrTeamRemoved
	}

	memberIndex := slices.IndexFunc(t.State.Members, func(c *Contestant) bool {
		return c.Id == *contenstantId
	})

	if memberIndex != -1 {
		member := t.State.Members[memberIndex]
		member.Team = nil
		t.State.Members = slices.Delete(t.State.Members, memberIndex, memberIndex+1)

		t.events = append(t.events, ContestantLeftTeam{TeamId: t.State.Id, ContestantId: *contenstantId})
	}
	return nil
}

func (t *Team) GetEvents() []any {
	return t.events
}

func (t *Team) Commit() {
	t.events = slices.Delete(t.events, 0, len(t.events))
}

type TeamDto struct {
	Id           string
	TournamentId string
	Name         string
	Number       int
	Members      []ContestantDto
}

type TeamReadRepository interface {
	FindAll(tournamentId *string) ([]TeamDto, error)
	FindById(tournamentId *string, teamId *string) (*TeamDto, error)
	FindByName(tournamentId *string, name *string) (*TeamDto, error)
	FindByNumber(tournamentId *string, number int) (*TeamDto, error)
	FindByMemberId(tournamentId *string, memberId *string) (*TeamDto, error)
}
