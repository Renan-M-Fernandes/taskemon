package task

import (
	"database/sql"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Renan-M-Fernandes/taskemon/internal/database"
	_ "modernc.org/sqlite"
)

type UserTotal struct {
	UserID string
	Total  int
}

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	maxPokemonSpecies = 1025

	dsn := "file:" + strings.ReplaceAll(t.Name(), "/", "_") + "?mode=memory&cache=shared"

	db, err := sql.Open(
		"sqlite",
		dsn,
	)
	if err != nil {
		t.Fatal(err)
	}

	db.SetMaxOpenConns(1)

	if err := database.Migrate(db); err != nil {
		t.Fatalf("Migrate test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func timePtrString(t *time.Time) string {
	if t == nil {
		return "<nil>"
	}

	return t.Format("2006-01-02 15:04:05")
}

func setupService(db *sql.DB) *Service {
	repo := NewRepository(db)
	service := NewService(repo)
	return service
}

func loadTasks(t *testing.T) []Task {
	t.Helper()

	data, err := os.ReadFile("testdata/tasks.json")
	if err != nil {
		t.Fatal(err)
	}

	var tasks []Task

	if err := json.Unmarshal(data, &tasks); err != nil {
		t.Fatal(err)
	}

	return tasks
}

func createBulkTasks(t *testing.T, service *Service, count int) []Task {
	t.Helper()
	tasks := loadTasks(t)

	if count == 0 || count > len(tasks) {
		count = len(tasks)
	}

	created := make([]Task, 0, count)

	for i := 0; i < count; i++ {
		task, err := service.CreateTask(tasks[i])
		if err != nil {
			t.Fatalf("Create Task (%d): %v", i, err)
		}

		created = append(created, task)
	}

	return created
}

func completeBulkTasks(t *testing.T, service *Service, tasks []Task, limit int) []int {
	t.Helper()

	if limit <= 0 || limit > len(tasks) {
		limit = len(tasks)
	}

	completed := make([]int, 0, limit)

	for i := 0; i < limit; i++ {
		task := tasks[i]

		_, err := service.CompleteTask(task.ID, task.UserID)
		if err != nil {
			t.Fatalf("CompleteTask(%d): %v", task.ID, err)
		}

		completed = append(completed, task.ID)
	}

	return completed
}

func getUserCount(t *testing.T, tasks []Task) []UserTotal {
	t.Helper()

	counts := make(map[string]int)

	for _, tt := range tasks {
		counts[tt.UserID]++
	}

	var userTotals []UserTotal

	for userID, total := range counts {
		userTotals = append(userTotals, UserTotal{
			UserID: userID,
			Total:  total,
		})
	}

	return userTotals
}

func getUserCountExcludesCompleted(t *testing.T, service *Service, usertotals []UserTotal) []UserTotal {
	t.Helper()

	var totals []UserTotal
	for _, user := range usertotals {
		task, err := service.ListTasksByUser(user.UserID)
		if err != nil {
			t.Fatalf("List tasks by user")
		}
		for _, tt := range task {
			if tt.Completed {
				user.Total--
			}
		}
		totals = append(totals, user)
	}

	return totals
}

func getUserCountCompleted(t *testing.T, service *Service, usertotals []UserTotal) []UserTotal {
	t.Helper()

	var totals []UserTotal
	for _, user := range usertotals {
		task, err := service.ListTasksByUser(user.UserID)
		if err != nil {
			t.Fatalf("List tasks by user")
		}
		for _, tt := range task {
			if !tt.Completed {
				user.Total--
			}
		}
		totals = append(totals, user)
	}

	return totals
}

func assertTaskMatchesInput(t *testing.T, got Task, expect Task, completed bool) {
	t.Helper()

	if got.ID == 0 {
		t.Fatal("ID should't be 0")
	}

	if got.UserID != expect.UserID {
		t.Fatalf("UserID mismatch: got %q, expect %q", got.UserID, expect.UserID)
	}

	if got.Title != expect.Title {
		t.Fatalf("Title mismatch: got %q, expect %q", got.Title, expect.Title)
	}

	if got.Description != expect.Description {
		t.Fatalf("Description mismatch: got %q, expect %q", got.Description, expect.Description)
	}

	expectedTag := expect.Tag
	if expectedTag == "" {
		expectedTag = "misc"
	}

	if got.Tag != expectedTag {
		t.Fatalf("Tag mismatch: got %q, expect %q", got.Tag, expectedTag)
	}

	if timePtrString(got.DueAt) != timePtrString(expect.DueAt) {
		t.Fatalf("DueAt mismatch: got %s, expect %s", timePtrString(got.DueAt), timePtrString(expect.DueAt))
	}

	if got.CreatedAt.IsZero() {
		t.Fatal("Expected createdAt to be set")
	}

	if completed {
		if !got.Completed {
			t.Fatal("Expected completed to be true")
		}

		if got.CompletedAt == nil {
			t.Fatal("Expected completedAt to be set")
		}
	} else {
		if got.Completed {
			t.Fatal("Expected completed to be false")
		}
		if timePtrString(got.CompletedAt) != "<nil>" {
			t.Fatalf("Expected createdAt to be empty")
		}
	}

}

func assertRewardMatchesInput(t *testing.T, got TaskReward, expect TaskReward, completed bool) {
	t.Helper()

	if got.ID == 0 {
		t.Fatal("ID should't be 0")
	}

	if got.TaskID != expect.TaskID {
		t.Fatalf("TaskID mismatch: got %q, expect %q", got.TaskID, expect.TaskID)
	}

	if got.PokemonID != expect.PokemonID {
		t.Fatalf("PokemonID mismatch: got %q, expect %q", got.PokemonID, expect.PokemonID)
	}

	if got.PokemonName != expect.PokemonName {
		t.Fatalf("PokemonName mismatch: got %q, expect %q", got.PokemonName, expect.PokemonName)
	}

	if got.Rarity != expect.Rarity {
		t.Fatalf("Rarity mismatch: got %q, expect %q", got.Rarity, expect.Rarity)
	}

	if got.Shiny != expect.Shiny {
		t.Fatalf("Shiny mismatch: got %v, expect %v", got.Shiny, expect.Shiny)
	}

	if got.GeneratedAt != expect.GeneratedAt {
		t.Fatalf("GeneratedAt mismatch: got %q, expect %q", got.GeneratedAt, expect.GeneratedAt)
	}

	if completed {
		if !got.Revealed {
			t.Fatal("Expected revealed to be true")
		}

		if got.RevealedAt == nil {
			t.Fatal("Expected revealedAt to be set")
		}
	} else {
		if got.Revealed {
			t.Fatal("Expected revealed to be false")
		}
		if timePtrString(got.RevealedAt) != "<nil>" {
			t.Fatalf("Expected revealedAt to be empty")
		}
	}
}

func assertNewReward(t *testing.T, r TaskReward, ID int) {
	t.Helper()

	if r.TaskID != ID {
		t.Fatalf("Id mismatch: task %q, reward %q", ID, r.TaskID)
	}
	if r.PokemonID == 0 {
		t.Fatalf("Pokemon id shouldn't be zero")
	}
	if r.PokemonName == "" {
		t.Fatalf("Pokemon name shouldn't be spaces")
	}
	if r.Sprite == "" {
		t.Fatalf("Sprite shouldn't be spaces")
	}
	if r.Revealed {
		t.Fatalf("Revealed shouldn't be true")
	}
	if r.GeneratedAt.IsZero() {
		t.Fatalf("Generated at shouldn't be zero")
	}
	if r.RevealedAt != nil {
		t.Fatalf("Revealed at at should be nil")
	}
}

func assertNewStatistics(t *testing.T, us UserStatistic) {
	t.Helper()

	if us.TasksCompleted != 0 {
		t.Fatalf("Tasks completed should be 0 %d", us.TasksCompleted)
	}
	if us.TasksDeleted != 0 {
		t.Fatalf("Tasks deleted should be 0 %d", us.TasksDeleted)
	}
	if us.PokemonCaught != 0 {
		t.Fatalf("Pokemon caught should be 0 %d", us.PokemonCaught)
	}
	if us.ShinyCaught != 0 {
		t.Fatalf("Shiny caught should be 0 %d", us.ShinyCaught)
	}
	if us.UniquePokemon != 0 {
		t.Fatalf("Unique pokemon should be 0 %d", us.UniquePokemon)
	}
	if us.CurrentStreak != 0 {
		t.Fatalf("Current streak should be 0 %d", us.CurrentStreak)
	}
	if us.LongestStreak != 0 {
		t.Fatalf("Longest streak should be 0 %d", us.LongestStreak)
	}
}

func contains(ids []int, id int) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}

