package api

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/Renan-M-Fernandes/taskemon/internal/task"
)

func TestToTaskResponseHidesUnrevealedReward(t *testing.T) {
	response := ToTaskResponse(task.Task{
		ID:    1,
		Title: "Catch Pikachu",
		Reward: task.TaskReward{
			PokemonID:   25,
			PokemonName: "pikachu",
			Sprite:      "pikachu.png",
			Rarity:      1,
			Shiny:       true,
			Revealed:    false,
		},
	})

	body, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Marshal response: %v", err)
	}

	jsonBody := string(body)
	for _, hidden := range []string{"pokemonID", "pokemonName", "pokemonSprite", "pikachu"} {
		if strings.Contains(jsonBody, hidden) {
			t.Fatalf("hidden reward leaked %q in JSON: %s", hidden, jsonBody)
		}
	}

	if !strings.Contains(jsonBody, `"revealed":false`) {
		t.Fatalf("revealed flag missing from JSON: %s", jsonBody)
	}
}

func TestToTaskResponseShowsRevealedReward(t *testing.T) {
	revealedAt := time.Now().Truncate(time.Second)
	response := ToTaskResponse(task.Task{
		ID:    1,
		Title: "Catch Pikachu",
		Reward: task.TaskReward{
			PokemonID:   25,
			PokemonName: "pikachu",
			Sprite:      "pikachu.png",
			Rarity:      1,
			Shiny:       false,
			Revealed:    true,
			RevealedAt:  &revealedAt,
		},
	})

	if response.Reward.PokemonID == nil || *response.Reward.PokemonID != 25 {
		t.Fatalf("PokemonID mismatch: got %v, expect 25", response.Reward.PokemonID)
	}

	if response.Reward.PokemonName == nil || *response.Reward.PokemonName != "pikachu" {
		t.Fatalf("PokemonName mismatch: got %v, expect pikachu", response.Reward.PokemonName)
	}

	if response.Reward.PokemonSprite == nil || *response.Reward.PokemonSprite != "pikachu.png" {
		t.Fatalf("PokemonSprite mismatch: got %v, expect pikachu.png", response.Reward.PokemonSprite)
	}

	if response.Reward.RevealedAt == nil {
		t.Fatal("expected RevealedAt to be exposed")
	}
}

func TestToTaskResponseSliceEmptyReturnsEmptyArray(t *testing.T) {
	response := ToTaskResponseSlice(nil)

	body, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Marshal response: %v", err)
	}

	if string(body) != "[]" {
		t.Fatalf("empty response mismatch: got %s, expect []", body)
	}
}

func TestToCollectionResponseSliceEmptyReturnsEmptyArray(t *testing.T) {
	response := ToCollectionResponseSlice(nil)

	body, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Marshal response: %v", err)
	}

	if string(body) != "[]" {
		t.Fatalf("empty response mismatch: got %s, expect []", body)
	}
}

func TestToUserStatisticResponse(t *testing.T) {
	stats := task.UserStatistic{
		UserID:         "ash",
		TasksCompleted: 2,
		TasksOpened:    3,
		TasksDeleted:   1,
		PokemonCaught:  2,
		ShinyCaught:    1,
		UniquePokemon:  2,
		CurrentStreak:  4,
		LongestStreak:  7,
	}

	response := ToUserStatisticResponse(stats)

	if response.UserID != stats.UserID {
		t.Fatalf("UserID mismatch: got %q, expect %q", response.UserID, stats.UserID)
	}
	if response.TasksCompleted != stats.TasksCompleted {
		t.Fatalf("TasksCompleted mismatch: got %d, expect %d", response.TasksCompleted, stats.TasksCompleted)
	}
	if response.LongestStreak != stats.LongestStreak {
		t.Fatalf("LongestStreak mismatch: got %d, expect %d", response.LongestStreak, stats.LongestStreak)
	}
}
