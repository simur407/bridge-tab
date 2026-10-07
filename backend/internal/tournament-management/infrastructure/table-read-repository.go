package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
	"context"
	"database/sql"
)

type PostgresTableReadRepository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

func (r *PostgresTableReadRepository) FindAll(tournamentId *string) ([]domain.TableDto, error) {
	rows, err := r.Tx.QueryContext(r.Ctx, `
	SELECT id, tournament_id, number FROM tournament_management.tournament_table
	WHERE tournament_id = $1
	ORDER BY number ASC`, tournamentId)
	if err != nil {
		return nil, err
	}

	var tables []domain.TableDto
	for rows.Next() {
		var table domain.TableDto
		err := rows.Scan(&table.Id, &table.TournamentId, &table.Number)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func (r *PostgresTableReadRepository) FindById(tournamentId *string, tableId *string) (*domain.TableDto, error) {
	row := r.Tx.QueryRowContext(r.Ctx, `
	SELECT id, tournament_id, number FROM tournament_management.tournament_table
	WHERE tournament_id = $1 AND id = $2`, tournamentId, tableId)
	var table domain.TableDto
	err := row.Scan(&table.Id, &table.TournamentId, &table.Number)
	if err != nil {
		return nil, err
	}

	return &table, nil
}

func (r *PostgresTableReadRepository) FindByNumber(tournamentId *string, number int) (*domain.TableDto, error) {
	row := r.Tx.QueryRowContext(r.Ctx, `
	SELECT id, tournament_id, number FROM tournament_management.tournament_table
	WHERE tournament_id = $1 AND number = $2`, tournamentId, number)
	var table domain.TableDto
	err := row.Scan(&table.Id, &table.TournamentId, &table.Number)
	if err != nil {
		return nil, err
	}

	return &table, nil
}
