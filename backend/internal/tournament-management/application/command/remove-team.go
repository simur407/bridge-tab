package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type RemoveTeamCommand struct {
	TournamentId   	string
	TeamId					string
}

// Execute executes the command
func (c *RemoveTeamCommand) Execute(repo domain.TournamentRepository) error {
	id := domain.TournamentId(c.TournamentId)
	teamId := domain.TeamId(c.TeamId)
	// Load the Tournament from the repository
	t, err := repo.Load(&id)
	if err != nil {
		return err
	}

	if err := t.DeleteTeam(&teamId); err != nil {
		return err
	}

	// Save the Tournament to the repository
	return repo.Save(t)
}
