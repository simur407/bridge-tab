package tournament_management

import (
	"errors"

	domain "bridge-tab/internal/tournament-management/domain"
	infra "bridge-tab/internal/tournament-management/infrastructure"
)

type CreateTournamentCommand struct {
	TournamentId   	string
	Name 						string
}

func (c *CreateTournamentCommand) Execute(repo domain.TournamentRepository) error {
	// Validate input
	if err := validate(c); err != nil {
		return err
	}

	tournamentId := domain.TournamentId(c.TournamentId)

	t, err := repo.Load(&tournamentId)
	if err != nil {
		if !errors.Is(err, infra.ErrTournamentNotFound) {
			return err
		}
	}

	if t != nil {
		return errors.New("tournament already exists")
	}

	t = domain.CreateTournament(tournamentId, c.Name)
	return repo.Save(t)
}

func validate(c *CreateTournamentCommand) error {
	if c.TournamentId == "" || c.Name == "" {
		return errors.New("invalid input")
	}
	return nil
}
