package tournament_management

import (
	"database/sql"
)

// constraint naming: https://stackoverflow.com/a/4108266

func Migrate(db *sql.DB) {
	m0001_initial(db)
}

func m0001_initial(db *sql.DB) {
	_, err := db.Exec(`
		CREATE SCHEMA IF NOT EXISTS tournament_management;

		CREATE TABLE IF NOT EXISTS tournament_management.tournament (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			name VARCHAR(100) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			started_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS tournament_management.team (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			tournament_id UUID NOT NULL REFERENCES tournament_management.tournament (id),
			name VARCHAR(100) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

			CONSTRAINT team_name_unique UNIQUE (tournament_id, name)
		);

		CREATE TABLE IF NOT EXISTS tournament_management.contestant (
			id UUID PRIMARY KEY NOT NULL,
			tournament_id UUID NOT NULL REFERENCES tournament_management.tournament (id),
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

			CONSTRAINT contestant_tournament_unique UNIQUE (tournament_id, id)
		);

		CREATE TABLE IF NOT EXISTS tournament_management.team_contestant (
			team_id UUID NOT NULL REFERENCES tournament_management.team (id),
			contestant_id UUID NOT NULL REFERENCES tournament_management.contestant (id),
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (team_id, contestant_id)
		);
	`)
	if err != nil {
		panic(err)
	}
}
