package task

import (
	"errors"
	"testing"
	"time"
)

func TestCreateTask(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	tests := []struct {
		name    string
		input   Task
		wantErr error
	}{
		{
			name: "Create task with required fields",
			input: Task{
				UserID: "ash",
				Title:  "Catch Pikachu",
			},
			wantErr: nil,
		},
		{
			name: "Create task with all fields",
			input: Task{
				UserID:      "ash",
				Title:       "Catch Pikachu",
				Description: "Find and catch a Pikachu in Viridian Forest",
				Tag:         "pokemon",
				DueAt:       ptrTime(time.Now()),
			},
			wantErr: nil,
		},
		{
			name: "Create task with empty title",
			input: Task{
				UserID:      "ash",
				Description: "Find and catch a Pikachu in Viridian Forest",
				Tag:         "pokemon",
				DueAt:       ptrTime(time.Now()),
			},
			wantErr: ErrEmptyTitle,
		},
		{
			name: "Create task with empty tag",
			input: Task{
				UserID:      "ash",
				Title:       "Catch Pikachu",
				Description: "Find and catch a Pikachu in Viridian Forest",
				Tag:         "pokemon",
				DueAt:       ptrTime(time.Now()),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			task, err := service.CreateTask(tt.input)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Create task: expected %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Create task start: returned error: %v", err)
			}

			assertTaskMatchesInput(t, task, tt.input, false)
		})
	}
}

func TestCreateTaskCreatesHiddenReward(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)
	task := createBulkTasks(t, service, 50)

	for _, tt := range task {
		t.Run("Teste create task hidden reward", func(t *testing.T) {
			assertNewReward(t, tt.Reward, tt.ID)
		})
	}
}

func TestCreateTaskCreatesUserStatisticIfMissing(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)
	task := createBulkTasks(t, service, 50)
	userTotals := getUserCount(t, task)

	for _, c := range userTotals {
		us, err := service.GetStatistic(c.UserID)
		if err != nil {
			t.Fatalf("Get user statistic at: %q", c.UserID)
		}

		if us.TasksOpened != c.Total {
			t.Fatalf("Tasks opened mismatch for %q: got %q, expect %q", c.UserID, us.TasksOpened, c.UserID)
		}
		if us.UserID != c.UserID {
			t.Fatalf("Userid mismatch: got %q, expect %q", c.UserID, us.UserID)
		}

		assertNewStatistics(t, us)
	}
}

func TestGetTask(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)
	task := createBulkTasks(t, service, 50)
	completeList := completeBulkTasks(t, service, task, 20)

	for _, tt := range task {
		t.Run("Test get task", func(t *testing.T) {
			task, err := service.GetTask(tt.ID, tt.UserID)
			if err != nil {
				t.Fatalf("Get task (%v) returned error: %v", tt.ID, err)
			}
			completed := contains(completeList, tt.ID)
			assertTaskMatchesInput(t, task, tt, completed)

			assertRewardMatchesInput(t, task.Reward, tt.Reward, completed)
		})
	}
}

func TestGetTaskNotFound(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	_, err := service.GetTask(1, "renan")

	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("Expected %v, got %v", ErrTaskNotFound, err)
	}
}

func TestListTasks(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)
	task := createBulkTasks(t, service, 50)
	userTotals := getUserCount(t, task)
	_ = completeBulkTasks(t, service, task, 20)

	for _, c := range userTotals {
		userNotFound := true
		tt, err := service.ListTasksByUser(c.UserID)
		if err != nil {
			t.Fatalf("List task at: %q", c.UserID)
		}
		ttotals := getUserCount(t, tt)
		for _, total := range ttotals {
			if total.UserID == c.UserID {
				userNotFound = false
				if total.Total != c.Total {
					t.Fatalf("List task total mismatch for %q: got %q, expect %q", c.UserID, total.Total, c.Total)
				}
			}
		}
		if userNotFound {
			t.Fatalf("List task user not found at: %q", c.UserID)
		}
	}
}

