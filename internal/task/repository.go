package task

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Handles Tasks

func (r Repository) CreateTask(t Task) (Task, error) {
	var dueTime any

	if t.DueAt != nil && !t.DueAt.IsZero() {
		dueTime = t.DueAt.Format("2006-01-02 15:04:05")
	} else {
		dueTime = nil
	}

	res, err := r.db.Exec(
		"INSERT INTO tasks (user_id, title, description, completed, tag, due_at, created_at) VALUES (?, ?, ?, ?, ?, ?,?)",
		t.UserID,
		t.Title,
		t.Description,
		t.Completed,
		t.Tag,
		dueTime,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return Task{}, fmt.Errorf("insert task: exec: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Task{}, fmt.Errorf("insert task: last insert id: %w", err)
	}

	t.ID = int(id)
	return t, nil
}

func (r Repository) DeleteTask(ID int) error {
	result, err := r.db.Exec(
		"DELETE FROM tasks WHERE ID = ?",
		ID,
	)
	if err != nil {
		return fmt.Errorf("delete task: exec: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("delete task: rows affected: %w", err)
	}
	return nil
}

func (r Repository) CompleteTask(t Task) error {
	result, err := r.db.Exec(
		"UPDATE tasks SET completed = true, completed_at = ? WHERE id = ?",
		time.Now().Format("2006-01-02 15:04:05"),
		t.ID,
	)
	if err != nil {
		return fmt.Errorf("complete task: exec: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("complete task: row affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("complete task: not found: %w", err)
	}
	return nil
}

func (r Repository) UpdateTask(t Task) error {
	res, err := r.db.Exec(`
		UPDATE tasks
		SET
			title = ?,
			description = ?,
			due_at = ?,
			tag = ?
		WHERE id = ? AND completed = 0
	`,
		t.Title,
		t.Description,
		t.DueAt,
		t.Tag,
		t.ID,
	)
	if err != nil {
		return fmt.Errorf("update task: exec: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update task: rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (r Repository) GetTask(t Task) (Task, error) {
	err := r.db.QueryRow(
		"SELECT id,user_id,title,description,due_at,completed,tag,created_at,completed_at FROM tasks WHERE ID = ? and user_id = ?",
		t.ID,
		t.UserID,
	).Scan(
		&t.ID,
		&t.UserID,
		&t.Title,
		&t.Description,
		&t.DueAt,
		&t.Completed,
		&t.Tag,
		&t.CreatedAt,
		&t.CompletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrTaskNotFound
		}
		return Task{}, fmt.Errorf("get task: exec: %w", err)
	}

	err = r.db.QueryRow(
		"SELECT id,task_id,pokemon_id,pokemon_name,pokemon_sprite,rarity,shiny,revealed,generated_at,revealed_at FROM task_rewards WHERE task_id = ?",
		t.ID,
	).Scan(
		&t.Reward.ID,
		&t.Reward.TaskID,
		&t.Reward.PokemonID,
		&t.Reward.PokemonName,
		&t.Reward.Sprite,
		&t.Reward.Rarity,
		&t.Reward.Shiny,
		&t.Reward.Revealed,
		&t.Reward.GeneratedAt,
		&t.Reward.RevealedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrTaskRewardNotFound
		}
		return Task{}, fmt.Errorf("get task reward: exec: %w", err)
	}

	return t, err
}

func (r Repository) ExistTask(t Task) (Task, error) {
	err := r.db.QueryRow(
		"SELECT id, user_id, completed FROM tasks WHERE ID = ? and user_id = ?",
		t.ID,
		t.UserID,
	).Scan(
		&t.ID,
		&t.UserID,
		&t.Completed,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrTaskNotFound
		}
		return Task{}, fmt.Errorf("exist task: exec: %w", err)
	}

	err = r.db.QueryRow(
		"SELECT id,task_id,revealed from task_rewards WHERE task_id = ?",
		t.ID,
	).Scan(
		&t.Reward.ID,
		&t.Reward.TaskID,
		&t.Reward.Revealed,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrTaskRewardNotFound
		}
		return Task{}, fmt.Errorf("exist task reward: exec: %w", err)
	}

	return t, nil
}

func (r Repository) ListTasksByUser(t Task) ([]Task, error) {
	rows, err := r.db.Query(`
	SELECT
		t.id,
		t.user_id,
		t.title,
		t.description,
		t.completed,
		t.due_at,
		t.tag,
		t.created_at,
		t.completed_at,

		tr.id,
		tr.task_id,
		tr.pokemon_id,
		tr.pokemon_name,
		tr.pokemon_sprite,
		tr.rarity,
		tr.shiny,
		tr.revealed,
		tr.generated_at,
		tr.revealed_at

	FROM tasks t
	LEFT JOIN task_rewards tr
		ON tr.task_id = t.id
	WHERE t.user_id = ?
	`, t.UserID)
	if err != nil {
		return nil, fmt.Errorf("list tasks: query: %w", err)
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Title,
			&t.Description,
			&t.Completed,
			&t.DueAt,
			&t.Tag,
			&t.CreatedAt,
			&t.CompletedAt,

			&t.Reward.ID,
			&t.Reward.TaskID,
			&t.Reward.PokemonID,
			&t.Reward.PokemonName,
			&t.Reward.Sprite,
			&t.Reward.Rarity,
			&t.Reward.Shiny,
			&t.Reward.Revealed,
			&t.Reward.GeneratedAt,
			&t.Reward.RevealedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("list tasks: scan: %w", err)
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list tasks: rows iteration: %w", err)
	}

	return tasks, nil
}

func (r Repository) ListCompletedTasks() ([]Task, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, title, description, due_at, completed, tag, created_at, completed_at FROM tasks",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var t Task

		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Title,
			&t.Description,
			&t.DueAt,
			&t.Completed,
			&t.Tag,
			&t.CreatedAt,
			&t.CompletedAt,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Handles Task Rewards

func (r Repository) CreateTaskReward(tr TaskReward) error {
	_, err := r.db.Exec(
		"INSERT INTO task_rewards (task_id,pokemon_id,pokemon_name,pokemon_sprite,rarity,shiny,generated_at) VALUES (?,?,?,?,?,?,?)",
		tr.TaskID,
		tr.PokemonID,
		tr.PokemonName,
		tr.Sprite,
		tr.Rarity,
		tr.Shiny,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return fmt.Errorf("insert task reward: exec: %w", err)
	}
	return nil
}

func (r Repository) RevealPokemon(tr TaskReward) error {
	result, err := r.db.Exec(
		"UPDATE task_rewards SET revealed = 1, revealed_at = ? WHERE id = ?",
		time.Now().Format("2006-01-02 15:04:05"),
		tr.ID,
	)
	if err != nil {
		return fmt.Errorf("reveal pokemon: exec: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("reveal pokemon: rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("reveal pokemon: not found: %w", err)
	}
	return nil
}

func (r Repository) DeleteTaskReward(ID int) error {
	result, err := r.db.Exec(
		"DELETE FROM task_rewards WHERE ID = ? and Revealed = 0",
		ID,
	)
	if err != nil {
		return fmt.Errorf("delete task reward: exec: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("delete task reward: rows affected: %w", err)
	}

	return nil
}

func (r Repository) GetTaskReward(t Task) (TaskReward, error) {
	var tr TaskReward
	err := r.db.QueryRow(
		"SELECT id,task_id,pokemon_id,pokemon_name,pokemon_sprite,rarity,shiny,revealed,generated_at,revealed_at FROM task_rewards WHERE task_id = ?",
		t.ID,
	).Scan(
		&tr.ID,
		&tr.TaskID,
		&tr.PokemonID,
		&tr.PokemonName,
		&tr.Sprite,
		&tr.Rarity,
		&tr.Shiny,
		&tr.Revealed,
		&tr.GeneratedAt,
		&tr.RevealedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TaskReward{}, nil
		}
		return TaskReward{}, err
	}
	return tr, nil
}

func (r Repository) ListRevealedPokemons() ([]TaskReward, error) {
	rows, err := r.db.Query(
		"SELECT id,task_id,pokemon_id,pokemon_name,pokemon_sprite,rarity,shiny,revealed,generated_at,revealed_at FROM task_rewards WHERE revealed = 1",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taskRewards []TaskReward

	for rows.Next() {
		var tr TaskReward

		err := rows.Scan(
			&tr.ID,
			&tr.TaskID,
			&tr.PokemonID,
			&tr.PokemonName,
			&tr.Sprite,
			&tr.Rarity,
			&tr.Shiny,
			&tr.Revealed,
			&tr.GeneratedAt,
			&tr.RevealedAt,
		)
		if err != nil {
			return nil, err
		}

		taskRewards = append(taskRewards, tr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return taskRewards, nil
}

// Handles Collection Entries

func (r Repository) CreateCollectionEntry(ce CollectionEntry) error {
	_, err := r.db.Exec(
		"INSERT INTO collection_entries (user_id,pokemon_id,pokemon_name,rarity,shiny,first_caught_at) VALUES (?,?,?,?,?,?)",
		ce.UserID,
		ce.PokemonID,
		ce.PokemonName,
		ce.Rarity,
		ce.Shiny,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return fmt.Errorf("create collection: exec: %w", err)
	}
	return err
}

func (r Repository) ExistCollectionEntry(userID string, pokemonID int) (int, error) {
	var ce CollectionEntry
	err := r.db.QueryRow(
		"SELECT id FROM collection_entries WHERE user_id = ? AND pokemon_id = ?",
		userID,
		pokemonID,
	).Scan(
		&ce.ID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("exist collection: exec: %w", err)
	}
	return ce.ID, nil
}

func (r Repository) UpdateCollectionEntry(pokemonID int, shiny bool, userID string) error {
	res, err := r.db.Exec(
		"UPDATE CollectionEntry SET count = count + 1, shiny = ? WHERE user_id = ? AND pokemon_id = ?",
		shiny,
		userID,
		pokemonID,
	)
	if err != nil {
		return fmt.Errorf("update collection: exec: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update collection: rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("update collection: not found: %w", err)
	}

	return nil
}

func (r Repository) ListCollection(userID string) ([]CollectionEntry, error) {
	rows, err := r.db.Query(`
		SELECT pokemon_id, pokemon_name, count, rarity, shiny,first_caught_at,last_caught_at 
		FROM collection_entries
		WHERE user_id = ?
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list collection: query: %w", err)
	}
	defer rows.Close()

	var collection []CollectionEntry

	for rows.Next() {
		var c CollectionEntry

		err := rows.Scan(
			&c.PokemonID,
			&c.PokemonName,
			&c.Count,
			&c.Rarity,
			&c.Shiny,
			&c.FirstCaughtAt,
			&c.LastCaughtAt,
		)
		if err != nil {
			return nil, fmt.Errorf("list collection: scan: %w", err)
		}

		collection = append(collection, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list collection: rows iteration: %w", err)
	}

	return collection, nil
}

// Handles Statistics

func (r Repository) CreateUserStatistic(us UserStatistic) error {
	_, err := r.db.Exec(
		"INSERT INTO user_statistics (user_id,tasks_completed,tasks_opened,tasks_deleted,pokemon_caught,shiny_caught,unique_pokemon,current_streak,longest_streak) VALUES (?,?,?,?,?,?,?,?,?)",
		us.UserID,
		us.TasksCompleted,
		us.TasksOpened,
		us.TasksDeleted,
		us.PokemonCaught,
		us.ShinyCaught,
		us.UniquePokemon,
		us.CurrentStreak,
		us.LongestStreak,
	)
	if err != nil {
		return fmt.Errorf("create statistic: exec: %w", err)
	}
	return nil
}

func (r Repository) UpdateUserStatisticOnClose(userID string, shinyCaught, uniquePokemon, currentStreak, longestStreak int) error {
	res, err := r.db.Exec(`
		UPDATE user_statistics
		SET
			tasks_completed = tasks_completed + 1,
			pokemon_caught = pokemon_caught + 1,
			shiny_caught = ?,
			unique_pokemon = ?,
			current_streak = ?,
			longest_streak = ?
		WHERE user_id = ?
	`,
		shinyCaught,
		uniquePokemon,
		currentStreak,
		longestStreak,
		userID,
	)

	if err != nil {
		return fmt.Errorf("update statistic close: exec: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update statistic close: rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("update statistic close: not found: %w", err)
	}

	return nil
}

func (r Repository) UpdateUserStatisticOnCreate(userID string) error {
	res, err := r.db.Exec(`
		UPDATE user_statistics
		SET tasks_opened = tasks_opened + 1
		WHERE user_id = ?
	`, userID)

	if err != nil {
		return fmt.Errorf("update statistic create: exec: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update statistic create: rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("update statistic create: not found: %w", err)
	}

	return nil
}

func (r Repository) UpdateUserStatisticOnDelete(userID string) error {
	res, err := r.db.Exec(`
		UPDATE user_statistics
		SET tasks_deleted = tasks_deleted + 1
		WHERE user_id = ?
	`, userID)

	if err != nil {
		return fmt.Errorf("update statistic delete: exec: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update statistic delete: rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("update statistic delete: not found: %w", err)
	}

	return nil
}

func (r Repository) GetDataForStatistic(userID string) (int, int, int, []time.Time, error) {
	//Get reduced statistics
	var (
		longestStreak      int
		shinyCaughtTotal   int
		uniquePokemonTotal int
	)

	// Get longest streak
	err := r.db.QueryRow(
		`SELECT longest_streak
		 FROM user_statistics
		 WHERE user_id = ?`,
		userID,
	).Scan(&longestStreak)
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("get statistic data: longest streak: %w", err)
	}

	// Get total shiny Pokémon caught
	err = r.db.QueryRow(`
		SELECT COUNT(*)
		FROM task_rewards tr
		JOIN tasks t ON t.id = tr.task_id
		WHERE t.user_id = ?
		  AND tr.shiny = 1
	`, userID).Scan(&shinyCaughtTotal)
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("get statistic data: shiny count: %w", err)
	}

	// Get total unique Pokémon
	err = r.db.QueryRow(`
		SELECT COUNT(*)
		FROM collection_entries
		WHERE user_id = ?
	`, userID).Scan(&uniquePokemonTotal)
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("get statistic data: unique pokemon count: %w", err)
	}

	// Get completion dates (newest first)
	rows, err := r.db.Query(`
		SELECT DISTINCT DATE(completed_at)
		FROM tasks
		WHERE user_id = ?
		  AND completed = 1
		ORDER BY DATE(completed_at) DESC
	`, userID)
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("get statistic data: date query: %w", err)
	}
	defer rows.Close()

	var dates []time.Time

	for rows.Next() {
		var dateStr string

		if err := rows.Scan(&dateStr); err != nil {
			return 0, 0, 0, nil, fmt.Errorf("get statistic data: date scan: %w", err)
		}

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return 0, 0, 0, nil, fmt.Errorf("get statistic data: parse date: %w", err)
		}

		dates = append(dates, date)
	}

	if err := rows.Err(); err != nil {
		return 0, 0, 0, nil, fmt.Errorf("get statistic data: iteration date: %w", err)
	}

	return longestStreak, shinyCaughtTotal, uniquePokemonTotal, dates, nil
}

func (r Repository) GetStatistic(userID string) (UserStatistic, error) {
	var us UserStatistic
	err := r.db.QueryRow(`
		SELECT user_id,tasks_completed,tasks_opened,tasks_deleted,pokemon_caught,shiny_caught,unique_pokemon,current_streak,longest_streak 
		FROM user_statistics
		WHERE user_id = ?`,
		userID,
	).Scan(
		&us.UserID,
		&us.TasksOpened,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserStatistic{}, fmt.Errorf("get statistic: not found: %w", err)
		}
		return UserStatistic{}, fmt.Errorf("get statistic: exec: %w", err)
	}
	return us, nil
}

func (r Repository) ExistStatistic(userID string) (UserStatistic, error) {
	var us UserStatistic
	err := r.db.QueryRow(`
		SELECT user_id,tasks_opened
		FROM user_statistics
		WHERE user_id = ?`,
		userID,
	).Scan(
		&us.UserID,
		&us.TasksOpened,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserStatistic{}, nil
		}
		return UserStatistic{}, fmt.Errorf("exist statistic: exec: %w", err)
	}
	return us, nil
}
