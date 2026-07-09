package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/Renan-M-Fernandes/taskemon/internal/task"
	"github.com/Renan-M-Fernandes/taskemon/internal/taskprint"
)

type Handler struct {
	taskService *task.Service
	taskPrint   *taskprint.Service
}

func NewHandler(taskService *task.Service, taskPrint *taskprint.Service) *Handler {
	return &Handler{
		taskService: taskService,
		taskPrint:   taskPrint,
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

	task, err := h.taskService.GetTask(taskID, userID)

	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToTaskResponse(task))
}

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")

	task, err := h.taskService.ListTasksByUser(userID)

	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToTaskResponseSlice(task))
}

func (h *Handler) ListTasksCompleted(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")

	task, err := h.taskService.ListTasksByUserCompleted(userID)

	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToTaskResponseSlice(task))
}

func (h *Handler) ListTasksNotCompleted(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")

	task, err := h.taskService.ListTasksByUserNotCompleted(userID)

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

	_, err = h.taskService.CreateTask(task.Task{
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

	err = h.taskService.DeleteTask(taskID, userID)

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

	_, err = h.taskService.CompleteTask(taskID, userID)

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
	json.NewEncoder(w).Encode(ToCollectionResponseSlice(collection))
}

func (h *Handler) GetStatistic(w http.ResponseWriter, r *http.Request) {

	userID := r.PathValue("userID")

	stats, err := h.taskService.GetStatistic(userID)

	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToUserStatisticResponse(stats))
}

func (h *Handler) PrintTask(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	taskID, err := strconv.Atoi(r.PathValue("ID"))
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	err = h.taskPrint.PrintTask(r.Context(), taskID, userID)

	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
