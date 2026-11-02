package tournament_management

import (
	"context"
	"database/sql"
	"errors"

	domain "bridge-tab/internal/tournament-management/domain"
)

type PostgresBoardProtocolReadRepository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

func (r *PostgresBoardProtocolReadRepository) FindAll(tournamentId *string) ([]domain.BoardProtocolDto, error) {
	boardProtocolRows, err := r.Tx.QueryContext(r.Ctx, `
		SELECT board_no, vulnerable, north, east, south, west
		FROM tournament_management.board_protocol
		WHERE tournament_id = $1
		ORDER BY board_no`, tournamentId)
	if err != nil {
		return nil, err
	}
	defer boardProtocolRows.Close()

	var boardProtocols []domain.BoardProtocolDto
	for boardProtocolRows.Next() {
		dto, err := scanBoardProtocolDto(boardProtocolRows)
		if err != nil {
			return nil, err
		}
		boardProtocols = append(boardProtocols, dto)
	}
	if err := boardProtocolRows.Err(); err != nil {
		return nil, err
	}

	return boardProtocols, nil
}

func (r *PostgresBoardProtocolReadRepository) FindByBoardNo(tournamentId *string, boardNo int) (*domain.BoardProtocolDto, error) {
	row := r.Tx.QueryRowContext(r.Ctx, `
		SELECT board_no, vulnerable, north, east, south, west
		FROM tournament_management.board_protocol
		WHERE tournament_id = $1 AND board_no = $2`, tournamentId, boardNo)
	dto, err := scanBoardProtocolDto(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &dto, nil
}

type boardProtocolScanner interface {
	Scan(dest ...any) error
}

func scanBoardProtocolDto(row boardProtocolScanner) (domain.BoardProtocolDto, error) {
	var boardProtocol domain.BoardProtocolDto
	var vulnerable domain.Vulnerable
	var north, east, south, west sql.NullString
	err := row.Scan(&boardProtocol.BoardNo, &vulnerable, &north, &east, &south, &west)
	if err != nil {
		return domain.BoardProtocolDto{}, err
	}
	boardProtocol.Vulnerable = vulnerable.String()
	boardProtocol.North = north.String
	boardProtocol.East = east.String
	boardProtocol.South = south.String
	boardProtocol.West = west.String
	return boardProtocol, nil
}
