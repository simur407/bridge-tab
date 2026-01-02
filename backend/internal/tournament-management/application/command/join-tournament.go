package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type JoinTournamentCommand struct {
	TournamentId   string
	ContestantId 	string
}

// Execute executes the command
func (c *JoinTournamentCommand) Execute(repo domain.TournamentRepository) error {
	id := domain.TournamentId(c.TournamentId)
	contestantId := domain.ContestantId(c.ContestantId)

	// Load the Tournament from the repository
	t, err := repo.Load(&id)
	if err != nil {
		return err
	}

	if err := t.JoinTournament(&contestantId); err != nil {
		return err
	}

	// Save the Tournament to the repository
	return repo.Save(t)
}