func TestListCompletedTasks(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)
	task := createBulkTasks(t, service, 50)
	_ = completeBulkTasks(t, service, task, 20)
	userTotals := getUserCount(t, task)
	userTotals = getUserCountCompleted(t, service, userTotals)

	for _, c := range userTotals {
		userNotFound := true
		tt, err := service.ListTasksByUserCompleted(c.UserID)
		if err != nil {
			t.Fatalf("List task completed error at: %q", c.UserID)
		}
		ttotals := getUserCount(t, tt)
		for _, total := range ttotals {
			if total.UserID == c.UserID {
				userNotFound = false
				if total.Total != c.Total {
					t.Fatalf("List task completed total mismatch for %q: got %q, expect %q", c.UserID, total.Total, c.Total)
				}
			}
		}
		if userNotFound {
			t.Fatalf("List task completed user not found at: %q", c.UserID)
		}
	}
}

func TestListTasksExcludesCompletedTasks(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)
	task := createBulkTasks(t, service, 50)
	_ = completeBulkTasks(t, service, task, 20)
	userTotals := getUserCount(t, task)
	userTotals = getUserCountExcludesCompleted(t, service, userTotals)

	for _, c := range userTotals {
		userNotFound := true

		tt, err := service.ListTasksByUserNotCompleted(c.UserID)
		if err != nil {
			t.Fatalf("List task excludes completed error at: %q", c.UserID)
		}
		ttotals := getUserCount(t, tt)
		for _, total := range ttotals {
			if total.UserID == c.UserID {
				userNotFound = false
				if total.Total != c.Total {
					t.Fatalf("List task excludes completed total mismatch for %q: got %q, expect %q", c.UserID, total.Total, c.Total)
				}
			}
		}
		if userNotFound {
			if c.Total != 0 {
				t.Fatalf("List task excludes completed user not found at: %q", c.UserID)
			}
		}
	}
}

//

func TestUpdateTask(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	task := createTestTask(t, service, Task{
		UserID:      "ash",
		Title:       "Catch Pikachu",
		Description: "Original description",
		Tag:         "pokemon",
	})

	newTitle := "Catch Raichu"
	newDescription := "Updated description"
	newTag := "electric"
	newDueAt := time.Now().Add(24 * time.Hour).Truncate(time.Second)

	err := service.UpdateTask(true, TaskUpdate{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       &newTitle,
		Description: &newDescription,
		Tag:         &newTag,
		DueAt:       &newDueAt,
	})
	if err != nil {
		t.Fatalf("UpdateTask: %v", err)
	}

	got, err := service.GetTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("GetTask: %v", err)
	}

	expected := task
	expected.Title = newTitle
	expected.Description = newDescription
	expected.Tag = newTag
	expected.DueAt = &newDueAt

	assertTaskMatchesInput(t, got, expected, false)
}

func TestUpdateTaskNotFound(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	title := "Missing task"

	err := service.UpdateTask(false, TaskUpdate{
		ID:     999,
		UserID: "ash",
		Title:  &title,
	})

	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected %v, got %v", ErrTaskNotFound, err)
	}
}

func TestUpdateTaskKeepsCreatedAt(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	task := createTestTask(t, service, Task{
		UserID: "ash",
		Title:  "Catch Pikachu",
		Tag:    "pokemon",
	})

	before, err := service.GetTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("GetTask before update: %v", err)
	}

	newTitle := "Catch Bulbasaur"

	err = service.UpdateTask(false, TaskUpdate{
		ID:     task.ID,
		UserID: task.UserID,
		Title:  &newTitle,
	})
	if err != nil {
		t.Fatalf("UpdateTask: %v", err)
	}

	after, err := service.GetTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("GetTask after update: %v", err)
	}

	if !after.CreatedAt.Equal(before.CreatedAt) {
		t.Fatalf("CreatedAt changed: got %s, expect %s", after.CreatedAt, before.CreatedAt)
	}
}

