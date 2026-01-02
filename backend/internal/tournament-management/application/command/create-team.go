package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type CreateTeamCommand struct {
	TournamentId   	string
	TeamId					string
	Name				 		string
}

// Execute executes the command
func (c *CreateTeamCommand) Execute(repo domain.TournamentRepository) error {
	id := domain.TournamentId(c.TournamentId)
	teamId := domain.TeamId(c.TeamId)
	// Load the Tournament from the repository
	t, err := repo.Load(&id)
	if err != nil {
		return err
	}

	if err := t.CreateTeam(&teamId, c.Name); err != nil {
		return err
	}

	// Save the Tournament to the repository
	return repo.Save(t)
}
