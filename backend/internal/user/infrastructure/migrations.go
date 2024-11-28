package user

import (
	"database/sql"
)

// constraint naming: https://stackoverflow.com/a/4108266

func Migrate(db *sql.DB) {
	m0001_initial(db)
}

func m0001_initial(db *sql.DB) {
	_, err := db.Exec(`
		CREATE SCHEMA IF NOT EXISTS user;

		CREATE TABLE IF NOT EXISTS user.users (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			name TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		);
	`)
	if err != nil {
		panic(err)
	}
}
