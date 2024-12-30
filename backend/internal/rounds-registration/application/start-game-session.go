package rounds_registration

import (
	domain "bridge-tab/internal/rounds-registration/domain"
	infra "bridge-tab/internal/rounds-registration/infrastructure"
	"errors"
	"slices"
)

type TeamDto struct {
	Id      string
	Players []string
}

type RoundDto struct {
	DealNo int
	NsTeam string
	EwTeam string
}

type StartGameSessionCommand struct {
	Id     string
	Teams  []TeamDto
	Rounds []RoundDto
}

func (c *StartGameSessionCommand) Execute(repository domain.GameSessionRepository) error {
	gameSessionId := domain.GameSessionId(c.Id)

	t, err := repository.Load(&gameSessionId)
	if err != nil {
		if !errors.Is(err, infra.ErrGameSessionNotFound) {
			return err
		}
	}

	if t != nil {
		return errors.New("game session already exists")
	}

	var teams []*domain.Team
	for _, team := range c.Teams {
		players := make([]domain.PlayerId, len(team.Players))
		for i, playerId := range team.Players {
			players[i] = domain.PlayerId(playerId)
		}

		teams = append(teams, &domain.Team{Id: domain.TeamId(team.Id), Players: players})
	}

	var rounds []*domain.Round
	for _, round := range c.Rounds {
		nsTeamIndex := slices.IndexFunc(teams, func(t *domain.Team) bool {
			return t.Id == domain.TeamId(round.NsTeam)
		})

		ewTeamIndex := slices.IndexFunc(teams, func(t *domain.Team) bool {
			return t.Id == domain.TeamId(round.EwTeam)
		})

		if nsTeamIndex == -1 || ewTeamIndex == -1 {
			return infra.ErrMisconfiguredGameSession
		}

		rounds = append(rounds, &domain.Round{DealNo: round.DealNo, NsTeam: &teams[nsTeamIndex].Id, EwTeam: &teams[ewTeamIndex].Id})
	}

	t, err = domain.StartGameSession(gameSessionId, teams, rounds)

	if err != nil {
		return err
	}

	return repository.Save(t)
}
