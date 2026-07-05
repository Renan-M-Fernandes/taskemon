package database

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func setupDatabaseTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", "file:"+t.Name()+"?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("Open test database: %v", err)
	}
	db.SetMaxOpenConns(1)

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestMigrateCreatesTables(t *testing.T) {
	db := setupDatabaseTestDB(t)

	if err := Migrate(db); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	tables := []string{
		"tasks",
		"task_rewards",
		"collection_entries",
		"user_statistics",
	}

	for _, table := range tables {
		t.Run(table, func(t *testing.T) {
			var name string
			err := db.QueryRow(`
				SELECT name
				FROM sqlite_master
				WHERE type = 'table'
				  AND name = ?
			`, table).Scan(&name)
			if err != nil {
				t.Fatalf("table %s was not created: %v", table, err)
			}
		})
	}
}

func TestTaskRewardCascadeDelete(t *testing.T) {
	db := setupDatabaseTestDB(t)

	if err := Migrate(db); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	_, err := db.Exec(`
		INSERT INTO tasks (id, user_id, title, description, tag)
		VALUES (1, 'ash', 'Catch Pikachu', '', 'pokemon')
	`)
	if err != nil {
		t.Fatalf("insert task: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO task_rewards (task_id, pokemon_id, pokemon_name, rarity)
		VALUES (1, 25, 'pikachu', 1)
	`)
	if err != nil {
		t.Fatalf("insert reward: %v", err)
	}

	_, err = db.Exec(`DELETE FROM tasks WHERE id = 1`)
	if err != nil {
		t.Fatalf("delete task: %v", err)
	}

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM task_rewards WHERE task_id = 1`).Scan(&count)
	if err != nil {
		t.Fatalf("count rewards: %v", err)
	}

	if count != 0 {
		t.Fatalf("reward count mismatch: got %d, expect 0", count)
	}
}

func TestCollectionAllowsNormalAndShinySamePokemon(t *testing.T) {
	db := setupDatabaseTestDB(t)

	if err := Migrate(db); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	_, err := db.Exec(`
		INSERT INTO collection_entries (user_id, pokemon_id, pokemon_name, shiny)
		VALUES ('ash', 25, 'pikachu', 0)
	`)
	if err != nil {
		t.Fatalf("insert normal pokemon: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO collection_entries (user_id, pokemon_id, pokemon_name, shiny)
		VALUES ('ash', 25, 'pikachu', 1)
	`)
	if err != nil {
		t.Fatalf("insert shiny pokemon: %v", err)
	}

	var count int
	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM collection_entries
		WHERE user_id = 'ash'
		  AND pokemon_id = 25
	`).Scan(&count)
	if err != nil {
		t.Fatalf("count collection entries: %v", err)
	}

	if count != 2 {
		t.Fatalf("collection count mismatch: got %d, expect 2", count)
	}
}
