package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type CreateBoardProtocol struct {
	TournamentId string
	BoardNo      int
	Vulnerable   int
}

func (c *CreateBoardProtocol) Execute(repo domain.TournamentRepository) error {
	tournamentId := domain.TournamentId(c.TournamentId)
	t, err := repo.Load(&tournamentId)
	if err != nil {
		return err
	}

	if err := t.CreateBoardProtocol(c.BoardNo, domain.Vulnerable(c.Vulnerable)); err != nil {
		return err
	}

	return repo.Save(t)
}
