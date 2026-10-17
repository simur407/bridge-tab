package tournament_management

import (
	"context"
	"database/sql"

	domain "bridge-tab/internal/tournament-management/domain"
)

type PostgresScheduleReadRepository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

func (r *PostgresScheduleReadRepository) FindAll(tournamentId *string) ([]domain.ScheduleDto, error) {
	rows, err := r.Tx.QueryContext(r.Ctx, `
		SELECT schedule.round_no,
			schedule.table_id,
			tournament_table.number,
			schedule.team_ns_id,
			ns.number,
			schedule.team_ew_id,
			ew.number,
			schedule.set_id,
			tournament_set.label
		FROM tournament_management.schedule
		JOIN tournament_management.tournament_table
			ON schedule.table_id = tournament_table.id
		JOIN tournament_management.team AS ns
			ON schedule.team_ns_id = ns.id
		JOIN tournament_management.team AS ew
			ON schedule.team_ew_id = ew.id
		JOIN tournament_management.tournament_set
			ON schedule.set_id = tournament_set.id
		WHERE schedule.tournament_id = $1
		ORDER BY schedule.round_no, tournament_table.number`, tournamentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedule []domain.ScheduleDto
	for rows.Next() {
		var entry domain.ScheduleDto
		if err := rows.Scan(&entry.RoundNo, &entry.TableId, &entry.TableNumber, &entry.NsTeamId, &entry.NsNumber, &entry.EwTeamId, &entry.EwNumber, &entry.SetId, &entry.SetLabel); err != nil {
			return nil, err
		}
		schedule = append(schedule, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return schedule, nil
}
