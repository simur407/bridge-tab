package play

import (
	"context"
	"database/sql"
	"errors"

	domain "bridge-tab/internal/play/domain"
)

type PostgresBoardResultReadRepository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

func (r *PostgresBoardResultReadRepository) FindAll(tournamentId string) ([]domain.BoardResultDto, error) {
	rows, err := r.Tx.QueryContext(r.Ctx, `
		SELECT board_result.deal_no, board_result.ns_team_id, board_result.ew_team_id,
			ns_team.number, ew_team.number,
			board_result.contract, board_result.tricks, board_result.declarer, board_result.opening_lead
		FROM play.board_result AS board_result
		LEFT JOIN play.team AS ns_team
			ON board_result.ns_team_id = ns_team.team_id
			AND board_result.tournament_id = ns_team.tournament_id
		LEFT JOIN play.team AS ew_team
			ON board_result.ew_team_id = ew_team.team_id
			AND board_result.tournament_id = ew_team.tournament_id
		WHERE board_result.tournament_id = $1
		ORDER BY board_result.deal_no, ns_team.number, ew_team.number`, tournamentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.BoardResultDto
	for rows.Next() {
		var result domain.BoardResultDto
		if err := rows.Scan(
			&result.DealNo, &result.NsTeamId, &result.EwTeamId,
			&result.NsTeamNumber, &result.EwTeamNumber,
			&result.Contract, &result.Tricks, &result.Declarer, &result.OpeningLead,
		); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, rows.Err()
}

func (r *PostgresBoardResultReadRepository) FindByTeams(tournamentId string, dealNo int, teamAId string, teamBId string) (*domain.BoardResultDto, error) {
	row := r.Tx.QueryRowContext(r.Ctx, `
		SELECT board_result.deal_no, board_result.ns_team_id, board_result.ew_team_id,
			ns_team.number, ew_team.number,
			board_result.contract, board_result.tricks, board_result.declarer, board_result.opening_lead
		FROM play.board_result AS board_result
		LEFT JOIN play.team AS ns_team
			ON board_result.ns_team_id = ns_team.team_id
			AND board_result.tournament_id = ns_team.tournament_id
		LEFT JOIN play.team AS ew_team
			ON board_result.ew_team_id = ew_team.team_id
			AND board_result.tournament_id = ew_team.tournament_id
		WHERE board_result.tournament_id = $1 AND board_result.deal_no = $2
			AND ((board_result.ns_team_id = $3 AND board_result.ew_team_id = $4) OR (board_result.ns_team_id = $4 AND board_result.ew_team_id = $3))`,
		tournamentId, dealNo, teamAId, teamBId)

	var result domain.BoardResultDto
	if err := row.Scan(
		&result.DealNo, &result.NsTeamId, &result.EwTeamId,
		&result.NsTeamNumber, &result.EwTeamNumber,
		&result.Contract, &result.Tricks, &result.Declarer, &result.OpeningLead,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBoardResultNotFound
		}
		return nil, err
	}
	return &result, nil
}
