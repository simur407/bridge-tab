package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type CreateBoardProtocol struct {
	TournamentId string
	BoardNo      int
	Vulnerable   int
	TeamPairs    [](struct {
		NS string
		EW string
	})
}

func (c *CreateBoardProtocol) Execute(repo domain.TournamentRepository) error {
	id := domain.TournamentId(c.TournamentId)
	vulnerable := domain.Vulnerable(c.Vulnerable)
	teamPairs := make([]domain.TeamPairs, len(c.TeamPairs))
	for i, pair := range c.TeamPairs {
		NS, EW := domain.TeamId(pair.NS), domain.TeamId(pair.EW)
		teamPairs[i] = domain.TeamPairs{
			NS: NS,
			EW: EW,
		}
	}
	// Load the Tournament from the repository
	t, err := repo.Load(&id)
	if err != nil {
		return err
	}

	if err := t.CreateBoardProtocol(c.BoardNo, vulnerable, teamPairs); err != nil {
		return err
	}

	// Save the Tournament to the repository
	return repo.Save(t)
}
