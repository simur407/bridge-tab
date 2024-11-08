package tournament_management

import (
	"errors"
	"slices"
	"time"
)

type TournamentId string

type TournamentState struct {
	Id   					TournamentId
	Name 					string
	StartedAt 		*time.Time
	Teams					[]*Team
	Contestants 	[]*Contestant
	removed   		bool
}

type Tournament struct {
	State   TournamentState
	events  []any
}

// events
type TournamentCreated struct {
	TournamentId   TournamentId
	Name 					string
}

type TournamentRemoved struct {
	TournamentId   TournamentId
	Name 					string
}

type TournamentStarted struct {
	TournamentId		TournamentId
	StartedAt 		time.Time
}

type ContestantJoinedTournament struct {
	TournamentId 	TournamentId
	ContestantId	ContestantId
}

type ContestantLeftTournament struct {
	TournamentId   	TournamentId
	ContestantId 		ContestantId
}

type TeamCreated struct {
	TournamentId 	TournamentId
	TeamId				TeamId
	Name 					string
}

type TeamRemoved struct {
	TournamentId 	TournamentId
	TeamId				TeamId
}

// errors
var ErrTournamentRemoved = errors.New("Tournament removed")
var ErrTournamentNotStarted = errors.New("Tournament not started")
var ErrTournamentAlreadyStarted = errors.New("Tournament already started")
var ErrSomeTeamHasNoMembers = errors.New("one of the teams has no members")
var ErrContestantNotJoinedTournament = errors.New("contestant not joined Tournament")
var ErrNoSuchTeamInTournament = errors.New("no such team in Tournament")
var ErrTeamAlreadyExists = errors.New("team already exists")
var ErrContestantAlreadyInOtherTeam = errors.New("contestant already in other team")

func CreateTournament(id TournamentId, name string) *Tournament {
	return &Tournament{
		State:  TournamentState{Id: id, Name: name, removed: false, Teams: []*Team{}, Contestants: []*Contestant{}},
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

	for _, team := range t.State.Teams {
		if err := t.DeleteTeam(&team.State.Id); err != nil {
			return err
		}
	}
	for _, contestant := range t.State.Contestants {
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

func (t *Tournament) JoinTournament(contestantId *ContestantId) error {
	if t.State.removed {
		return ErrTournamentRemoved
	}

	if t.State.StartedAt != nil {
		return ErrTournamentAlreadyStarted
	}

	if slices.ContainsFunc(t.State.Contestants, func(c *Contestant) bool {
		return c.Id == *contestantId
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

	contestantIndex := slices.IndexFunc(t.State.Contestants, func(c *Contestant) bool{
		return c.Id == *contestantId
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

func (t *Tournament) CreateTeam(teamId *TeamId, name string) error {
	if t.State.removed {
		return ErrTournamentRemoved
	}

	if t.State.StartedAt != nil {
		return ErrTournamentAlreadyStarted
	}

	if slices.ContainsFunc(t.State.Teams, func(tt *Team) bool {
		return tt.State.Id == *teamId
	}) {
		return ErrTeamAlreadyExists
	}

	team := CreateTeam(*teamId, t.State.Id, name)

	t.State.Teams = append(t.State.Teams, team)
	t.events = append(t.events, TeamCreated{TournamentId: t.State.Id, TeamId: *teamId, Name: name})
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
		return tt.State.Id == *teamId
	})
	teamToRemove := t.State.Teams[teamIndex]


	if teamToRemove != nil {
		if err := teamToRemove.Remove() ; err != nil {
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
		return c.Id == *contestantId
	})
	if contestantIndex == -1 {
		return ErrContestantNotJoinedTournament
	}
	contestant := t.State.Contestants[contestantIndex]

	contestantHasTeam := contestant.Team != nil
	itsDifferentTeam := contestant.Team != nil && contestant.Team.State.Id != *teamId

	if contestantHasTeam && itsDifferentTeam {
		return ErrContestantAlreadyInOtherTeam
	}

	teamIndex := slices.IndexFunc(t.State.Teams, func(tt *Team) bool {
		return tt.State.Id == *teamId
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
		return c.Id == *contestantId
	})
	if contestantIndex == -1 {
		return ErrContestantNotJoinedTournament
	}

	contestant := t.State.Contestants[contestantIndex]

	if (contestant.Team != nil) {
		team := contestant.Team

		if err := team.Leave(contestantId); err != nil {
			return err
		}
	
		t.events = append(t.events, team.GetEvents()...)
		team.Commit()
	}
	
	return nil
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
	Id   string
	Name string
}

type TournamentReadRepository interface {
	FindAll() ([]TournamentDto, error)
	FindAllContestants(id *TournamentId) ([]ContestantDto, error)
}

// type InMemoryTournamentRepository struct {
// 	Tournaments []*Tournament
// }

// func (r *InMemoryTournamentRepository) Load(Id TournamentId) (*Tournament, error) {
// 	for _, t := range r.Tournaments {
// 		if t.State.Id == Id {
// 			return t, nil
// 		}
// 	}
// 	return nil, nil
// }

// func (r *InMemoryTournamentRepository) Save(t *Tournament) error {
// 	for i, Tournament := range r.Tournaments {
// 		if Tournament.State.Id == t.State.Id {
// 			r.Tournaments[i] = t
// 			return nil
// 		}
// 	}
// 	r.Tournaments = append(r.Tournaments, t)
// 	return nil
// }

// func (r *InMemoryTournamentRepository) FindAll() ([]TournamentDto, error) {
// 	var Tournaments []TournamentDto
// 	for _, t := range r.Tournaments {
// 		Tournaments = append(Tournaments, TournamentDto{Id: string(t.State.Id), Name: t.State.Name})
// 	}
// 	return Tournaments, nil
// }
