package api

import "net/http"

func NewRouter(handler *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	RegisterRoutes(mux, handler)
	return mux
}

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc(
		"GET    /api/v1/health",
		handler.GetHealth,
	)

	mux.HandleFunc(
		"GET    /api/v1/tasks/{userID}/{ID}",
		handler.GetTasks,
	)

	mux.HandleFunc(
		"GET    /api/v1/tasks/{userID}",
		handler.ListTasks,
	)

	mux.HandleFunc(
		"GET    /api/v1/tasks/completed/{userID}",
		handler.ListTasksCompleted,
	)

	mux.HandleFunc(
		"GET    /api/v1/tasks/open/{userID}",
		handler.ListTasksNotCompleted,
	)

	mux.HandleFunc(
		"POST   /api/v1/tasks/{userID}",
		handler.CreateTask,
	)

	mux.HandleFunc(
		"DELETE /api/v1/tasks/{userID}/{ID}",
		handler.DeleteTask,
	)

	mux.HandleFunc(
		"PATCH  /api/v1/tasks/{userID}/{ID}",
		handler.UpdateTask,
	)

	mux.HandleFunc(
		"POST   /api/v1/tasks/{userID}/{ID}/complete",
		handler.CompleteTask,
	)

	mux.HandleFunc(
		"GET    /api/v1/users/{userID}/collection",
		handler.ListCollection,
	)

	mux.HandleFunc(
		"GET    /api/v1/users/{userID}/stats",
		handler.GetStatistic,
	)

	mux.HandleFunc(
		"POST /api/v1/tasks/{userID}/{ID}/print",
		handler.PrintTask,
	)
}
