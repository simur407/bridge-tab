package tournament_management

import (
	"database/sql"
)

// constraint naming: https://stackoverflow.com/a/4108266

func Migrate(db *sql.DB) {
	m0001_initial(db)
	m0002_remove_primary_key_constraint_contestant(db)
	m0003_add_team_number(db)
	m0004_optional_team_name(db)
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

func m0003_add_team_number(db *sql.DB) {
	_, err := db.Exec(`
		ALTER TABLE tournament_management.team
		ADD COLUMN IF NOT EXISTS number INTEGER;

		UPDATE tournament_management.team
		SET number = name::integer
		WHERE number IS NULL AND name ~ '^[0-9]+$';

		WITH missing AS (
			SELECT id,
				tournament_id,
				ROW_NUMBER() OVER (PARTITION BY tournament_id ORDER BY created_at, id) AS rn
			FROM tournament_management.team
			WHERE number IS NULL
		),
		bases AS (
			SELECT tournament_id, COALESCE(MAX(number), 0) AS base
			FROM tournament_management.team
			GROUP BY tournament_id
		)
		UPDATE tournament_management.team AS team
		SET number = bases.base + missing.rn
		FROM missing
		JOIN bases ON bases.tournament_id = missing.tournament_id
		WHERE team.id = missing.id;

		ALTER TABLE tournament_management.team
		ALTER COLUMN number SET NOT NULL;

		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'team_number_unique'
			) THEN
				ALTER TABLE tournament_management.team
				ADD CONSTRAINT team_number_unique UNIQUE (tournament_id, number);
			END IF;
		END $$;
	`)
	if err != nil {
		panic(err)
	}
}

func m0004_optional_team_name(db *sql.DB) {
	_, err := db.Exec(`
		ALTER TABLE tournament_management.team
		DROP CONSTRAINT IF EXISTS team_name_unique;

		CREATE UNIQUE INDEX IF NOT EXISTS team_name_unique
		ON tournament_management.team (tournament_id, name)
		WHERE name <> '';
	`)
	if err != nil {
		panic(err)
	}
}
