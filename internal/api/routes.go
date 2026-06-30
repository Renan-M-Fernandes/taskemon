package api

import "net/http"

func RegisterRoutes(handler *Handler) {
	http.HandleFunc(
		"GET    /api/v1/health",
		handler.GetHealth,
	)

	http.HandleFunc(
		"GET    /api/v1/tasks/{userID}/{ID}",
		handler.GetTasks,
	)

	http.HandleFunc(
		"GET    /api/v1/tasks/{userID}",
		handler.ListTasks,
	)

	http.HandleFunc(
		"POST   /api/v1/tasks/{userID}",
		handler.CreateTask,
	)

	http.HandleFunc(
		"DELETE /api/v1/tasks/{userID}/{ID}",
		handler.DeleteTask,
	)

	http.HandleFunc(
		"PATCH  /api/v1/tasks/{userID}/{ID}",
		handler.UpdateTask,
	)

	http.HandleFunc(
		"POST   /api/v1/tasks/{userID}/{ID}/complete",
		handler.CompleteTask,
	)

	http.HandleFunc(
		"GET    /api/v1/users/{userID}/collection",
		handler.ListCollection,
	)

	http.HandleFunc(
		"GET    /api/v1/users/{userID}/stats",
		handler.GetStatistic,
	)
}
