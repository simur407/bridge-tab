package tournament_management

import (
	"database/sql"
)

// constraint naming: https://stackoverflow.com/a/4108266

func Migrate(db *sql.DB) {
	m0001_initial(db)
	m0002_remove_primary_key_constraint_contestant(db)
}

func m0001_initial(db *sql.DB) {
	_, err := db.Exec(`
		CREATE SCHEMA IF NOT EXISTS tournament_management;

		CREATE TABLE IF NOT EXISTS tournament_management.tournament (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			name TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			started_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS tournament_management.team (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			tournament_id UUID NOT NULL REFERENCES tournament_management.tournament (id),
			name TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

			CONSTRAINT team_name_unique UNIQUE (tournament_id, name)
		);

		CREATE TABLE IF NOT EXISTS tournament_management.contestant (
			id UUID PRIMARY KEY NOT NULL,
			tournament_id UUID NOT NULL REFERENCES tournament_management.tournament (id),
			created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

			CONSTRAINT contestant_tournament_unique UNIQUE (tournament_id, id)
		);

		CREATE TABLE IF NOT EXISTS tournament_management.team_contestant (
			team_id UUID NOT NULL REFERENCES tournament_management.team (id),
			contestant_id UUID NOT NULL REFERENCES tournament_management.contestant (id),
			created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (team_id, contestant_id)
		);

		CREATE TABLE IF NOT EXISTS tournament_management.board_protocol (
			board_no INT NOT NULL,
			tournament_id UUID NOT NULL REFERENCES tournament_management.tournament (id),
			vulnerable TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (board_no, tournament_id)
		);

		CREATE TABLE IF NOT EXISTS tournament_management.board_protocol_team_pairs (
			board_no INT NOT NULL,
			tournament_id UUID NOT NULL REFERENCES tournament_management.tournament (id),
			team_ns_id UUID NOT NULL REFERENCES tournament_management.team (id),
			team_ew_id UUID NOT NULL REFERENCES tournament_management.team (id),
			created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (board_no, tournament_id, team_ns_id, team_ew_id)
		)
	`)
	if err != nil {
		panic(err)
	}
}

func m0002_remove_primary_key_constraint_contestant(db *sql.DB) {
	_, err := db.Exec(`
		ALTER TABLE tournament_management.team_contestant
		DROP CONSTRAINT IF EXISTS team_contestant_contestant_id_fkey;
		ALTER TABLE tournament_management.contestant
		DROP CONSTRAINT IF EXISTS contestant_pkey;
	`)
	if err != nil {
		panic(err)
	}
}
