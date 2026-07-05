package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Renan-M-Fernandes/taskemon/internal/api"
	"github.com/Renan-M-Fernandes/taskemon/internal/database"
	"github.com/Renan-M-Fernandes/taskemon/internal/task"
)

func main() {
	db, err := database.Connect("taskemon.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = database.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = task.LoadPokemonSpeciesCount(ctx)
	if err != nil {
		log.Fatal(err)
	}

	repo := task.NewRepository(db)
	service := task.NewService(repo)
	handler := api.NewHandler(service)

	mux := api.NewRouter(handler)

	log.Fatal(
		http.ListenAndServe(
			":8080",
			api.EnableCORS(mux),
		),
	)
}