func createTestTask(t *testing.T, service *Service, input Task) Task {
	t.Helper()

	task, err := service.CreateTask(input)
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	return task
}

func createTestTasksForUser(t *testing.T, service *Service, userID string, count int) []Task {
	t.Helper()

	tasks := make([]Task, 0, count)

	for i := 0; i < count; i++ {
		task := createTestTask(t, service, Task{
			UserID:      userID,
			Title:       "Test task",
			Description: "Test task description",
			Tag:         "test",
		})

		tasks = append(tasks, task)
	}

	return tasks
}

func getRewardByTaskID(t *testing.T, db *sql.DB, taskID int) TaskReward {
	t.Helper()

	var reward TaskReward

	err := db.QueryRow(`
		SELECT
			id,
			task_id,
			pokemon_id,
			pokemon_name,
			sprite,
			rarity,
			shiny,
			revealed,
			generated_at,
			revealed_at
		FROM task_rewards
		WHERE task_id = ?
	`, taskID).Scan(
		&reward.ID,
		&reward.TaskID,
		&reward.PokemonID,
		&reward.PokemonName,
		&reward.Sprite,
		&reward.Rarity,
		&reward.Shiny,
		&reward.Revealed,
		&reward.GeneratedAt,
		&reward.RevealedAt,
	)

	if err != nil {
		t.Fatalf("get reward by task ID %d: %v", taskID, err)
	}

	return reward
}

