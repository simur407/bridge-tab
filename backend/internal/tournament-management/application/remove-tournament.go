package tournament_management

import domain "bridge-tab/internal/tournament-management/domain"

type RemoveTournamentCommand struct {
	TournamentId   string
}

// Execute executes the command
func (c *RemoveTournamentCommand) Execute(repo domain.TournamentRepository) error {
	id := domain.TournamentId(c.TournamentId)
	// Load the Tournament from the repository
	t, err := repo.Load(&id)
	if err != nil {
		return err
	}

	if err := t.Remove(); err != nil {
		return err
	}

	// Save the Tournament to the repository
	return repo.Save(t)
}
