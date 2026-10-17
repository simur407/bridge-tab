package tournament_management

import (
	"context"
	"database/sql"

	domain "bridge-tab/internal/tournament-management/domain"
)

type PostgresBoardProtocolReadRepository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

func (r *PostgresBoardProtocolReadRepository) FindAll(tournamentId *string) ([]domain.BoardProtocolDto, error) {
	boardProtocolRows, err := r.Tx.QueryContext(r.Ctx, "SELECT board_no, vulnerable FROM tournament_management.board_protocol WHERE tournament_id = $1", tournamentId)
	if err != nil {
		return nil, err
	}
	defer boardProtocolRows.Close()

	var boardProtocols []domain.BoardProtocolDto
	for boardProtocolRows.Next() {
		var boardProtocol domain.BoardProtocolDto
		var vulnerable domain.Vulnerable
		err := boardProtocolRows.Scan(&boardProtocol.BoardNo, &vulnerable)
		if err != nil {
			return nil, err
		}
		boardProtocol.Vulnerable = vulnerable.String()
		boardProtocols = append(boardProtocols, boardProtocol)
	}
	if err := boardProtocolRows.Err(); err != nil {
		return nil, err
	}

	return boardProtocols, nil
}
