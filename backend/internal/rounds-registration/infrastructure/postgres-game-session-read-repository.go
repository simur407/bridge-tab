package rounds_registration

import (
	domain "bridge-tab/internal/rounds-registration/domain"
	"context"
	"database/sql"
)

type PostgresGameSessionReadRepository struct {
	Tx  *sql.Tx
	Ctx context.Context
}

func (p *PostgresGameSessionReadRepository) FindRound(gameSessionId *string, dealNo int, playerTeamId string, versusTeamId string) (*domain.RoundDto, error) {
	row := p.Tx.QueryRowContext(p.Ctx, `
	SELECT deal_no, ns_team.name AS ns_team_name, ew_team.name AS ew_team_name 
	FROM rounds_registration.round 
	LEFT JOIN tournament_management.team AS ns_team 
		ON ns_team_id = ns_team.id
	LEFT JOIN tournament_management.team AS ew_team
		ON ew_team_id = ew_team.id
	WHERE game_session_id = $1 AND deal_no = $2 
	AND ((ns_team_id = $3 AND ew_team_id = $4) OR (ns_team_id = $4 AND ew_team_id = $3))`, gameSessionId, dealNo, playerTeamId, versusTeamId)

	var round domain.RoundDto
	if err := row.Scan(&round.DealNo, &round.NsTeamName, &round.EwTeamName); err != nil {
		return nil, err
	}

	if round.DealNo == 0 {
		return nil, nil
	}

	return &round, nil
}

func (p *PostgresGameSessionReadRepository) FindAllRounds(gameSessionId *string) ([]domain.PlayedRoundDto, error) {
	rows, err := p.Tx.QueryContext(p.Ctx, `
		SELECT deal_no, ns_team.name AS ns_team_name, ew_team.name AS ew_team_name, contract, tricks, declarer, opening_lead
		FROM rounds_registration.round 
		LEFT JOIN tournament_management.team AS ns_team 
			ON ns_team_id = ns_team.id
		LEFT JOIN tournament_management.team AS ew_team
			ON ew_team_id = ew_team.id
		WHERE game_session_id = $1
		ORDER BY deal_no, round.updated_at`, gameSessionId)

	if err != nil {
		return nil, err
	}

	var rounds []domain.PlayedRoundDto
	for rows.Next() {
		var round domain.PlayedRoundDto
		var contract sql.NullString
		var openingLead sql.NullString
		var declarer sql.NullString
		var tricks sql.NullInt64
		if err := rows.Scan(&round.DealNo, &round.NsTeamName, &round.EwTeamName, &contract, &tricks, &declarer, &openingLead); err != nil {
			return nil, err
		}

		if contract.Valid {
			round.Contract = contract.String
		}

		if openingLead.Valid {
			round.OpeningLead = openingLead.String
		}

		if declarer.Valid {
			round.Declarer = declarer.String
		}

		if tricks.Valid {
			round.Tricks = int(tricks.Int64)
		}

		rounds = append(rounds, round)
	}

	return rounds, nil
}
