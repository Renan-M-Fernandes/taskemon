package task

import (
	"database/sql"
	"errors"
	"testing"
	"time"
)

func setupRepository(t *testing.T) (*Repository, *sql.DB) {
	t.Helper()

	db := setupTestDB(t)
	repo := NewRepository(db)

	return repo, db
}

func createRepoTask(t *testing.T, repo *Repository, input Task) Task {
	t.Helper()

	task, err := repo.CreateTask(input)
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	return task
}

func createRepoReward(t *testing.T, repo *Repository, taskID int, pokemonID int, pokemonName string, shiny bool) TaskReward {
	t.Helper()
	err := repo.CreateTaskReward(TaskReward{
		TaskID:      taskID,
		PokemonID:   pokemonID,
		PokemonName: pokemonName,
		Sprite:      "https://example.com/" + pokemonName + ".png",
		Rarity:      1,
		Shiny:       shiny,
	})
	if err != nil {
		t.Fatalf("CreateTaskReward: %v", err)
	}

	reward, err := repo.GetTaskReward(taskID)
	if err != nil {
		t.Fatalf("GetTaskReward: %v", err)
	}

	return reward
}

func createRepoTaskWithReward(t *testing.T, repo *Repository, input Task, pokemonID int, pokemonName string, shiny bool) Task {
	t.Helper()

	task := createRepoTask(t, repo, input)
	task.Reward = createRepoReward(t, repo, task.ID, pokemonID, pokemonName, shiny)

	return task
}

func createRepoStatistic(t *testing.T, repo *Repository, input UserStatistic) UserStatistic {
	t.Helper()

	err := repo.CreateUserStatistic(input)
	if err != nil {
		t.Fatalf("CreateUserStatistic: %v", err)
	}

	stat, err := repo.GetStatistic(input.UserID)
	if err != nil {
		t.Fatalf("GetStatistic: %v", err)
	}

	return stat
}

func TestRepositoryCreateTask(t *testing.T) {
	repo, _ := setupRepository(t)

	dueAt := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	input := Task{
		UserID:      "ash",
		Title:       "Catch Pikachu",
		Description: "Find a Pikachu in Viridian Forest",
		Tag:         "pokemon",
		DueAt:       &dueAt,
	}

	created, err := repo.CreateTask(input)
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	if created.ID == 0 {
		t.Fatal("expected ID to be set")
	}

	if created.UserID != input.UserID {
		t.Fatalf("UserID mismatch: got %q, expect %q", created.UserID, input.UserID)
	}

	if created.Title != input.Title {
		t.Fatalf("Title mismatch: got %q, expect %q", created.Title, input.Title)
	}

	if timePtrString(created.DueAt) != timePtrString(input.DueAt) {
		t.Fatalf("DueAt mismatch: got %s, expect %s", timePtrString(created.DueAt), timePtrString(input.DueAt))
	}
}

func TestRepositoryGetTask(t *testing.T) {
	repo, _ := setupRepository(t)

	expected := createRepoTaskWithReward(t, repo, Task{
		UserID:      "ash",
		Title:       "Catch Pikachu",
		Description: "Find a Pikachu in Viridian Forest",
		Tag:         "pokemon",
	}, 25, "pikachu", false)

	got, err := repo.GetTask(expected.ID, expected.UserID)
	if err != nil {
		t.Fatalf("GetTask: %v", err)
	}

	assertTaskMatchesInput(t, got, expected, false)
	assertRewardMatchesInput(t, got.Reward, expected.Reward, false)
}

func TestRepositoryGetTaskNotFound(t *testing.T) {
	repo, _ := setupRepository(t)

	_, err := repo.GetTask(999, "ash")
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected %v, got %v", ErrTaskNotFound, err)
	}
}

func TestRepositoryGetTaskWithoutReward(t *testing.T) {
	repo, _ := setupRepository(t)

	task := createRepoTask(t, repo, Task{
		UserID: "ash",
		Title:  "Catch Pikachu",
		Tag:    "pokemon",
	})

	_, err := repo.GetTask(task.ID, task.UserID)
	if !errors.Is(err, ErrTaskRewardNotFound) {
		t.Fatalf("expected %v, got %v", ErrTaskRewardNotFound, err)
	}
}

