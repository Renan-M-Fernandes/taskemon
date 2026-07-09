package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Renan-M-Fernandes/taskemon/internal/config"
	"github.com/Renan-M-Fernandes/taskemon/internal/database"
	"github.com/Renan-M-Fernandes/taskemon/internal/printer"
	"github.com/Renan-M-Fernandes/taskemon/internal/task"
	"github.com/Renan-M-Fernandes/taskemon/internal/taskprint"
	_ "modernc.org/sqlite"
)

type apiTestServer struct {
	db      *sql.DB
	repo    *task.Repository
	service *task.Service
	printer *recordingPrinter
	handler *Handler
	router  *http.ServeMux
}

type recordingPrinter struct {
	tickets []printer.Ticket
	err     error
}

func (p *recordingPrinter) PrintTicket(ctx context.Context, ticket printer.Ticket) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if p.err != nil {
		return p.err
	}
	p.tickets = append(p.tickets, ticket)
	return nil
}

func loadConfig(t *testing.T) config.Config {
	t.Helper()
	return config.Default()
}

func setupAPITestServer(t *testing.T) apiTestServer {
	t.Helper()
	cfg := loadConfig(t)

	dsn := "file:" + strings.ReplaceAll(t.Name(), "/", "_") + "?mode=memory&cache=shared"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("Open test database: %v", err)
	}
	db.SetMaxOpenConns(1)

	if err := database.Migrate(db); err != nil {
		t.Fatalf("Migrate test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	repo := task.NewRepository(db)
	service := task.NewService(repo)
	testPrinter := &recordingPrinter{}
	taskPrintService := taskprint.NewService(service, testPrinter, taskprint.Config{
		QRMode:  cfg.Printer.QRMode,
		BaseURL: cfg.Printer.BaseURL,
	})
	handler := NewHandler(service, taskPrintService)
	router := NewRouter(handler)

	return apiTestServer{
		db:      db,
		repo:    repo,
		service: service,
		printer: testPrinter,
		handler: handler,
		router:  router,
	}
}

func performRequest(t *testing.T, router http.Handler, method string, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var requestBody *bytes.Reader
	if body == nil {
		requestBody = bytes.NewReader(nil)
	} else {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Marshal request body: %v", err)
		}
		requestBody = bytes.NewReader(data)
	}

	req := httptest.NewRequest(method, path, requestBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func performInvalidRequest(router http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func createAPITask(t *testing.T, repo *task.Repository, input task.Task) task.Task {
	t.Helper()

	created, err := repo.CreateTask(input)
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	return created
}

func createAPIReward(t *testing.T, repo *task.Repository, taskID int, pokemonID int, pokemonName string, shiny bool) task.TaskReward {
	t.Helper()

	err := repo.CreateTaskReward(task.TaskReward{
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

func createAPITaskWithReward(t *testing.T, repo *task.Repository, input task.Task, pokemonID int, pokemonName string, shiny bool) task.Task {
	t.Helper()

	created := createAPITask(t, repo, input)
	created.Reward = createAPIReward(t, repo, created.ID, pokemonID, pokemonName, shiny)

	return created
}

func createAPIStatistic(t *testing.T, repo *task.Repository, input task.UserStatistic) task.UserStatistic {
	t.Helper()

	if err := repo.CreateUserStatistic(input); err != nil {
		t.Fatalf("CreateUserStatistic: %v", err)
	}

	stats, err := repo.GetStatistic(input.UserID)
	if err != nil {
		t.Fatalf("GetStatistic: %v", err)
	}

	return stats
}

func createAPICollectionEntry(t *testing.T, repo *task.Repository, input task.CollectionEntry) {
	t.Helper()

	if err := repo.CreateCollectionEntry(input); err != nil {
		t.Fatalf("CreateCollectionEntry: %v", err)
	}
}

func parseJSONResponse[T any](t *testing.T, rr *httptest.ResponseRecorder) T {
	t.Helper()

	var response T
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Unmarshal response: %v\nbody: %s", err, rr.Body.String())
	}

	return response
}