func TestCompleteTask(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	task := createTestTask(t, service, Task{
		UserID: "ash",
		Title:  "Catch Pikachu",
		Tag:    "pokemon",
	})

	completedTask, err := service.CompleteTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("CompleteTask: %v", err)
	}

	assertCompletedTask(t, completedTask)

	storedTask, err := service.GetTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("GetTask: %v", err)
	}

	assertCompletedTask(t, storedTask)
}

func TestCompleteTaskRevealsReward(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	task := createTestTask(t, service, Task{
		UserID: "ash",
		Title:  "Catch Pikachu",
		Tag:    "pokemon",
	})

	_, err := service.CompleteTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("CompleteTask: %v", err)
	}

	got, err := service.GetTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("GetTask: %v", err)
	}

	assertRevealedReward(t, got.Reward)
}

func TestCompleteTaskUpdatesCollection(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	task := createTestTask(t, service, Task{
		UserID: "ash",
		Title:  "Catch Pikachu",
		Tag:    "pokemon",
	})

	task.Reward = setRewardForTask(t, db, task.ID, 25, "pikachu", false)

	_, err := service.CompleteTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("CompleteTask: %v", err)
	}

	collection, err := service.ListCollection(task.UserID)
	if err != nil {
		t.Fatalf("ListCollection: %v", err)
	}

	entry := findCollectionEntry(t, collection, 25, false)

	if entry.Count != 1 {
		t.Fatalf("collection count mismatch: got %d, expect 1", entry.Count)
	}
}

func TestCompleteTaskUpdatesStatistics(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	task := createTestTask(t, service, Task{
		UserID: "ash",
		Title:  "Catch Pikachu",
		Tag:    "pokemon",
	})

	_, err := service.CompleteTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("CompleteTask: %v", err)
	}

	stats, err := service.GetStatistic(task.UserID)
	if err != nil {
		t.Fatalf("GetStatistic: %v", err)
	}

	assertStatisticCounts(t, stats, 1, 1, 0, 1)
}

func TestCompleteTaskAlreadyCompleted(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	task := createTestTask(t, service, Task{
		UserID: "ash",
		Title:  "Catch Pikachu",
		Tag:    "pokemon",
	})

	_, err := service.CompleteTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("first CompleteTask: %v", err)
	}

	_, err = service.CompleteTask(task.ID, task.UserID)
	if !errors.Is(err, ErrTaskAlreadyCompleted) {
		t.Fatalf("expected %v, got %v", ErrTaskAlreadyCompleted, err)
	}
}

func TestCompleteTaskNotFound(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	_, err := service.CompleteTask(999, "ash")

	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected %v, got %v", ErrTaskNotFound, err)
	}
}

func TestDeleteTask(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	task := createTestTask(t, service, Task{
		UserID: "ash",
		Title:  "Catch Pikachu",
		Tag:    "pokemon",
	})

	err := service.DeleteTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("DeleteTask: %v", err)
	}

	_, err = service.GetTask(task.ID, task.UserID)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected %v after delete, got %v", ErrTaskNotFound, err)
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	err := service.DeleteTask(999, "ash")

	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected %v, got %v", ErrTaskNotFound, err)
	}
}

func TestDeleteCompletedTask(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	task := createTestTask(t, service, Task{
		UserID: "ash",
		Title:  "Catch Pikachu",
		Tag:    "pokemon",
	})

	_, err := service.CompleteTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("CompleteTask: %v", err)
	}

	err = service.DeleteTask(task.ID, task.UserID)

	if !errors.Is(err, ErrTaskCompletedDelete) {
		t.Fatalf("expected %v, got %v", ErrTaskCompletedDelete, err)
	}
}

func TestDeleteTaskRemovesReward(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	task := createTestTask(t, service, Task{
		UserID: "ash",
		Title:  "Catch Pikachu",
		Tag:    "pokemon",
	})

	err := service.DeleteTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("DeleteTask: %v", err)
	}

	_, err = service.GetTaskReward(task.ID, task.UserID)
	if !errors.Is(err, ErrTaskRewardNotFound) {
		t.Fatalf("expected %v, got %v", ErrTaskRewardNotFound, err)
	}
}