func TestRepositoryExistTask(t *testing.T) {
	repo, _ := setupRepository(t)

	expected := createRepoTaskWithReward(t, repo, Task{
		UserID: "ash",
		Title:  "Catch Pikachu",
		Tag:    "pokemon",
	}, 25, "pikachu", false)

	got, err := repo.ExistTask(expected.ID, expected.UserID)
	if err != nil {
		t.Fatalf("ExistTask: %v", err)
	}

	if got.ID != expected.ID {
		t.Fatalf("ID mismatch: got %d, expect %d", got.ID, expected.ID)
	}

	if got.UserID != expected.UserID {
		t.Fatalf("UserID mismatch: got %q, expect %q", got.UserID, expected.UserID)
	}

	if got.Reward.ID != expected.Reward.ID {
		t.Fatalf("Reward ID mismatch: got %d, expect %d", got.Reward.ID, expected.Reward.ID)
	}
}

func TestRepositoryExistTaskNotFound(t *testing.T) {
	repo, _ := setupRepository(t)

	_, err := repo.ExistTask(999, "ash")
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected %v, got %v", ErrTaskNotFound, err)
	}
}

func TestRepositoryUpdateTask(t *testing.T) {
	repo, _ := setupRepository(t)

	task := createRepoTaskWithReward(t, repo, Task{
		UserID:      "ash",
		Title:       "Catch Pikachu",
		Description: "Original description",
		Tag:         "pokemon",
	}, 25, "pikachu", false)

	dueAt := time.Now().Add(48 * time.Hour).Truncate(time.Second)
	task.Title = "Catch Raichu"
	task.Description = "Updated description"
	task.Tag = "electric"
	task.DueAt = &dueAt

	err := repo.UpdateTask(task)
	if err != nil {
		t.Fatalf("UpdateTask: %v", err)
	}

	got, err := repo.GetTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("GetTask: %v", err)
	}

	assertTaskMatchesInput(t, got, task, false)
}

func TestRepositoryUpdateTaskWrongUser(t *testing.T) {
	repo, _ := setupRepository(t)

	task := createRepoTaskWithReward(t, repo, Task{
		UserID:      "ash",
		Title:       "Catch Pikachu",
		Description: "Original description",
		Tag:         "pokemon",
	}, 25, "pikachu", false)

	task.UserID = "misty"
	task.Title = "Steal Pikachu"

	err := repo.UpdateTask(task)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected %v, got %v", ErrTaskNotFound, err)
	}
}

func TestRepositoryUpdateTaskNotFound(t *testing.T) {
	repo, _ := setupRepository(t)

	err := repo.UpdateTask(Task{
		ID:     999,
		UserID: "ash",
		Title:  "Missing task",
		Tag:    "pokemon",
	})

	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected %v, got %v", ErrTaskNotFound, err)
	}
}

func TestRepositoryUpdateCompletedTask(t *testing.T) {
	repo, _ := setupRepository(t)

	task := createRepoTaskWithReward(t, repo, Task{
		UserID: "ash",
		Title:  "Catch Pikachu",
		Tag:    "pokemon",
	}, 25, "pikachu", false)

	err := repo.CompleteTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("CompleteTask: %v", err)
	}

	task.Title = "Catch Raichu"
	err = repo.UpdateTask(task)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected %v, got %v", ErrTaskNotFound, err)
	}
}

func TestRepositoryCompleteTask(t *testing.T) {
	repo, _ := setupRepository(t)

	task := createRepoTaskWithReward(t, repo, Task{
		UserID: "ash",
		Title:  "Catch Pikachu",
		Tag:    "pokemon",
	}, 25, "pikachu", false)

	err := repo.CompleteTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("CompleteTask: %v", err)
	}

	got, err := repo.GetTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("GetTask: %v", err)
	}

	assertCompletedTask(t, got)
}

