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

type Match struct {
	DealNo     int
	NsTeamId   TeamId
	EwTeamId   TeamId
	Vulnerable string
}

type PlaySessionState struct {
	TournamentId TournamentId
	Teams        []*Team
	Matches      []*Match
}

type PlaySessionProvisioned struct {
	TournamentId TournamentId
	Teams        []Team
	Matches      []Match
}

var ErrTeamsAreEmpty = errors.New("no teams in play")
var ErrMatchesAreEmpty = errors.New("no matches in play")
var ErrPlaySessionNotProvisioned = errors.New("play session is not provisioned")
var ErrTeamNotFound = errors.New("team not found")
var ErrMatchNotFound = errors.New("match not found")

type PlaySession struct {
	State  PlaySessionState
	events []any
}

func ProvisionPlaySession(id TournamentId, teams []*Team, matches []*Match) (*PlaySession, error) {
	if len(teams) == 0 {
		return nil, ErrTeamsAreEmpty
	}
	if len(matches) == 0 {
		return nil, ErrMatchesAreEmpty
	}

	teamValues := make([]Team, len(teams))
	for i, team := range teams {
		teamValues[i] = *team
	}
	matchValues := make([]Match, len(matches))
	for i, match := range matches {
		matchValues[i] = *match
	}

	return &PlaySession{
		State: PlaySessionState{TournamentId: id, Teams: teams, Matches: matches},
		events: []any{PlaySessionProvisioned{
			TournamentId: id,
			Teams:        teamValues,
			Matches:      matchValues,
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
	for _, match := range p.State.Matches {
		if match.DealNo != dealNo {
			continue
		}
		sameOrder := idutil.SameId(match.NsTeamId, teamA) && idutil.SameId(match.EwTeamId, teamB)
		oppositeOrder := idutil.SameId(match.NsTeamId, teamB) && idutil.SameId(match.EwTeamId, teamA)
		if sameOrder || oppositeOrder {
			return match, nil
		}
	}
	return nil, ErrMatchNotFound
}

func (p *PlaySession) GetEvents() []any {
	return p.events
}

func (p *PlaySession) Commit() {
	p.events = slices.Delete(p.events, 0, len(p.events))
}

type PlaySessionRepository interface {
	Save(session *PlaySession) error
	Load(id TournamentId) (*PlaySession, error)
}
