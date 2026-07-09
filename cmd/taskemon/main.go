package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Renan-M-Fernandes/taskemon/internal/api"
	"github.com/Renan-M-Fernandes/taskemon/internal/config"
	"github.com/Renan-M-Fernandes/taskemon/internal/database"
	"github.com/Renan-M-Fernandes/taskemon/internal/printer"
	"github.com/Renan-M-Fernandes/taskemon/internal/task"
	"github.com/Renan-M-Fernandes/taskemon/internal/taskprint"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg.Database.Path)
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

	taskPrinter, err := printer.NewFromConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	taskPrintService := taskprint.NewService(service, taskPrinter, taskprint.Config{
		QRMode:  cfg.Printer.QRMode,
		BaseURL: cfg.Printer.BaseURL,
	})

	handler := api.NewHandler(service, taskPrintService)

	mux := api.NewRouter(handler)

	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	log.Fatal(
		http.ListenAndServe(
			serverAddr,
			api.EnableCORS(mux, cfg.Server.CORSOrigins),
		),
	)
}
