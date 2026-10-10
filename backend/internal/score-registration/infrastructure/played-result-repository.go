package score_registration

import (
	"context"
	"database/sql"

	domain "bridge-tab/internal/score-registration/domain"
)

type PostgresPlayedResultRepository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

func (r *PostgresPlayedResultRepository) Replace(id domain.TournamentScoreId, results []domain.PlayedResult) error {
	_, err := r.Tx.ExecContext(r.Ctx, `
		DELETE FROM score_registration.played_result WHERE tournament_id = $1`, id)
	if err != nil {
		return err
	}

	for _, result := range results {
		_, err = r.Tx.ExecContext(r.Ctx, `
			INSERT INTO score_registration.played_result (
				tournament_id, deal_no, ns_team_id, ew_team_id,
				contract, tricks, declarer, ns_vulnerable, ew_vulnerable
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			id, result.DealNo, result.NsTeamId, result.EwTeamId,
			result.Contract, result.Tricks, result.Declarer, result.NsVulnerable, result.EwVulnerable,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresPlayedResultRepository) Find(id domain.TournamentScoreId) ([]domain.PlayedResult, error) {
	rows, err := r.Tx.QueryContext(r.Ctx, `
		SELECT deal_no, ns_team_id, ew_team_id, contract, tricks, declarer, ns_vulnerable, ew_vulnerable
		FROM score_registration.played_result
		WHERE tournament_id = $1
		ORDER BY deal_no, ns_team_id, ew_team_id`, id)
	if err != nil {
		return nil, err
	}

	var results []domain.PlayedResult
	for rows.Next() {
		var result domain.PlayedResult
		if err := rows.Scan(
			&result.DealNo, &result.NsTeamId, &result.EwTeamId,
			&result.Contract, &result.Tricks, &result.Declarer,
			&result.NsVulnerable, &result.EwVulnerable,
		); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
