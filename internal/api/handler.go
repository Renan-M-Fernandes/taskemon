package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/Renan-M-Fernandes/taskemon/internal/task"
)

type Handler struct {
	taskService *task.Service
}

func NewHandler(
	taskService *task.Service,
) *Handler {
	return &Handler{
		taskService: taskService,
	}
}

func (h *Handler) GetHealth(
	w http.ResponseWriter,
	r *http.Request,
) {
	health, err := h.taskService.GetHealth()

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {

	userID := r.PathValue("userID")
	taskID, err := strconv.Atoi(r.PathValue("ID"))
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.GetTask(task.Task{
		ID:     taskID,
		UserID: userID,
	})

	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToTaskResponse(task))
}

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")

	task, err := h.taskService.ListTasksByUser(task.Task{
		UserID: userID,
	})

	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToTaskResponseSlice(task))
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	var req TaskRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err = h.taskService.CreateTask(task.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DueAt:       req.DueAt,
		Tag:         req.Tag,
	},
	)

	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	taskID, err := strconv.Atoi(r.PathValue("ID"))
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	err = h.taskService.DeleteTask(task.Task{
		ID:     taskID,
		UserID: userID,
	},
	)

	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	taskID, err := strconv.Atoi(r.PathValue("ID"))
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	err = h.taskService.CompleteTask(task.Task{
		ID:     taskID,
		UserID: userID,
	},
	)

	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	taskID, err := strconv.Atoi(r.PathValue("ID"))
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	var req TaskRequestUpdate
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, dueAtSent := raw["dueAt"]

	err = h.taskService.UpdateTask(dueAtSent, task.TaskUpdate{
		ID:          taskID,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DueAt:       req.DueAt,
		Tag:         req.Tag,
	},
	)

	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ListCollection(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")

	collection, err := h.taskService.ListCollection(userID)

	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collection)
}

func (h *Handler) GetStatistic(w http.ResponseWriter, r *http.Request) {

	userID := r.PathValue("userID")

	stats, err := h.taskService.GetStatistic(userID)

	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
