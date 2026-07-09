package taskprint

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/Renan-M-Fernandes/taskemon/internal/printer"
	"github.com/Renan-M-Fernandes/taskemon/internal/task"
)

const (
	QRModePokemonPlaceholder = "pokemon_placeholder"
	QRModeTaskCompletion     = "task_completion"
)

type Config struct {
	QRMode  string
	BaseURL string
}

type Service struct {
	tasks   *task.Service
	printer printer.Printer
	config  Config
}

func NewService(tasks *task.Service, printer printer.Printer, cfg Config) *Service {
	return &Service{
		tasks:   tasks,
		printer: printer,
		config:  normalizeConfig(cfg),
	}
}

func (s Service) PrintTask(ctx context.Context, taskID int, userID string) error {
	t, err := s.tasks.GetTask(taskID, userID)
	if err != nil {
		return fmt.Errorf("print task ticket: get task for printing: %w", err)
	}

	ticket := printer.Ticket{
		TaskID:      t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Tag:         t.Tag,
		DueAt:       t.DueAt,
		CreatedAt:   t.CreatedAt,
		Shiny:       t.Reward.Shiny,
		QRValue:     buildTaskQRValue(t, s.config),
	}

	err = s.printer.PrintTicket(ctx, ticket)
	if err != nil {
		return fmt.Errorf("print task: print ticket %w", err)
	}

	return nil
}

func normalizeConfig(cfg Config) Config {
	cfg.QRMode = strings.ToLower(strings.TrimSpace(cfg.QRMode))
	cfg.BaseURL = strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")

	if cfg.QRMode == "" {
		cfg.QRMode = QRModePokemonPlaceholder
	}

	return cfg
}

func buildTaskQRValue(t task.Task, cfg Config) string {
	cfg = normalizeConfig(cfg)

	switch cfg.QRMode {
	case QRModeTaskCompletion:
		if cfg.BaseURL == "" {
			return fmt.Sprintf("taskemon://tasks/%s/%d/complete", url.PathEscape(t.UserID), t.ID)
		}

		return fmt.Sprintf("%s/tasks/%s/%d/complete", cfg.BaseURL, url.PathEscape(t.UserID), t.ID)

	case QRModePokemonPlaceholder:
		fallthrough
	default:
		return fmt.Sprintf("https://www.pokemon.com/us/pokedex/%d", t.Reward.PokemonID)
	}
}
