package play

import (
	"errors"
	"slices"

	"bridge-tab/internal/idutil"
	domain "bridge-tab/internal/play/domain"
)

var ErrPlaySessionAlreadyProvisioned = errors.New("play session is already provisioned")
var ErrMisconfiguredPlay = errors.New("play is not configured correctly")

type TeamDto struct {
	Id      string
	Number  int
	Players []string
}

type MatchDto struct {
	DealNo     int
	NsTeam     string
	EwTeam     string
	Vulnerable string
}

type ProvisionPlaySessionCommand struct {
	TournamentId string
	Teams        []TeamDto
	Matches      []MatchDto
}

func (c *ProvisionPlaySessionCommand) Execute(repository domain.PlaySessionRepository) error {
	tournamentId := domain.TournamentId(c.TournamentId)

	existing, err := repository.Load(tournamentId)
	if err != nil && !errors.Is(err, domain.ErrPlaySessionNotProvisioned) {
		return err
	}
	if existing != nil {
		return ErrPlaySessionAlreadyProvisioned
	}

	teams := make([]*domain.Team, len(c.Teams))
	for i, team := range c.Teams {
		players := make([]domain.PlayerId, len(team.Players))
		for j, playerId := range team.Players {
			players[j] = domain.PlayerId(playerId)
		}
		teams[i] = &domain.Team{Id: domain.TeamId(team.Id), Number: team.Number, Players: players}
	}

	matches := make([]*domain.Match, len(c.Matches))
	for i, match := range c.Matches {
		nsIndex := slices.IndexFunc(teams, func(team *domain.Team) bool {
			return idutil.SameId(team.Id, domain.TeamId(match.NsTeam))
		})
		ewIndex := slices.IndexFunc(teams, func(team *domain.Team) bool {
			return idutil.SameId(team.Id, domain.TeamId(match.EwTeam))
		})
		if nsIndex == -1 || ewIndex == -1 {
			return ErrMisconfiguredPlay
		}
		matches[i] = &domain.Match{
			DealNo:     match.DealNo,
			NsTeamId:   teams[nsIndex].Id,
			EwTeamId:   teams[ewIndex].Id,
			Vulnerable: match.Vulnerable,
		}
	}

	session, err := domain.ProvisionPlaySession(tournamentId, teams, matches)
	if err != nil {
		return err
	}
	return repository.Save(session)
}
