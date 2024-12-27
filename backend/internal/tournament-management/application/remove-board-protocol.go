package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type RemoveBoardProtocol struct {
	tournamentId string
	boardNo      int
}

func (c *RemoveBoardProtocol) Execute(repo domain.TournamentRepository) error {
	id := domain.TournamentId(c.tournamentId)

	// Load the Tournament from the repository
	t, err := repo.Load(&id)
	if err != nil {
		return err
	}

	if err := t.RemoveBoardProtocol(c.boardNo); err != nil {
		return err
	}

	// Save the Tournament to the repository
	return repo.Save(t)
}
