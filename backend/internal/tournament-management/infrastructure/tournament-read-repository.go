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
	rows, err := r.Tx.QueryContext(r.Ctx, "SELECT id, name, started_at FROM tournament_management.tournament")
	if err != nil {
		return nil, err
	}

	var Tournaments []domain.TournamentDto
	for rows.Next() {
		var Tournament domain.TournamentDto
		var StartedAt sql.NullString
		err := rows.Scan(&Tournament.Id, &Tournament.Name, &StartedAt)
		if err != nil {
			return nil, err
		}

		if StartedAt.Valid {
			Tournament.StartedAt = StartedAt.String
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

func (r *PostgresTournamentReadRepository) FindById(id string) (*domain.TournamentDto, error) {
	row := r.Tx.QueryRowContext(r.Ctx, "SELECT id, name, started_at FROM tournament_management.tournament WHERE id = $1", id)

	var Tournament domain.TournamentDto
	var StartedAt sql.NullString
	err := row.Scan(&Tournament.Id, &Tournament.Name, &StartedAt)
	if err != nil {
		return nil, err
	}
	if Tournament.Id == "" {
		return nil, nil
	}

	if StartedAt.Valid {
		Tournament.StartedAt = StartedAt.String
	}

	return &Tournament, nil
}
