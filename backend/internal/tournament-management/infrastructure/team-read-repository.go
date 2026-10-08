package tournament_management

import (
	"bridge-tab/internal/idutil"
	domain "bridge-tab/internal/tournament-management/domain"
	"context"
	"database/sql"
	"slices"
)

type PostgresTeamReadRepository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

type TeamRecord struct {
	id           string
	name         string
	number       int
	contestantId sql.NullString
}

func (r *PostgresTeamReadRepository) FindAll(tournamentId *string) ([]domain.TeamDto, error) {
	rows, err := r.Tx.QueryContext(r.Ctx, `
	SELECT team.id, team.name, team.number, team_contestant.contestant_id FROM tournament_management.team 
	LEFT JOIN tournament_management.team_contestant 
		ON team.id = team_contestant.team_id
	WHERE team.tournament_id = $1
	ORDER BY team.number ASC`, tournamentId)
	if err != nil {
		return nil, err
	}

	var teams []TeamRecord
	for rows.Next() {
		var team TeamRecord
		err := rows.Scan(&team.id, &team.name, &team.number, &team.contestantId)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	var teamDtos []domain.TeamDto
	for _, team := range teams {
		id := team.id
		teamIdx := slices.IndexFunc(teamDtos, func(t domain.TeamDto) bool {
			return idutil.SameId(t.Id, id)
		})

		if teamIdx != -1 {
			if team.contestantId.Valid {
				teamDtos[teamIdx].Members = append(teamDtos[teamIdx].Members, domain.ContestantDto{Id: domain.ContestantId(team.contestantId.String)})
			}
		} else {
			members := []domain.ContestantDto{}
			if team.contestantId.Valid {
				members = append(members, domain.ContestantDto{Id: domain.ContestantId(team.contestantId.String)})
			}
			teamDtos = append(teamDtos, domain.TeamDto{Id: id, Name: team.name, Number: team.number, Members: members})
		}
	}

	return teamDtos, nil
}

func (r *PostgresTeamReadRepository) FindById(tournamentId *string, teamId *string) (*domain.TeamDto, error) {
	row := r.Tx.QueryRowContext(r.Ctx, `
	SELECT team.id, team.name, team.number FROM tournament_management.team
	WHERE team.tournament_id = $1 AND team.id = $2`, tournamentId, teamId)
	var team domain.TeamDto
	err := row.Scan(&team.Id, &team.Name, &team.Number)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

func (r *PostgresTeamReadRepository) FindByName(tournamentId *string, name *string) (*domain.TeamDto, error) {
	row := r.Tx.QueryRowContext(r.Ctx, `
	SELECT team.id, team.name, team.number FROM tournament_management.team 
	WHERE team.tournament_id = $1 AND team.name = $2
	ORDER BY team.number ASC`, tournamentId, name)
	var team domain.TeamDto
	err := row.Scan(&team.Id, &team.Name, &team.Number)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

func (r *PostgresTeamReadRepository) FindByNumber(tournamentId *string, number int) (*domain.TeamDto, error) {
	row := r.Tx.QueryRowContext(r.Ctx, `
	SELECT team.id, team.name, team.number FROM tournament_management.team
	WHERE team.tournament_id = $1 AND team.number = $2`, tournamentId, number)
	var team domain.TeamDto
	err := row.Scan(&team.Id, &team.Name, &team.Number)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

func (r *PostgresTeamReadRepository) FindByMemberId(tournamentId *string, memberId *string) (*domain.TeamDto, error) {
	row := r.Tx.QueryRowContext(r.Ctx, `
	SELECT team.id, team.name, team.number FROM tournament_management.team AS team
	LEFT JOIN tournament_management.team_contestant AS team_contestant
		ON team.id = team_contestant.team_id
	WHERE team.tournament_id = $1
	AND team_contestant.contestant_id = $2
	ORDER BY team.number ASC`, tournamentId, memberId)
	var team domain.TeamDto
	err := row.Scan(&team.Id, &team.Name, &team.Number)
	if err != nil {
		return nil, err
	}

	return &team, nil
}
