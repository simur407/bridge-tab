package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
	"context"
	"database/sql"
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

	var boardProtocols []domain.BoardProtocolDto
	for boardProtocolRows.Next() {
		var boardProtocol domain.BoardProtocolDto
		var vulnerable domain.Vulnerable
		err := boardProtocolRows.Scan(&boardProtocol.BoardNo, &vulnerable)
		if err != nil {
			return nil, err
		}

		var readableVulnerable string
		switch vulnerable {
		case domain.None:
			readableVulnerable = "None"
		case domain.NS:
			readableVulnerable = "NS"
		case domain.EW:
			readableVulnerable = "EW"
		case domain.Both:
			readableVulnerable = "Both"
		default:
			readableVulnerable = "Unknown"
		}

		boardProtocol.Vulnerable = readableVulnerable
		boardProtocols = append(boardProtocols, boardProtocol)
	}

	teamPairRows, err := r.Tx.QueryContext(r.Ctx, "SELECT board_no, team_ns_id, team_ew_id FROM tournament_management.board_protocol_team_pairs WHERE tournament_id = $1", tournamentId)
	if err != nil {
		return nil, err
	}

	for teamPairRows.Next() {
		var boardNo int
		var teamPair domain.TeamPairsDto
		err := teamPairRows.Scan(&boardNo, &teamPair.NS, &teamPair.EW)
		if err != nil {
			return nil, err
		}

		for i := range boardProtocols {
			if boardProtocols[i].BoardNo == boardNo {
				boardProtocols[i].TeamPairs = append(boardProtocols[i].TeamPairs, teamPair)
			}
		}
	}

	return boardProtocols, nil
}
