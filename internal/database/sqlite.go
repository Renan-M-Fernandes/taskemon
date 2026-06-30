package database

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

func Connect() (*sql.DB, error) {
	return sql.Open(
		"sqlite",
		"./database/taskemon.db",
	)
}

func Migrate(db *sql.DB) error {
	migration, err := os.ReadFile("./internal/database/migrations.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(migration))
	return err
}