func TestDeleteTaskUpdatesStatistics(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	task := createTestTask(t, service, Task{
		UserID: "ash",
		Title:  "Catch Pikachu",
		Tag:    "pokemon",
	})

	err := service.DeleteTask(task.ID, task.UserID)
	if err != nil {
		t.Fatalf("DeleteTask: %v", err)
	}

	stats, err := service.GetStatistic(task.UserID)
	if err != nil {
		t.Fatalf("GetStatistic: %v", err)
	}

	assertStatisticCounts(t, stats, 1, 0, 1, 0)
}

func TestGetCollection(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	tasks := createTestTasksForUser(t, service, "ash", 3)

	setRewardForTask(t, db, tasks[0].ID, 25, "pikachu", false)
	setRewardForTask(t, db, tasks[1].ID, 4, "charmander", false)
	setRewardForTask(t, db, tasks[2].ID, 1, "bulbasaur", false)

	for _, task := range tasks {
		_, err := service.CompleteTask(task.ID, task.UserID)
		if err != nil {
			t.Fatalf("CompleteTask(%d): %v", task.ID, err)
		}
	}

	collection, err := service.ListCollection("ash")
	if err != nil {
		t.Fatalf("ListCollection: %v", err)
	}

	if len(collection) != 3 {
		t.Fatalf("collection length mismatch: got %d, expect 3", len(collection))
	}
}

func TestGetCollectionEmpty(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	collection, err := service.ListCollection("missing-user")
	if err != nil {
		t.Fatalf("ListCollection: %v", err)
	}

	if len(collection) != 0 {
		t.Fatalf("expected empty collection, got %d entries", len(collection))
	}
}

func TestCollectionIncrementsDuplicatePokemon(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	tasks := createTestTasksForUser(t, service, "ash", 2)

	setRewardForTask(t, db, tasks[0].ID, 25, "pikachu", false)
	setRewardForTask(t, db, tasks[1].ID, 25, "pikachu", false)

	for _, task := range tasks {
		_, err := service.CompleteTask(task.ID, task.UserID)
		if err != nil {
			t.Fatalf("CompleteTask(%d): %v", task.ID, err)
		}
	}

	collection, err := service.ListCollection("ash")
	if err != nil {
		t.Fatalf("ListCollection: %v", err)
	}

	if len(collection) != 1 {
		t.Fatalf("collection length mismatch: got %d, expect 1", len(collection))
	}

	entry := findCollectionEntry(t, collection, 25, false)

	if entry.Count != 2 {
		t.Fatalf("collection count mismatch: got %d, expect 2", entry.Count)
	}
}

func TestGetStatistics(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	_ = createTestTasksForUser(t, service, "ash", 3)

	stats, err := service.GetStatistic("ash")
	if err != nil {
		t.Fatalf("GetStatistic: %v", err)
	}

	assertStatisticCounts(t, stats, 3, 0, 0, 0)
}

func TestGetStatisticsNotFound(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	_, err := service.GetStatistic("missing-user")

	if !errors.Is(err, ErrStatisticNotFound) {
		t.Fatalf("expected %v, got %v", ErrStatisticNotFound, err)
	}
}

func TestStatisticsAfterTaskLifecycle(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	tasks := createTestTasksForUser(t, service, "ash", 3)

	setRewardForTask(t, db, tasks[0].ID, 25, "pikachu", false)
	setRewardForTask(t, db, tasks[1].ID, 4, "charmander", false)

	_, err := service.CompleteTask(tasks[0].ID, tasks[0].UserID)
	if err != nil {
		t.Fatalf("CompleteTask first: %v", err)
	}

	_, err = service.CompleteTask(tasks[1].ID, tasks[1].UserID)
	if err != nil {
		t.Fatalf("CompleteTask second: %v", err)
	}

	err = service.DeleteTask(tasks[2].ID, tasks[2].UserID)
	if err != nil {
		t.Fatalf("DeleteTask: %v", err)
	}

	stats, err := service.GetStatistic("ash")
	if err != nil {
		t.Fatalf("GetStatistic: %v", err)
	}

	assertStatisticCounts(t, stats, 3, 2, 1, 2)

	if stats.UniquePokemon != 2 {
		t.Fatalf("UniquePokemon mismatch: got %d, expect 2", stats.UniquePokemon)
	}
}

