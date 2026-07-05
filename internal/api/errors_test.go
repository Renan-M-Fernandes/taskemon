package api

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Renan-M-Fernandes/taskemon/internal/task"
)

func TestErrorCodeMapsDomainErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "bad request", err: task.ErrBadRequest, want: http.StatusBadRequest},
		{name: "empty title", err: task.ErrEmptyTitle, want: http.StatusBadRequest},
		{name: "empty tag", err: task.ErrEmptyTag, want: http.StatusBadRequest},
		{name: "pokemon unavailable", err: task.ErrPokemonSpeciesUnavailable, want: http.StatusServiceUnavailable},
		{name: "task not found", err: task.ErrTaskNotFound, want: http.StatusNotFound},
		{name: "statistic not found", err: task.ErrStatisticNotFound, want: http.StatusNotFound},
		{name: "reward not found", err: task.ErrTaskRewardNotFound, want: http.StatusNotFound},
		{name: "already completed", err: task.ErrTaskAlreadyCompleted, want: http.StatusConflict},
		{name: "completed task delete", err: task.ErrTaskCompletedDelete, want: http.StatusConflict},
		{name: "unknown", err: errors.New("boom"), want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ErrorCode(tt.err)
			if got != tt.want {
				t.Fatalf("status mismatch: got %d, expect %d", got, tt.want)
			}
		})
	}
}