func TestRepositoryCompleteTaskNotFound(t *testing.T) {
	repo, _ := setupRepository(t)

	err := repo.CompleteTask(999, "ash")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRepositoryDeleteTask(t *testing.T) {
	repo, _ := setupRepository(t)

	task := createRepoTask(t, repo, Task{
		UserID: "ash",
		Title:  "Catch Pikachu",
		Tag:    "pokemon",
	})

	err := repo.DeleteTask(task.ID)
	if err != nil {
		t.Fatalf("DeleteTask: %v", err)
	}

	_, err = repo.GetTask(task.ID, task.UserID)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected %v, got %v", ErrTaskNotFound, err)
	}
}

func TestRepositoryDeleteTaskNotFound(t *testing.T) {
	repo, _ := setupRepository(t)

	err := repo.DeleteTask(999)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected %v, got %v", ErrTaskNotFound, err)
	}
}

func TestRepositoryListTasksByUser(t *testing.T) {
	repo, _ := setupRepository(t)

	createRepoTaskWithReward(t, repo, Task{UserID: "ash", Title: "Catch Pikachu", Tag: "pokemon"}, 25, "pikachu", false)
	createRepoTaskWithReward(t, repo, Task{UserID: "ash", Title: "Catch Bulbasaur", Tag: "pokemon"}, 1, "bulbasaur", false)
	createRepoTaskWithReward(t, repo, Task{UserID: "misty", Title: "Catch Psyduck", Tag: "pokemon"}, 54, "psyduck", false)

	tasks, err := repo.ListTasksByUser("ash")
	if err != nil {
		t.Fatalf("ListTasksByUser: %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("task count mismatch: got %d, expect 2", len(tasks))
	}

	for _, task := range tasks {
		if task.UserID != "ash" {
			t.Fatalf("unexpected UserID: got %q, expect %q", task.UserID, "ash")
		}
	}
}

func TestRepositoryListTasksByUserNotCompleted(t *testing.T) {
	repo, _ := setupRepository(t)

	openTask := createRepoTaskWithReward(t, repo, Task{UserID: "ash", Title: "Catch Pikachu", Tag: "pokemon"}, 25, "pikachu", false)
	completedTask := createRepoTaskWithReward(t, repo, Task{UserID: "ash", Title: "Catch Bulbasaur", Tag: "pokemon"}, 1, "bulbasaur", false)

	err := repo.CompleteTask(completedTask.ID, completedTask.UserID)
	if err != nil {
		t.Fatalf("CompleteTask: %v", err)
	}

	tasks, err := repo.ListTasksByUserNotCompleted("ash")
	if err != nil {
		t.Fatalf("ListTasksByUserNotCompleted: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("task count mismatch: got %d, expect 1", len(tasks))
	}

	if tasks[0].ID != openTask.ID {
		t.Fatalf("ID mismatch: got %d, expect %d", tasks[0].ID, openTask.ID)
	}

	assertOpenTask(t, tasks[0])
}

func TestRepositoryListTasksByUserCompleted(t *testing.T) {
	repo, _ := setupRepository(t)

	openTask := createRepoTaskWithReward(t, repo, Task{UserID: "ash", Title: "Catch Pikachu", Tag: "pokemon"}, 25, "pikachu", false)
	completedTask := createRepoTaskWithReward(t, repo, Task{UserID: "ash", Title: "Catch Bulbasaur", Tag: "pokemon"}, 1, "bulbasaur", false)
	_ = openTask

	err := repo.CompleteTask(completedTask.ID, completedTask.UserID)
	if err != nil {
		t.Fatalf("CompleteTask: %v", err)
	}

	tasks, err := repo.ListTasksByUserCompleted("ash")
	if err != nil {
		t.Fatalf("ListTasksByUserCompleted: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("task count mismatch: got %d, expect 1", len(tasks))
	}

	if tasks[0].ID != completedTask.ID {
		t.Fatalf("ID mismatch: got %d, expect %d", tasks[0].ID, completedTask.ID)
	}

	assertCompletedTask(t, tasks[0])
}

func TestRepositoryCreateTaskReward(t *testing.T) {
	repo, _ := setupRepository(t)

	task := createRepoTask(t, repo, Task{UserID: "ash", Title: "Catch Pikachu", Tag: "pokemon"})
	reward := createRepoReward(t, repo, task.ID, 25, "pikachu", false)

	assertNewReward(t, reward, task.ID)
}

func TestRepositoryCreateTaskRewardRequiresTask(t *testing.T) {
	repo, _ := setupRepository(t)

	err := repo.CreateTaskReward(TaskReward{
		TaskID:      999,
		PokemonID:   25,
		PokemonName: "pikachu",
		Sprite:      "https://example.com/pikachu.png",
		Rarity:      1,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRepositoryGetTaskRewardNotFound(t *testing.T) {
	repo, _ := setupRepository(t)

	_, err := repo.GetTaskReward(999)
	if !errors.Is(err, ErrTaskRewardNotFound) {
		t.Fatalf("expected %v, got %v", ErrTaskRewardNotFound, err)
	}
}

func TestRepositoryRevealPokemon(t *testing.T) {
	repo, _ := setupRepository(t)

	task := createRepoTaskWithReward(t, repo, Task{UserID: "ash", Title: "Catch Pikachu", Tag: "pokemon"}, 25, "pikachu", false)

	err := repo.RevealPokemon(task.ID)
	if err != nil {
		t.Fatalf("RevealPokemon: %v", err)
	}

	reward, err := repo.GetTaskReward(task.ID)
	if err != nil {
		t.Fatalf("GetTaskReward: %v", err)
	}

	assertRevealedReward(t, reward)
}

func TestRepositoryRevealPokemonNotFound(t *testing.T) {
	repo, _ := setupRepository(t)

	err := repo.RevealPokemon(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRepositoryDeleteTaskReward(t *testing.T) {
	repo, _ := setupRepository(t)

	task := createRepoTaskWithReward(t, repo, Task{UserID: "ash", Title: "Catch Pikachu", Tag: "pokemon"}, 25, "pikachu", false)

	err := repo.DeleteTaskReward(task.ID)
	if err != nil {
		t.Fatalf("DeleteTaskReward: %v", err)
	}

	_, err = repo.GetTaskReward(task.ID)
	if !errors.Is(err, ErrTaskRewardNotFound) {
		t.Fatalf("expected %v, got %v", ErrTaskRewardNotFound, err)
	}
}

func TestRepositoryDeleteTaskRewardNotFound(t *testing.T) {
	repo, _ := setupRepository(t)

	err := repo.DeleteTaskReward(999)
	if !errors.Is(err, ErrTaskRewardNotFound) {
		t.Fatalf("expected %v, got %v", ErrTaskRewardNotFound, err)
	}
}

func TestRepositoryListRevealedPokemons(t *testing.T) {
	repo, _ := setupRepository(t)

	first := createRepoTaskWithReward(t, repo, Task{UserID: "ash", Title: "Catch Pikachu", Tag: "pokemon"}, 25, "pikachu", false)
	second := createRepoTaskWithReward(t, repo, Task{UserID: "ash", Title: "Catch Bulbasaur", Tag: "pokemon"}, 1, "bulbasaur", false)
	_ = second

	err := repo.RevealPokemon(first.ID)
	if err != nil {
		t.Fatalf("RevealPokemon: %v", err)
	}

	rewards, err := repo.ListRevealedPokemons()
	if err != nil {
		t.Fatalf("ListRevealedPokemons: %v", err)
	}

	if len(rewards) != 1 {
		t.Fatalf("reward count mismatch: got %d, expect 1", len(rewards))
	}

	if rewards[0].TaskID != first.ID {
		t.Fatalf("TaskID mismatch: got %d, expect %d", rewards[0].TaskID, first.ID)
	}
}

func TestRepositoryCreateCollectionEntry(t *testing.T) {
	repo, _ := setupRepository(t)

	input := CollectionEntry{
		UserID:      "ash",
		PokemonID:   25,
		PokemonName: "pikachu",
		Rarity:      1,
		Shiny:       false,
	}

	err := repo.CreateCollectionEntry(input)
	if err != nil {
		t.Fatalf("CreateCollectionEntry: %v", err)
	}

	collection, err := repo.ListCollection(input.UserID)
	if err != nil {
		t.Fatalf("ListCollection: %v", err)
	}

	if len(collection) != 1 {
		t.Fatalf("collection count mismatch: got %d, expect 1", len(collection))
	}

	entry := collection[0]
	if entry.PokemonID != input.PokemonID {
		t.Fatalf("PokemonID mismatch: got %d, expect %d", entry.PokemonID, input.PokemonID)
	}

	if entry.Count != 1 {
		t.Fatalf("Count mismatch: got %d, expect 1", entry.Count)
	}
}

func TestRepositoryExistCollectionEntry(t *testing.T) {
	repo, _ := setupRepository(t)

	err := repo.CreateCollectionEntry(CollectionEntry{
		UserID:      "ash",
		PokemonID:   25,
		PokemonName: "pikachu",
		Rarity:      1,
	})
	if err != nil {
		t.Fatalf("CreateCollectionEntry: %v", err)
	}

	id, err := repo.ExistCollectionEntry("ash", 25, false)
	if err != nil {
		t.Fatalf("ExistCollectionEntry: %v", err)
	}

	if id == 0 {
		t.Fatal("expected collection entry ID to be set")
	}
}

func TestRepositoryExistCollectionEntryNotFound(t *testing.T) {
	repo, _ := setupRepository(t)

	id, err := repo.ExistCollectionEntry("ash", 25, false)
	if err != nil {
		t.Fatalf("ExistCollectionEntry: %v", err)
	}

	if id != 0 {
		t.Fatalf("expected ID 0, got %d", id)
	}
}

func TestRepositoryUpdateCollectionEntry(t *testing.T) {
	repo, _ := setupRepository(t)

	err := repo.CreateCollectionEntry(CollectionEntry{
		UserID:      "ash",
		PokemonID:   25,
		PokemonName: "pikachu",
		Rarity:      1,
		Shiny:       false,
	})
	if err != nil {
		t.Fatalf("CreateCollectionEntry: %v", err)
	}

	err = repo.UpdateCollectionEntry(25, false, "ash")
	if err != nil {
		t.Fatalf("UpdateCollectionEntry: %v", err)
	}

	collection, err := repo.ListCollection("ash")
	if err != nil {
		t.Fatalf("ListCollection: %v", err)
	}

	entry := findCollectionEntry(t, collection, 25, false)
	if entry.Count != 2 {
		t.Fatalf("Count mismatch: got %d, expect 2", entry.Count)
	}
}

func TestRepositoryUpdateCollectionEntryNotFound(t *testing.T) {
	repo, _ := setupRepository(t)

	err := repo.UpdateCollectionEntry(25, false, "ash")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRepositoryListCollection(t *testing.T) {
	repo, _ := setupRepository(t)

	entries := []CollectionEntry{
		{UserID: "ash", PokemonID: 25, PokemonName: "pikachu", Rarity: 1},
		{UserID: "ash", PokemonID: 4, PokemonName: "charmander", Rarity: 1},
		{UserID: "misty", PokemonID: 54, PokemonName: "psyduck", Rarity: 1},
	}

	for _, entry := range entries {
		if err := repo.CreateCollectionEntry(entry); err != nil {
			t.Fatalf("CreateCollectionEntry: %v", err)
		}
	}

	collection, err := repo.ListCollection("ash")
	if err != nil {
		t.Fatalf("ListCollection: %v", err)
	}

	if len(collection) != 2 {
		t.Fatalf("collection count mismatch: got %d, expect 2", len(collection))
	}
}

func TestRepositoryCreateUserStatistic(t *testing.T) {
	repo, _ := setupRepository(t)

	input := UserStatistic{
		UserID:         "ash",
		TasksOpened:    3,
		TasksCompleted: 1,
		TasksDeleted:   1,
		PokemonCaught:  1,
		ShinyCaught:    1,
		UniquePokemon:  1,
		CurrentStreak:  1,
		LongestStreak:  2,
	}

	got := createRepoStatistic(t, repo, input)

	if got != input {
		t.Fatalf("statistic mismatch: got %+v, expect %+v", got, input)
	}
}

func TestRepositoryGetStatisticNotFound(t *testing.T) {
	repo, _ := setupRepository(t)

	_, err := repo.GetStatistic("missing-user")
	if !errors.Is(err, ErrStatisticNotFound) {
		t.Fatalf("expected %v, got %v", ErrStatisticNotFound, err)
	}
}

func TestRepositoryExistStatistic(t *testing.T) {
	repo, _ := setupRepository(t)

	input := createRepoStatistic(t, repo, UserStatistic{
		UserID:      "ash",
		TasksOpened: 3,
	})

	got, err := repo.ExistStatistic(input.UserID)
	if err != nil {
		t.Fatalf("ExistStatistic: %v", err)
	}

	if got.UserID != input.UserID {
		t.Fatalf("UserID mismatch: got %q, expect %q", got.UserID, input.UserID)
	}

	if got.TasksOpened != input.TasksOpened {
		t.Fatalf("TasksOpened mismatch: got %d, expect %d", got.TasksOpened, input.TasksOpened)
	}
}

func TestRepositoryExistStatisticNotFound(t *testing.T) {
	repo, _ := setupRepository(t)

	got, err := repo.ExistStatistic("missing-user")
	if err != nil {
		t.Fatalf("ExistStatistic: %v", err)
	}

	if got.UserID != "" {
		t.Fatalf("expected empty UserID, got %q", got.UserID)
	}
}

func TestRepositoryUpdateUserStatisticOnCreate(t *testing.T) {
	repo, _ := setupRepository(t)

	_ = createRepoStatistic(t, repo, UserStatistic{UserID: "ash"})

	err := repo.UpdateUserStatisticOnCreate("ash")
	if err != nil {
		t.Fatalf("UpdateUserStatisticOnCreate: %v", err)
	}

	got, err := repo.GetStatistic("ash")
	if err != nil {
		t.Fatalf("GetStatistic: %v", err)
	}

	if got.TasksOpened != 1 {
		t.Fatalf("TasksOpened mismatch: got %d, expect 1", got.TasksOpened)
	}
}

func TestRepositoryUpdateUserStatisticOnCreateNotFound(t *testing.T) {
	repo, _ := setupRepository(t)

	err := repo.UpdateUserStatisticOnCreate("missing-user")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRepositoryUpdateUserStatisticOnDelete(t *testing.T) {
	repo, _ := setupRepository(t)

	_ = createRepoStatistic(t, repo, UserStatistic{UserID: "ash"})

	err := repo.UpdateUserStatisticOnDelete("ash")
	if err != nil {
		t.Fatalf("UpdateUserStatisticOnDelete: %v", err)
	}

	got, err := repo.GetStatistic("ash")
	if err != nil {
		t.Fatalf("GetStatistic: %v", err)
	}

	if got.TasksDeleted != 1 {
		t.Fatalf("TasksDeleted mismatch: got %d, expect 1", got.TasksDeleted)
	}
}

func TestRepositoryUpdateUserStatisticOnClose(t *testing.T) {
	repo, _ := setupRepository(t)

	_ = createRepoStatistic(t, repo, UserStatistic{UserID: "ash"})

	err := repo.UpdateUserStatisticOnClose("ash", 1, 2, 3, 4)
	if err != nil {
		t.Fatalf("UpdateUserStatisticOnClose: %v", err)
	}

	got, err := repo.GetStatistic("ash")
	if err != nil {
		t.Fatalf("GetStatistic: %v", err)
	}

	if got.TasksCompleted != 1 {
		t.Fatalf("TasksCompleted mismatch: got %d, expect 1", got.TasksCompleted)
	}

	if got.PokemonCaught != 1 {
		t.Fatalf("PokemonCaught mismatch: got %d, expect 1", got.PokemonCaught)
	}

	if got.ShinyCaught != 1 {
		t.Fatalf("ShinyCaught mismatch: got %d, expect 1", got.ShinyCaught)
	}

	if got.UniquePokemon != 2 {
		t.Fatalf("UniquePokemon mismatch: got %d, expect 2", got.UniquePokemon)
	}

	if got.CurrentStreak != 3 {
		t.Fatalf("CurrentStreak mismatch: got %d, expect 3", got.CurrentStreak)
	}

	if got.LongestStreak != 4 {
		t.Fatalf("LongestStreak mismatch: got %d, expect 4", got.LongestStreak)
	}
}

func TestRepositoryGetDataForStatistic(t *testing.T) {
	repo, _ := setupRepository(t)

	_ = createRepoStatistic(t, repo, UserStatistic{
		UserID:        "ash",
		LongestStreak: 2,
	})

	first := createRepoTaskWithReward(t, repo, Task{UserID: "ash", Title: "Catch Pikachu", Tag: "pokemon"}, 25, "pikachu", true)
	second := createRepoTaskWithReward(t, repo, Task{UserID: "ash", Title: "Catch Bulbasaur", Tag: "pokemon"}, 1, "bulbasaur", false)
	otherUser := createRepoTaskWithReward(t, repo, Task{UserID: "misty", Title: "Catch Psyduck", Tag: "pokemon"}, 54, "psyduck", true)
	_ = otherUser

	if err := repo.CompleteTask(first.ID, first.UserID); err != nil {
		t.Fatalf("CompleteTask first: %v", err)
	}
	if err := repo.CompleteTask(second.ID, second.UserID); err != nil {
		t.Fatalf("CompleteTask second: %v", err)
	}
	if err := repo.CompleteTask(otherUser.ID, otherUser.UserID); err != nil {
		t.Fatalf("CompleteTask other user: %v", err)
	}

	if err := repo.RevealPokemon(first.ID); err != nil {
		t.Fatalf("CompleteTaskReward first: %v", err)
	}
	if err := repo.RevealPokemon(second.ID); err != nil {
		t.Fatalf("CompleteTaskReward second: %v", err)
	}
	if err := repo.RevealPokemon(otherUser.ID); err != nil {
		t.Fatalf("CompleteTaskReward other user: %v", err)
	}

	if err := repo.CreateCollectionEntry(CollectionEntry{UserID: "ash", PokemonID: 25, PokemonName: "pikachu", Rarity: 1, Shiny: true}); err != nil {
		t.Fatalf("CreateCollectionEntry first: %v", err)
	}
	if err := repo.CreateCollectionEntry(CollectionEntry{UserID: "ash", PokemonID: 1, PokemonName: "bulbasaur", Rarity: 1}); err != nil {
		t.Fatalf("CreateCollectionEntry second: %v", err)
	}
	if err := repo.CreateCollectionEntry(CollectionEntry{UserID: "misty", PokemonID: 54, PokemonName: "psyduck", Rarity: 1, Shiny: true}); err != nil {
		t.Fatalf("CreateCollectionEntry other user: %v", err)
	}

	longestStreak, shinyCaughtTotal, uniquePokemonTotal, dates, err := repo.GetDataForStatistic("ash")
	if err != nil {
		t.Fatalf("GetDataForStatistic: %v", err)
	}

	if longestStreak != 2 {
		t.Fatalf("longestStreak mismatch: got %d, expect 2", longestStreak)
	}

	if shinyCaughtTotal != 1 {
		t.Fatalf("shinyCaughtTotal mismatch: got %d, expect 1", shinyCaughtTotal)
	}

	if uniquePokemonTotal != 2 {
		t.Fatalf("uniquePokemonTotal mismatch: got %d, expect 2", uniquePokemonTotal)
	}

	if len(dates) == 0 {
		t.Fatal("expected at least one completion date")
	}
}
