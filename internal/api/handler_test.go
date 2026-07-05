package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/Renan-M-Fernandes/taskemon/internal/task"
)

func performRawRequest(router http.Handler, method string, path string, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return performInvalidRequest(router, req)
}

func TestGetHealthHandler(t *testing.T) {
	server := setupAPITestServer(t)

	rr := performRequest(t, server.router, http.MethodGet, "/api/v1/health", nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("status mismatch: got %d, expect %d", rr.Code, http.StatusOK)
	}

	response := parseJSONResponse[task.Health](t, rr)
	if response.Status == "" {
		t.Fatal("expected health status")
	}
}

func TestGetTaskHandler(t *testing.T) {
	server := setupAPITestServer(t)
	created := createAPITaskWithReward(t, server.repo, task.Task{
		UserID:      "ash",
		Title:       "Catch Pikachu",
		Description: "Find Pikachu",
		Tag:         "pokemon",
	}, 25, "pikachu", false)

	rr := performRequest(t, server.router, http.MethodGet, "/api/v1/tasks/ash/"+strconv.Itoa(created.ID), nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("status mismatch: got %d, expect %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	response := parseJSONResponse[TaskResponse](t, rr)
	if response.ID != created.ID {
		t.Fatalf("ID mismatch: got %d, expect %d", response.ID, created.ID)
	}

	if response.Reward.PokemonName != nil {
		t.Fatal("unrevealed reward should not expose PokemonName")
	}
}

func TestGetTaskHandlerInvalidID(t *testing.T) {
	server := setupAPITestServer(t)

	rr := performRequest(t, server.router, http.MethodGet, "/api/v1/tasks/ash/not-an-id", nil)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status mismatch: got %d, expect %d", rr.Code, http.StatusBadRequest)
	}
}

func TestGetTaskHandlerNotFound(t *testing.T) {
	server := setupAPITestServer(t)

	rr := performRequest(t, server.router, http.MethodGet, "/api/v1/tasks/ash/999", nil)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status mismatch: got %d, expect %d; body=%s", rr.Code, http.StatusNotFound, rr.Body.String())
	}
}

