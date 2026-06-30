package api

import (
	"errors"
	"net/http"
)

var ErrBadRequest = errors.New("bad request")
var ErrTaskNotFound = errors.New("task not found")
var ErrStatisticNotFound = errors.New("user statistic not found")
var ErrTaskRewardNotFound = errors.New("task reward not found")
var ErrTaskAlreadyCompleted = errors.New("task already completed")
var ErrTaskCompletedDelete = errors.New("cannot delete completed task")
var ErrEmptyTitle = errors.New("title cannot be empty")
var ErrEmptyTag = errors.New("tag cannot be empty")

func ErrorCode(err error) int {
	switch {
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest

	case errors.Is(err, ErrEmptyTitle):
		return http.StatusBadRequest

	case errors.Is(err, ErrEmptyTag):
		return http.StatusBadRequest

	case errors.Is(err, ErrTaskNotFound):
		return http.StatusNotFound

	case errors.Is(err, ErrStatisticNotFound):
		return http.StatusNotFound

	case errors.Is(err, ErrTaskRewardNotFound):
		return http.StatusNotFound

	case errors.Is(err, ErrTaskAlreadyCompleted):
		return http.StatusConflict

	case errors.Is(err, ErrTaskCompletedDelete):
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}
