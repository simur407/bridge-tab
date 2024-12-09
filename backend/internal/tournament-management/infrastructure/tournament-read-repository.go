package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
	"context"
	"database/sql"
)

type PostgresTournamentReadRepository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

func (r *PostgresTournamentReadRepository) FindAll() ([]domain.TournamentDto, error) {
	rows, err := r.Tx.QueryContext(r.Ctx, "SELECT id, name FROM tournament_management.tournament")
	if err != nil {
		return nil, err
	}

	var Tournaments []domain.TournamentDto
	for rows.Next() {
		var Tournament domain.TournamentDto
		err := rows.Scan(&Tournament.Id, &Tournament.Name)
		if err != nil {
			return nil, err
		}
		Tournaments = append(Tournaments, Tournament)
	}
	return Tournaments, nil
}

func (r *PostgresTournamentReadRepository) FindAllContestants(id *domain.TournamentId) ([]domain.ContestantDto, error) {
	rows, err := r.Tx.QueryContext(r.Ctx, "SELECT id FROM tournament_management.contestant WHERE tournament_id = $1", id)
	if err != nil {
		return nil, err
	}

	var contestants []domain.ContestantDto
	for rows.Next() {
		var contestant domain.ContestantDto
		err := rows.Scan(&contestant.Id)
		if err != nil {
			return nil, err
		}
		contestants = append(contestants, contestant)
	}
	return contestants, nil
}
