package main

import (
	// "flag"

	"log"
	"net/http"

	"github.com/Renan-M-Fernandes/taskemon/internal/api"
	"github.com/Renan-M-Fernandes/taskemon/internal/database"
	"github.com/Renan-M-Fernandes/taskemon/internal/task"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = database.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	task.LoadPokemonSpeciesCount()

	repo := task.NewRepository(db)
	service := task.NewService(repo)
	handler := api.NewHandler(service)

	api.RegisterRoutes(handler)

	log.Fatal(
		http.ListenAndServe(
			":8080",
			api.EnableCORS(http.DefaultServeMux),
		),
	)
}
