package database

import (
	"database/sql"
	_ "embed"

	_ "modernc.org/sqlite"
)

func Connect(path string) (*sql.DB, error) {
	return sql.Open(
		"sqlite",
		path,
	)
}

//go:embed migrations.sql
var migrationsSQL string

func Migrate(db *sql.DB) error {
	_, err := db.Exec(migrationsSQL)
	return err
}
