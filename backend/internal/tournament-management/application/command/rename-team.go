package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type RenameTeamCommand struct {
	TournamentId string
	TeamId       string
	Name         string
}

func (c *RenameTeamCommand) Execute(repo domain.TournamentRepository) error {
	id := domain.TournamentId(c.TournamentId)
	teamId := domain.TeamId(c.TeamId)
	t, err := repo.Load(&id)
	if err != nil {
		return err
	}

	if err := t.RenameTeam(&teamId, c.Name); err != nil {
		return err
	}

	return repo.Save(t)
}