func TestCreateTaskValidation(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	tests := []struct {
		name    string
		input   Task
		wantErr error
	}{
		{
			name: "Empty title",
			input: Task{
				UserID: "ash",
				Tag:    "pokemon",
			},
			wantErr: ErrEmptyTitle,
		},
		{
			name: "Empty tag defaults to Misc",
			input: Task{
				UserID: "ash",
				Title:  "Catch Pikachu",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := service.CreateTask(tt.input)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("CreateTask: %v", err)
			}

			assertTaskMatchesInput(t, task, tt.input, false)
		})
	}
}

func TestUpdateTaskValidation(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	task := createTestTask(t, service, Task{
		UserID: "ash",
		Title:  "Catch Pikachu",
		Tag:    "pokemon",
	})

	emptyTitle := ""
	err := service.UpdateTask(false, TaskUpdate{
		ID:     task.ID,
		UserID: task.UserID,
		Title:  &emptyTitle,
	})
	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("expected %v, got %v", ErrEmptyTitle, err)
	}

	emptyTag := ""
	err = service.UpdateTask(false, TaskUpdate{
		ID:     task.ID,
		UserID: task.UserID,
		Tag:    &emptyTag,
	})
	if !errors.Is(err, ErrEmptyTag) {
		t.Fatalf("expected %v, got %v", ErrEmptyTag, err)
	}
}

func TestBulkCreateTasks(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	tasks := createBulkTasks(t, service, 50)

	if len(tasks) != 50 {
		t.Fatalf("created task count mismatch: got %d, expect 50", len(tasks))
	}

	for _, task := range tasks {
		assertTaskMatchesInput(t, task, task, false)
		assertNewReward(t, task.Reward, task.ID)
	}
}

func TestBulkListTasks(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	created := createBulkTasks(t, service, 50)
	userTotals := getUserCount(t, created)

	for _, user := range userTotals {
		tasks, err := service.ListTasksByUser(user.UserID)
		if err != nil {
			t.Fatalf("ListTasksByUser(%s): %v", user.UserID, err)
		}

		if len(tasks) != user.Total {
			t.Fatalf("task count mismatch for %s: got %d, expect %d", user.UserID, len(tasks), user.Total)
		}
	}
}

func TestCollectionAllowsNormalAndShinySamePokemon(t *testing.T) {
	db := setupTestDB(t)
	service := setupService(db)

	tasks := createTestTasksForUser(t, service, "ash", 2)

	setRewardForTask(t, db, tasks[0].ID, 25, "pikachu", false)
	setRewardForTask(t, db, tasks[1].ID, 25, "pikachu", true)

	for _, task := range tasks {
		_, err := service.CompleteTask(task.ID, task.UserID)
		if err != nil {
			t.Fatalf("CompleteTask(%d): %v", task.ID, err)
		}
	}

	collection, err := service.ListCollection("ash")
	if err != nil {
		t.Fatalf("ListCollection: %v", err)
	}

	if len(collection) != 2 {
		t.Fatalf("collection length mismatch: got %d, expect 2", len(collection))
	}

	normal := findCollectionEntry(t, collection, 25, false)
	shiny := findCollectionEntry(t, collection, 25, true)

	if normal.Count != 1 {
		t.Fatalf("normal count mismatch: got %d, expect 1", normal.Count)
	}

	if shiny.Count != 1 {
		t.Fatalf("shiny count mismatch: got %d, expect 1", shiny.Count)
	}
}
