package tournament_management

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	"bridge-tab/internal/idutil"
	domain "bridge-tab/internal/tournament-management/domain"
)

type PostgresTournamentRepository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

var ErrTournamentNotFound = errors.New("tournament not found")

func (r *PostgresTournamentRepository) Load(Id *domain.TournamentId) (*domain.Tournament, error) {
	var Tournament domain.Tournament
	var StartedAt sql.NullString
	var FinishedAt sql.NullString
	row := r.Tx.QueryRowContext(r.Ctx, "SELECT id, name, started_at, finished_at FROM tournament_management.tournament WHERE id = $1", Id)
	err := row.Scan(&Tournament.State.Id, &Tournament.State.Name, &StartedAt, &FinishedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", ErrTournamentNotFound, Id)
		}
		return nil, err
	}

	contestantRows, err := r.Tx.QueryContext(r.Ctx, "SELECT id FROM tournament_management.contestant WHERE tournament_id = $1", Id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	var Contestants []*domain.Contestant
	for contestantRows.Next() {
		var contestant domain.Contestant
		err = contestantRows.Scan(&contestant.Id)
		if err != nil {
			return nil, err
		}
		Contestants = append(Contestants, &contestant)
	}

	teamRows, err := r.Tx.QueryContext(r.Ctx, "SELECT id, name, number FROM tournament_management.team WHERE Tournament_id = $1", Id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	var Teams []*domain.Team
	for teamRows.Next() {
		var team domain.Team
		err = teamRows.Scan(&team.State.Id, &team.State.Name, &team.State.Number)
		if err != nil {
			return nil, err
		}
		team.State.TournamentId = *Id
		Teams = append(Teams, &team)
	}

	tableRows, err := r.Tx.QueryContext(r.Ctx, "SELECT id, number FROM tournament_management.tournament_table WHERE tournament_id = $1", Id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	var Tables []*domain.Table
	for tableRows.Next() {
		var table domain.Table
		err = tableRows.Scan(&table.State.Id, &table.State.Number)
		if err != nil {
			return nil, err
		}
		table.State.TournamentId = *Id
		Tables = append(Tables, &table)
	}

	teamContestantRows, err := r.Tx.QueryContext(r.Ctx, `
	SELECT team_id, contestant_id FROM tournament_management.team_contestant 
	INNER JOIN tournament_management.team ON team_contestant.team_id = team.id 
	WHERE team.tournament_id = $1`, Id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	for teamContestantRows.Next() {
		var teamContestant struct {
			ContestantId string
			TeamId       string
		}
		err = teamContestantRows.Scan(&teamContestant.TeamId, &teamContestant.ContestantId)
		if err != nil {
			return nil, err
		}
		for _, team := range Teams {
			if idutil.SameId(team.State.Id, domain.TeamId(teamContestant.TeamId)) {
				contestantIndex := slices.IndexFunc(Contestants, func(c *domain.Contestant) bool {
					return idutil.SameId(c.Id, domain.ContestantId(teamContestant.ContestantId))
				})

				if contestantIndex != -1 {
					team.State.Members = append(team.State.Members, Contestants[contestantIndex])
					Contestants[contestantIndex].Team = team
				}
			}
		}
	}

	setRows, err := r.Tx.QueryContext(r.Ctx, "SELECT id, label FROM tournament_management.tournament_set WHERE tournament_id = $1", Id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	var Sets []*domain.Set
	for setRows.Next() {
		var set domain.Set
		err = setRows.Scan(&set.Id, &set.Label)
		if err != nil {
			return nil, err
		}
		set.BoardNos = make([]int, 0)
		set.TeamPairs = make([]domain.TeamPairs, 0)
		Sets = append(Sets, &set)
	}

	boardProtocolRows, err := r.Tx.QueryContext(r.Ctx, "SELECT board_no, vulnerable, set_id FROM tournament_management.board_protocol WHERE tournament_id = $1", Id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	var BoardProtocols []*domain.BoardProtocol
	for boardProtocolRows.Next() {
		var boardProtocol domain.BoardProtocol
		var setId sql.NullString
		err = boardProtocolRows.Scan(&boardProtocol.BoardNo, &boardProtocol.Vulnerable, &setId)
		if err != nil {
			return nil, err
		}
		boardProtocol.TeamPairs = make([]domain.TeamPairs, 0)
		if setId.Valid {
			id := domain.SetId(setId.String)
			boardProtocol.SetId = &id
			for _, set := range Sets {
				if idutil.SameId(set.Id, id) {
					set.BoardNos = append(set.BoardNos, boardProtocol.BoardNo)
				}
			}
		}
		BoardProtocols = append(BoardProtocols, &boardProtocol)
	}

	teamBoardProtocolRows, err := r.Tx.QueryContext(r.Ctx, `
	SELECT team_ns_id, team_ew_id, board_no, table_id FROM tournament_management.board_protocol_team_pairs 
	WHERE tournament_id = $1`, Id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	for teamBoardProtocolRows.Next() {
		var teamBoardProtocol struct {
			TeamNsId string
			TeamEwId string
			BoardNo  int
			TableId  sql.NullString
		}
		err = teamBoardProtocolRows.Scan(&teamBoardProtocol.TeamNsId, &teamBoardProtocol.TeamEwId, &teamBoardProtocol.BoardNo, &teamBoardProtocol.TableId)
		if err != nil {
			return nil, err
		}
		// verify if team ns and team ew exist
		teamNsExists := slices.ContainsFunc(Teams, func(t *domain.Team) bool {
			return idutil.SameId(t.State.Id, domain.TeamId(teamBoardProtocol.TeamNsId))
		})
		teamEwExists := slices.ContainsFunc(Teams, func(t *domain.Team) bool {
			return idutil.SameId(t.State.Id, domain.TeamId(teamBoardProtocol.TeamEwId))
		})

		if teamNsExists && teamEwExists {
			for _, boardProtocol := range BoardProtocols {
				if boardProtocol.BoardNo == teamBoardProtocol.BoardNo {
					pair := domain.TeamPairs{
						NS: domain.TeamId(teamBoardProtocol.TeamNsId),
						EW: domain.TeamId(teamBoardProtocol.TeamEwId),
					}
					if teamBoardProtocol.TableId.Valid {
						tableId := domain.TableId(teamBoardProtocol.TableId.String)
						pair.Table = &tableId
					}
					boardProtocol.TeamPairs = append(boardProtocol.TeamPairs, pair)
				}
			}
		}
	}

	for _, set := range Sets {
		for _, boardProtocol := range BoardProtocols {
			if boardProtocol.SetId != nil && idutil.SameId(*boardProtocol.SetId, set.Id) {
				if len(set.TeamPairs) == 0 {
					set.TeamPairs = append([]domain.TeamPairs(nil), boardProtocol.TeamPairs...)
				}
				break
			}
		}
	}

	Tournament.State.Contestants = Contestants
	Tournament.State.Teams = Teams
	Tournament.State.Tables = Tables
	Tournament.State.BoardProtocols = BoardProtocols
	Tournament.State.Sets = Sets
	if StartedAt.Valid {
		startedAtTime, err := time.Parse(time.RFC3339Nano, StartedAt.String)

		if err == nil {
			Tournament.State.StartedAt = &startedAtTime
		}
	}
	if FinishedAt.Valid {
		finishedAtTime, err := time.Parse(time.RFC3339Nano, FinishedAt.String)

		if err == nil {
			Tournament.State.FinishedAt = &finishedAtTime
		}
	}

	return &Tournament, nil
}

func (r *PostgresTournamentRepository) Save(t *domain.Tournament) error {
	for _, event := range t.GetEvents() {
		var err error
		switch event := event.(type) {
		case domain.TournamentCreated:
			err = r.TournamentCreated(event)
		case domain.TournamentRemoved:
			err = r.TournamentRemoved(event)
		case domain.TournamentStarted:
			err = r.TournamentStarted(event)
		case domain.TournamentFinished:
			err = r.TournamentFinished(event)
		case domain.SetCreated:
			err = r.setCreated(event)
		case domain.SetRemoved:
			err = r.setRemoved(event)
		case domain.ContestantJoinedTournament:
			err = r.contestantJoinedTournament(event)
		case domain.ContestantLeftTournament:
			err = r.contestantLeftTournament(event)
		case domain.TeamCreated:
			err = r.teamCreated(event)
		case domain.TeamRenamed:
			err = r.teamRenamed(event)
		case domain.TeamRemoved:
			err = r.teamRemoved(event)
		case domain.TableCreated:
			err = r.tableCreated(event)
		case domain.TableRemoved:
			err = r.tableRemoved(event)
		case domain.ContestantJoinedTeam:
			err = r.contestantJoinedTeam(event)
		case domain.ContestantLeftTeam:
			err = r.contestantLeftTeam(event)
		case domain.BoardProtocolCreated:
			err = r.boardProtocolCreated(event)
		case domain.BoardProtocolRemoved:
			err = r.boardProtocolRemoved(event)

		default:
			return errors.New("unknown event")
		}
		if err != nil {
			return err
		}
	}
	t.Commit()
	return nil
}

func (r *PostgresTournamentRepository) TournamentCreated(event domain.TournamentCreated) error {
	_, err := r.Tx.ExecContext(r.Ctx, "INSERT INTO tournament_management.tournament (id, name) VALUES ($1, $2)", event.TournamentId, event.Name)

	return err
}

func (r *PostgresTournamentRepository) TournamentRemoved(event domain.TournamentRemoved) error {
	_, err := r.Tx.ExecContext(r.Ctx, "DELETE FROM tournament_management.tournament WHERE id = $1", event.TournamentId)

	return err
}

func (r *PostgresTournamentRepository) TournamentStarted(event domain.TournamentStarted) error {
	_, err := r.Tx.ExecContext(r.Ctx, "UPDATE tournament_management.tournament SET started_at = $1 WHERE id = $2", event.StartedAt, event.TournamentId)

	return err
}

func (r *PostgresTournamentRepository) TournamentFinished(event domain.TournamentFinished) error {
	_, err := r.Tx.ExecContext(r.Ctx, "UPDATE tournament_management.tournament SET finished_at = $1 WHERE id = $2", event.FinishedAt, event.TournamentId)

	return err
}

func (r *PostgresTournamentRepository) setCreated(event domain.SetCreated) error {
	_, err := r.Tx.ExecContext(r.Ctx, "INSERT INTO tournament_management.tournament_set (id, tournament_id, label) VALUES ($1, $2, $3)", event.SetId, event.TournamentId, event.Label)
	if err != nil {
		return err
	}

	for _, boardNo := range event.BoardNos {
		_, err = r.Tx.ExecContext(r.Ctx, "UPDATE tournament_management.board_protocol SET set_id = $1 WHERE tournament_id = $2 AND board_no = $3", event.SetId, event.TournamentId, boardNo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PostgresTournamentRepository) setRemoved(event domain.SetRemoved) error {
	_, err := r.Tx.ExecContext(r.Ctx, "UPDATE tournament_management.board_protocol SET set_id = NULL WHERE tournament_id = $1 AND set_id = $2", event.TournamentId, event.SetId)
	if err != nil {
		return err
	}

	_, err = r.Tx.ExecContext(r.Ctx, "DELETE FROM tournament_management.tournament_set WHERE id = $1 AND tournament_id = $2", event.SetId, event.TournamentId)
	return err
}

func (r *PostgresTournamentRepository) contestantJoinedTournament(event domain.ContestantJoinedTournament) error {
	_, err := r.Tx.ExecContext(r.Ctx, "INSERT INTO tournament_management.contestant (id, tournament_id) VALUES ($1, $2)", event.ContestantId, event.TournamentId)

	return err
}

func (r *PostgresTournamentRepository) contestantLeftTournament(event domain.ContestantLeftTournament) error {
	_, err := r.Tx.ExecContext(r.Ctx, "DELETE FROM tournament_management.contestant WHERE id = $1 AND tournament_id = $2", event.ContestantId, event.TournamentId)

	return err
}

func (r *PostgresTournamentRepository) teamCreated(event domain.TeamCreated) error {
	_, err := r.Tx.ExecContext(r.Ctx, "INSERT INTO tournament_management.team (id, Tournament_id, name, number) VALUES ($1, $2, $3, $4)", event.TeamId, event.TournamentId, event.Name, event.Number)

	return err
}

func (r *PostgresTournamentRepository) teamRenamed(event domain.TeamRenamed) error {
	_, err := r.Tx.ExecContext(r.Ctx, "UPDATE tournament_management.team SET name = $1 WHERE id = $2", event.Name, event.TeamId)
	return err
}

func (r *PostgresTournamentRepository) teamRemoved(event domain.TeamRemoved) error {
	_, err := r.Tx.ExecContext(r.Ctx, "DELETE FROM tournament_management.team WHERE id = $1 AND tournament_id = $2", event.TeamId, event.TournamentId)

	return err
}

func (r *PostgresTournamentRepository) tableCreated(event domain.TableCreated) error {
	_, err := r.Tx.ExecContext(r.Ctx, "INSERT INTO tournament_management.tournament_table (id, tournament_id, number) VALUES ($1, $2, $3)", event.TableId, event.TournamentId, event.Number)

	return err
}

func (r *PostgresTournamentRepository) tableRemoved(event domain.TableRemoved) error {
	_, err := r.Tx.ExecContext(r.Ctx, "DELETE FROM tournament_management.tournament_table WHERE id = $1 AND tournament_id = $2", event.TableId, event.TournamentId)

	return err
}

func (r *PostgresTournamentRepository) contestantJoinedTeam(event domain.ContestantJoinedTeam) error {
	_, err := r.Tx.ExecContext(r.Ctx, "INSERT INTO tournament_management.team_contestant (team_id, contestant_id) VALUES ($1, $2)", event.TeamId, event.ContestantId)

	return err
}

func (r *PostgresTournamentRepository) contestantLeftTeam(event domain.ContestantLeftTeam) error {
	_, err := r.Tx.ExecContext(r.Ctx, "DELETE FROM tournament_management.team_contestant WHERE team_id = $1 AND contestant_id = $2", event.TeamId, event.ContestantId)

	return err
}

func (r *PostgresTournamentRepository) boardProtocolCreated(event domain.BoardProtocolCreated) error {
	_, err := r.Tx.ExecContext(r.Ctx, "INSERT INTO tournament_management.board_protocol (tournament_id, board_no, vulnerable) VALUES ($1, $2, $3)", event.TournamentId, event.BoardNo, event.Vulnerable)

	if err != nil {
		return err
	}

	if len(event.TeamPairs) == 0 {
		return nil
	}

	var valuesQuery string
	var values []interface{}
	for i, teamPair := range event.TeamPairs {
		var tableId interface{}
		if teamPair.Table != nil {
			tableId = *teamPair.Table
		}
		values = append(values, event.TournamentId, event.BoardNo, teamPair.NS, teamPair.EW, tableId)
		valuesQuery += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d), ", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
	}
	valuesQuery = valuesQuery[:len(valuesQuery)-2]

	_, err = r.Tx.ExecContext(r.Ctx, fmt.Sprintf("INSERT INTO tournament_management.board_protocol_team_pairs (tournament_id, board_no, team_ns_id, team_ew_id, table_id) VALUES %s", valuesQuery), values...)

	return err
}

func (r *PostgresTournamentRepository) boardProtocolRemoved(event domain.BoardProtocolRemoved) error {
	_, err := r.Tx.ExecContext(r.Ctx, "DELETE FROM tournament_management.board_protocol_team_pairs WHERE tournament_id = $1 AND board_no = $2", event.TournamentId, event.BoardNo)

	if err != nil {
		return err
	}

	_, err = r.Tx.ExecContext(r.Ctx, "DELETE FROM tournament_management.board_protocol WHERE tournament_id = $1 AND board_no = $2", event.TournamentId, event.BoardNo)

	return err
}
