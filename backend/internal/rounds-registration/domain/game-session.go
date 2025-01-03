package rounds_registration

import (
	"errors"
	"slices"
)

type GameSessionId string

type GameSessionState struct {
	Id     GameSessionId
	Rounds []*Round
	Teams  []*Team
}

type GameSessionStarted struct {
	GameSessionId GameSessionId
	Teams         []Team
	Rounds        []Round
}

type RoundPlayed struct {
	GameSessionId GameSessionId
	DealNo        int
	NsTeamId      TeamId
	EwTeamId      TeamId
	Contract      string
	Tricks        int
	Declarer      string
	OpeningLead   string
}

var ErrRoundNotFound = errors.New("round not found")
var ErrRoundAlreadyPlayed = errors.New("round already played")
var ErrTeamsAreEmpty = errors.New("no teams in game session")
var ErrRoundsAreEmpty = errors.New("no rounds in game session")

type GameSession struct {
	State  GameSessionState
	events []any
}

func StartGameSession(id GameSessionId, teams []*Team, rounds []*Round) (*GameSession, error) {
	if len(teams) == 0 {
		return nil, ErrTeamsAreEmpty
	}

	if len(rounds) == 0 {
		return nil, ErrRoundsAreEmpty
	}

	teamValues := make([]Team, len(teams))
	for i, team := range teams {
		teamValues[i] = *team
	}
	roundValues := make([]Round, len(rounds))
	for i, round := range rounds {
		roundValues[i] = *round
	}
	return &GameSession{
		State:  GameSessionState{Id: id, Teams: teams, Rounds: rounds},
		events: []any{GameSessionStarted{GameSessionId: id, Teams: teamValues, Rounds: roundValues}},
	}, nil
}

func (g *GameSession) AddRoundScore(dealNo int, teamId TeamId, versusTeamId TeamId, contract string, tricks int, declarer string, openingLead string) error {
	index := slices.IndexFunc(g.State.Rounds, func(r *Round) bool {
		return r.DealNo == dealNo &&
			(*r.NsTeam == teamId || *r.EwTeam == teamId) &&
			(*r.NsTeam == versusTeamId || *r.EwTeam == versusTeamId)
	})

	if index == -1 {
		return ErrRoundNotFound
	}

	round := g.State.Rounds[index]

	if round.IsPlayed() {
		return ErrRoundAlreadyPlayed
	}

	round.Contract = contract
	round.Tricks = tricks
	round.Declarer = declarer
	round.OpeningLead = openingLead

	g.events = append(g.events, RoundPlayed{
		GameSessionId: g.State.Id,
		DealNo:        dealNo,
		NsTeamId:      *round.NsTeam,
		EwTeamId:      *round.EwTeam,
		Contract:      contract,
		Tricks:        tricks,
		Declarer:      declarer,
		OpeningLead:   openingLead,
	})

	return nil
}

func (g *GameSession) GetEvents() []any {
	return g.events
}

type GameSessionRepository interface {
	Save(gameSession *GameSession) error
	Load(id *GameSessionId) (*GameSession, error)
}

type RoundDto struct {
	DealNo     int
	NsTeamName string
	EwTeamName string
}

type GameSessionReadRepository interface {
	FindRound(gameSessionId *string, dealNo int, playerTeamId string, versusTeamId string) (*RoundDto, error)
}
