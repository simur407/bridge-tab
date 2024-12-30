package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type JoinTeamCommand struct {
	TournamentId   string
	TeamId				string
	ContestantId 	string
}

// Execute executes the command
func (c *JoinTeamCommand) Execute(repo domain.TournamentRepository) error {
	id := domain.TournamentId(c.TournamentId)
	teamId := domain.TeamId(c.TeamId)
	contestantId := domain.ContestantId(c.ContestantId)
	// Load the Tournament from the repository
	t, err := repo.Load(&id)
	if err != nil {
		return err
	}

	if err := t.JoinTeam(&teamId, &contestantId); err != nil {
		return err
	}

	// Save the Tournament to the repository
	return repo.Save(t)
}