func TestListTasksHandler(t *testing.T) {
	server := setupAPITestServer(t)
	createAPITaskWithReward(t, server.repo, task.Task{UserID: "ash", Title: "Task 1", Tag: "pokemon"}, 25, "pikachu", false)
	createAPITaskWithReward(t, server.repo, task.Task{UserID: "ash", Title: "Task 2", Tag: "pokemon"}, 4, "charmander", false)
	createAPITaskWithReward(t, server.repo, task.Task{UserID: "misty", Title: "Task 3", Tag: "water"}, 54, "psyduck", false)

	rr := performRequest(t, server.router, http.MethodGet, "/api/v1/tasks/ash", nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("status mismatch: got %d, expect %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	response := parseJSONResponse[[]TaskResponse](t, rr)
	if len(response) != 2 {
		t.Fatalf("task length mismatch: got %d, expect 2", len(response))
	}
}

func TestListTasksCompletedHandler(t *testing.T) {
	server := setupAPITestServer(t)
	created := createAPITaskWithReward(t, server.repo, task.Task{UserID: "ash", Title: "Task 1", Tag: "pokemon"}, 25, "pikachu", false)
	if err := server.repo.CompleteTask(created.ID, created.UserID); err != nil {
		t.Fatalf("CompleteTask setup: %v", err)
	}
	created = createAPITaskWithReward(t, server.repo, task.Task{UserID: "ash", Title: "Task 2", Tag: "pokemon"}, 4, "charmander", false)
	if err := server.repo.CompleteTask(created.ID, created.UserID); err != nil {
		t.Fatalf("CompleteTask setup: %v", err)
	}
	createAPITaskWithReward(t, server.repo, task.Task{UserID: "ash", Title: "Task 3", Tag: "pokemon"}, 4, "entei", false)
	createAPITaskWithReward(t, server.repo, task.Task{UserID: "misty", Title: "Task 3", Tag: "water"}, 54, "psyduck", false)

	rr := performRequest(t, server.router, http.MethodGet, "/api/v1/tasks/completed/ash", nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("status mismatch: got %d, expect %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	response := parseJSONResponse[[]TaskResponse](t, rr)
	if len(response) != 2 {
		t.Fatalf("task length mismatch: got %d, expect 2", len(response))
	}
}

func TestListTasksNotCompletedHandler(t *testing.T) {
	server := setupAPITestServer(t)
	created := createAPITaskWithReward(t, server.repo, task.Task{UserID: "ash", Title: "Task 1", Tag: "pokemon"}, 25, "pikachu", false)
	if err := server.repo.CompleteTask(created.ID, created.UserID); err != nil {
		t.Fatalf("CompleteTask setup: %v", err)
	}
	createAPITaskWithReward(t, server.repo, task.Task{UserID: "ash", Title: "Task 2", Tag: "pokemon"}, 4, "charmander", false)
	createAPITaskWithReward(t, server.repo, task.Task{UserID: "ash", Title: "Task 3", Tag: "pokemon"}, 4, "entei", false)
	createAPITaskWithReward(t, server.repo, task.Task{UserID: "misty", Title: "Task 3", Tag: "water"}, 54, "psyduck", false)

	rr := performRequest(t, server.router, http.MethodGet, "/api/v1/tasks/open/ash", nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("status mismatch: got %d, expect %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	response := parseJSONResponse[[]TaskResponse](t, rr)
	if len(response) != 2 {
		t.Fatalf("task length mismatch: got %d, expect 2", len(response))
	}
}

func TestCreateTaskHandlerInvalidJSON(t *testing.T) {
	server := setupAPITestServer(t)

	rr := performRawRequest(server.router, http.MethodPost, "/api/v1/tasks/ash", "{")

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status mismatch: got %d, expect %d", rr.Code, http.StatusBadRequest)
	}
}

func TestCreateTaskHandlerValidationError(t *testing.T) {
	server := setupAPITestServer(t)

	rr := performRequest(t, server.router, http.MethodPost, "/api/v1/tasks/ash", TaskRequest{})

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status mismatch: got %d, expect %d; body=%s", rr.Code, http.StatusBadRequest, rr.Body.String())
	}
}

func TestUpdateTaskHandler(t *testing.T) {
	server := setupAPITestServer(t)
	created := createAPITaskWithReward(t, server.repo, task.Task{UserID: "ash", Title: "Catch Pikachu", Tag: "pokemon"}, 25, "pikachu", false)

	rr := performRequest(t, server.router, http.MethodPatch, "/api/v1/tasks/ash/"+strconv.Itoa(created.ID), TaskRequestUpdate{
		Title: stringPtr("Catch Raichu"),
		Tag:   stringPtr("electric"),
	})

	if rr.Code != http.StatusOK {
		t.Fatalf("status mismatch: got %d, expect %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	got, err := server.repo.GetTask(created.ID, "ash")
	if err != nil {
		t.Fatalf("GetTask: %v", err)
	}

	if got.Title != "Catch Raichu" {
		t.Fatalf("Title mismatch: got %q, expect %q", got.Title, "Catch Raichu")
	}

	if got.Tag != "electric" {
		t.Fatalf("Tag mismatch: got %q, expect %q", got.Tag, "electric")
	}
}

func TestUpdateTaskHandlerInvalidJSON(t *testing.T) {
	server := setupAPITestServer(t)
	created := createAPITaskWithReward(t, server.repo, task.Task{UserID: "ash", Title: "Catch Pikachu", Tag: "pokemon"}, 25, "pikachu", false)

	rr := performRawRequest(server.router, http.MethodPatch, "/api/v1/tasks/ash/"+strconv.Itoa(created.ID), "{")

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status mismatch: got %d, expect %d", rr.Code, http.StatusBadRequest)
	}
}

func TestUpdateTaskHandlerNotFound(t *testing.T) {
	server := setupAPITestServer(t)

	rr := performRequest(t, server.router, http.MethodPatch, "/api/v1/tasks/ash/999", TaskRequestUpdate{Title: stringPtr("Catch Raichu")})

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status mismatch: got %d, expect %d; body=%s", rr.Code, http.StatusNotFound, rr.Body.String())
	}
}

func TestCompleteTaskHandler(t *testing.T) {
	server := setupAPITestServer(t)
	created := createAPITaskWithReward(t, server.repo, task.Task{UserID: "ash", Title: "Catch Pikachu", Tag: "pokemon"}, 25, "pikachu", false)
	createAPIStatistic(t, server.repo, task.UserStatistic{UserID: "ash", TasksOpened: 1})

	rr := performRequest(t, server.router, http.MethodPost, "/api/v1/tasks/ash/"+strconv.Itoa(created.ID)+"/complete", nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("status mismatch: got %d, expect %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	got, err := server.repo.GetTask(created.ID, "ash")
	if err != nil {
		t.Fatalf("GetTask: %v", err)
	}

	if !got.Completed {
		t.Fatal("expected task to be completed")
	}

	if !got.Reward.Revealed {
		t.Fatal("expected reward to be revealed")
	}
}

func TestCompleteTaskHandlerAlreadyCompleted(t *testing.T) {
	server := setupAPITestServer(t)
	created := createAPITaskWithReward(t, server.repo, task.Task{UserID: "ash", Title: "Catch Pikachu", Tag: "pokemon"}, 25, "pikachu", false)
	if err := server.repo.CompleteTask(created.ID, created.UserID); err != nil {
		t.Fatalf("CompleteTask setup: %v", err)
	}

	rr := performRequest(t, server.router, http.MethodPost, "/api/v1/tasks/ash/"+strconv.Itoa(created.ID)+"/complete", nil)

	if rr.Code != http.StatusConflict {
		t.Fatalf("status mismatch: got %d, expect %d; body=%s", rr.Code, http.StatusConflict, rr.Body.String())
	}
}

func TestCompleteTaskHandlerNotFound(t *testing.T) {
	server := setupAPITestServer(t)

	rr := performRequest(t, server.router, http.MethodPost, "/api/v1/tasks/ash/999/complete", nil)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status mismatch: got %d, expect %d; body=%s", rr.Code, http.StatusNotFound, rr.Body.String())
	}
}

func TestDeleteTaskHandler(t *testing.T) {
	server := setupAPITestServer(t)
	created := createAPITaskWithReward(t, server.repo, task.Task{UserID: "ash", Title: "Catch Pikachu", Tag: "pokemon"}, 25, "pikachu", false)
	createAPIStatistic(t, server.repo, task.UserStatistic{UserID: "ash", TasksOpened: 1})

	rr := performRequest(t, server.router, http.MethodDelete, "/api/v1/tasks/ash/"+strconv.Itoa(created.ID), nil)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("status mismatch: got %d, expect %d; body=%s", rr.Code, http.StatusNoContent, rr.Body.String())
	}

	_, err := server.repo.GetTask(created.ID, "ash")
	if !errors.Is(err, task.ErrTaskNotFound) {
		t.Fatalf("expected %v after delete, got %v", task.ErrTaskNotFound, err)
	}
}

func TestDeleteTaskHandlerNotFound(t *testing.T) {
	server := setupAPITestServer(t)

	rr := performRequest(t, server.router, http.MethodDelete, "/api/v1/tasks/ash/999", nil)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status mismatch: got %d, expect %d; body=%s", rr.Code, http.StatusNotFound, rr.Body.String())
	}
}

func TestDeleteTaskHandlerCompletedTask(t *testing.T) {
	server := setupAPITestServer(t)
	created := createAPITaskWithReward(t, server.repo, task.Task{UserID: "ash", Title: "Catch Pikachu", Tag: "pokemon"}, 25, "pikachu", false)
	if err := server.repo.CompleteTask(created.ID, created.UserID); err != nil {
		t.Fatalf("CompleteTask setup: %v", err)
	}

	rr := performRequest(t, server.router, http.MethodDelete, "/api/v1/tasks/ash/"+strconv.Itoa(created.ID), nil)

	if rr.Code != http.StatusConflict {
		t.Fatalf("status mismatch: got %d, expect %d; body=%s", rr.Code, http.StatusConflict, rr.Body.String())
	}
}

func TestListCollectionHandler(t *testing.T) {
	server := setupAPITestServer(t)
	createAPICollectionEntry(t, server.repo, task.CollectionEntry{UserID: "ash", PokemonID: 25, PokemonName: "pikachu", Rarity: 1, Shiny: false})

	rr := performRequest(t, server.router, http.MethodGet, "/api/v1/users/ash/collection", nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("status mismatch: got %d, expect %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	response := parseJSONResponse[[]CollectionEntryResponse](t, rr)
	if len(response) != 1 {
		t.Fatalf("collection length mismatch: got %d, expect 1", len(response))
	}

	if response[0].PokemonName != "pikachu" {
		t.Fatalf("PokemonName mismatch: got %q, expect pikachu", response[0].PokemonName)
	}
}

func TestGetStatisticHandler(t *testing.T) {
	server := setupAPITestServer(t)
	createAPIStatistic(t, server.repo, task.UserStatistic{UserID: "ash", TasksOpened: 1})

	rr := performRequest(t, server.router, http.MethodGet, "/api/v1/users/ash/stats", nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("status mismatch: got %d, expect %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	response := parseJSONResponse[UserStatisticResponse](t, rr)
	if response.UserID != "ash" {
		t.Fatalf("UserID mismatch: got %q, expect ash", response.UserID)
	}
}

func TestGetStatisticHandlerNotFound(t *testing.T) {
	server := setupAPITestServer(t)

	rr := performRequest(t, server.router, http.MethodGet, "/api/v1/users/missing/stats", nil)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status mismatch: got %d, expect %d; body=%s", rr.Code, http.StatusNotFound, rr.Body.String())
	}
}

func stringPtr(s string) *string {
	return &s
}
