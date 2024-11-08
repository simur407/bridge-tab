package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type StartTurnamentCommand struct {
	TournamentId   string
}

// Execute executes the command
func (c *StartTurnamentCommand) Execute(repo domain.TournamentRepository) error {
	id := domain.TournamentId(c.TournamentId)
	// Load the turnament from the repository
	t, err := repo.Load(&id)
	if err != nil {
		return err
	}

	if err := t.Start(); err != nil {
		return err
	}

	// Save the turnament to the repository
	return repo.Save(t)
}
