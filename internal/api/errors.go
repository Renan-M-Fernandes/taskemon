package api

import (
	"errors"
	"net/http"

	"github.com/Renan-M-Fernandes/taskemon/internal/task"
)

func ErrorCode(err error) int {
	switch {
	case errors.Is(err, task.ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, task.ErrEmptyTitle):
		return http.StatusBadRequest
	case errors.Is(err, task.ErrEmptyTag):
		return http.StatusBadRequest
	case errors.Is(err, task.ErrPokemonSpeciesUnavailable):
		return http.StatusServiceUnavailable
	case errors.Is(err, task.ErrTaskNotFound):
		return http.StatusNotFound
	case errors.Is(err, task.ErrStatisticNotFound):
		return http.StatusNotFound
	case errors.Is(err, task.ErrTaskRewardNotFound):
		return http.StatusNotFound
	case errors.Is(err, task.ErrCollectionEntryNotFound):
		return http.StatusNotFound
	case errors.Is(err, task.ErrTaskAlreadyCompleted):
		return http.StatusConflict
	case errors.Is(err, task.ErrTaskCompletedDelete):
		return http.StatusConflict
	case errors.Is(err, task.ErrContextTimeout):
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}
