package taskprint

import (
	"testing"

	"github.com/Renan-M-Fernandes/taskemon/internal/task"
)

func TestBuildTaskQRValuePokemonPlaceholder(t *testing.T) {
	input := task.Task{
		Reward: task.TaskReward{
			PokemonID: 25,
		},
	}

	got := buildTaskQRValue(input, Config{QRMode: QRModePokemonPlaceholder})
	want := "https://www.pokemon.com/us/pokedex/25"

	if got != want {
		t.Fatalf("QR value mismatch: got %q, want %q", got, want)
	}
}

func TestBuildTaskQRValueTaskCompletion(t *testing.T) {
	input := task.Task{
		ID:     7,
		UserID: "ash ketchum",
	}

	got := buildTaskQRValue(input, Config{
		QRMode:  QRModeTaskCompletion,
		BaseURL: "http://localhost:8080/",
	})
	want := "http://localhost:8080/tasks/ash%20ketchum/7/complete"

	if got != want {
		t.Fatalf("QR value mismatch: got %q, want %q", got, want)
	}
}

func TestBuildTaskQRValueTaskCompletionWithoutBaseURL(t *testing.T) {
	input := task.Task{
		ID:     7,
		UserID: "ash ketchum",
	}

	got := buildTaskQRValue(input, Config{QRMode: QRModeTaskCompletion})
	want := "taskemon://tasks/ash%20ketchum/7/complete"

	if got != want {
		t.Fatalf("QR value mismatch: got %q, want %q", got, want)
	}
}

func TestBuildTaskQRValueDefaultsToPokemonPlaceholder(t *testing.T) {
	input := task.Task{
		Reward: task.TaskReward{
			PokemonID: 150,
		},
	}

	got := buildTaskQRValue(input, Config{})
	want := "https://www.pokemon.com/us/pokedex/150"

	if got != want {
		t.Fatalf("QR value mismatch: got %q, want %q", got, want)
	}
}
