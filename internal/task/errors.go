package task

import "errors"

var ErrBadRequest = errors.New("bad request")

var ErrTaskNotFound = errors.New("task not found")

var ErrStatisticNotFound = errors.New("user statistic not found")

var ErrTaskRewardNotFound = errors.New("task reward not found")

var ErrTaskAlreadyCompleted = errors.New("task already completed")

var ErrTaskCompletedDelete = errors.New("cannot delete completed task")

var ErrEmptyTitle = errors.New("title cannot be empty")

var ErrEmptyTag = errors.New("tag cannot be empty")