func setRewardForTask(
	t *testing.T,
	db *sql.DB,
	taskID int,
	pokemonID int,
	pokemonName string,
	shiny bool,
) TaskReward {
	t.Helper()

	_, err := db.Exec(`
		UPDATE task_rewards
		SET
			pokemon_id = ?,
			pokemon_name = ?,
			sprite = ?,
			rarity = ?,
			shiny = ?,
			revealed = 0,
			revealed_at = NULL
		WHERE task_id = ?
	`,
		pokemonID,
		pokemonName,
		"https://example.com/"+pokemonName+".png",
		1,
		shiny,
		taskID,
	)
	if err != nil {
		t.Fatalf("set reward for task %d: %v", taskID, err)
	}

	return getRewardByTaskID(t, db, taskID)
}

func findCollectionEntry(
	t *testing.T,
	collection []CollectionEntry,
	pokemonID int,
	shiny bool,
) CollectionEntry {
	t.Helper()

	for _, entry := range collection {
		if entry.PokemonID == pokemonID && entry.Shiny == shiny {
			return entry
		}
	}

	t.Fatalf("collection entry not found: pokemonID=%d shiny=%v", pokemonID, shiny)

	return CollectionEntry{}
}

func assertCompletedTask(t *testing.T, task Task) {
	t.Helper()

	if !task.Completed {
		t.Fatal("expected task to be completed")
	}

	if task.CompletedAt == nil {
		t.Fatal("expected CompletedAt to be set")
	}
}

func assertOpenTask(t *testing.T, task Task) {
	t.Helper()

	if task.Completed {
		t.Fatal("expected task to be open")
	}

	if task.CompletedAt != nil {
		t.Fatal("expected CompletedAt to be nil")
	}
}

func assertRevealedReward(t *testing.T, reward TaskReward) {
	t.Helper()

	if !reward.Revealed {
		t.Fatal("expected reward to be revealed")
	}

	if reward.RevealedAt == nil {
		t.Fatal("expected RevealedAt to be set")
	}

	if reward.PokemonID == 0 {
		t.Fatal("expected PokemonID to be set")
	}

	if reward.PokemonName == "" {
		t.Fatal("expected PokemonName to be set")
	}
}

func assertStatisticCounts(
	t *testing.T,
	got UserStatistic,
	tasksOpened int,
	tasksCompleted int,
	tasksDeleted int,
	pokemonCaught int,
) {
	t.Helper()

	if got.TasksOpened != tasksOpened {
		t.Fatalf("TasksOpened mismatch: got %d, expect %d", got.TasksOpened, tasksOpened)
	}

	if got.TasksCompleted != tasksCompleted {
		t.Fatalf("TasksCompleted mismatch: got %d, expect %d", got.TasksCompleted, tasksCompleted)
	}

	if got.TasksDeleted != tasksDeleted {
		t.Fatalf("TasksDeleted mismatch: got %d, expect %d", got.TasksDeleted, tasksDeleted)
	}

	if got.PokemonCaught != pokemonCaught {
		t.Fatalf("PokemonCaught mismatch: got %d, expect %d", got.PokemonCaught, pokemonCaught)
	}
}
