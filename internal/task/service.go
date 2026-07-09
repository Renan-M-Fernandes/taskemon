package task

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var maxPokemonSpecies int

var result struct {
	Count int `json:"count"`
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

var httpClient = &http.Client{
	Timeout: 20 * time.Second,
}

//Handles the tasks

func (s Service) CreateTask(t Task) (Task, error) {
	if t.Title == "" {
		return Task{}, ErrEmptyTitle
	}
	if t.Tag == "" {
		t.Tag = "misc"
	}
	t.Tag = strings.ToLower(t.Tag)

	t, err := s.repo.CreateTask(t)
	if err != nil {
		return Task{}, err
	}

	err = s.CreateTaskReward(t)
	if err != nil {
		return Task{}, err
	}

	us, err := s.repo.ExistStatistic(t.UserID)
	if err != nil {
		return Task{}, err
	}

	if us.UserID == "" {
		us.UserID = t.UserID
		us.TasksOpened++
		err := s.repo.CreateUserStatistic(us)
		if err != nil {
			return Task{}, err
		}
	} else {
		err = s.updateUserStatisticOnCreate(t.UserID)
		if err != nil {
			return Task{}, err
		}
	}
	t, err = s.repo.GetTask(t.ID, t.UserID)
	if err != nil {
		return Task{}, err
	}

	return t, nil
}

func (s Service) DeleteTask(ID int, userName string) error {
	t, err := s.repo.ExistTask(ID, userName)

	if err != nil {
		return err
	}

	if t.Completed {
		return ErrTaskCompletedDelete
	}

	err = s.repo.DeleteTaskReward(t.ID)

	if err != nil {
		return err
	}

	err = s.repo.DeleteTask(t.ID)

	if err != nil {
		return err
	}

	err = s.updateUserStatisticOnDelete(t.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) CompleteTask(ID int, userID string) (Task, error) {
	t, err := s.repo.GetTask(ID, userID)
	if err != nil {
		return Task{}, err
	}

	if t.Completed {
		return Task{}, ErrTaskAlreadyCompleted
	}

	err = s.repo.CompleteTask(ID, userID)

	if err != nil {
		return Task{}, err
	}

	_, err = s.revealPokemon(t)
	if err != nil {
		return Task{}, err
	}

	t, err = s.repo.GetTask(ID, userID)
	if err != nil {
		return Task{}, err
	}

	return t, nil
}

func (s Service) UpdateTask(dueAtSent bool, update TaskUpdate) error {
	t, err := s.GetTask(update.ID, update.UserID)
	if err != nil {
		return err
	}

	if t.Completed {
		return ErrTaskAlreadyCompleted
	}

	if update.Title != nil {
		if *update.Title == "" {
			return ErrEmptyTitle
		}
		t.Title = *update.Title
	}

	if update.Description != nil {
		t.Description = *update.Description
	}

	if update.DueAt != nil {
		t.DueAt = update.DueAt
	} else if dueAtSent {
		t.DueAt = nil
	}

	if update.Tag != nil {
		if *update.Tag == "" {
			return ErrEmptyTag
		}
		t.Tag = *update.Tag
		t.Tag = strings.ToLower(t.Tag)
	}

	err = s.repo.UpdateTask(t)

	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetTask(ID int, userID string) (Task, error) {
	return s.repo.GetTask(ID, userID)
}

func (s Service) ListTasksByUser(userID string) ([]Task, error) {
	return s.repo.ListTasksByUser(userID)
}

func (s Service) ListTasksByUserNotCompleted(userID string) ([]Task, error) {
	return s.repo.ListTasksByUserNotCompleted(userID)
}

func (s Service) ListTasksByUserCompleted(userID string) ([]Task, error) {
	return s.repo.ListTasksByUserCompleted(userID)
}

//Handles the rewards

func (s Service) CreateTaskReward(t Task) error {
	tr, err := getPokemonFromAPIs()
	if err != nil {
		return err
	}
	err = s.repo.CreateTaskReward(TaskReward{
		TaskID:      t.ID,
		PokemonID:   tr.PokemonID,
		PokemonName: tr.PokemonName,
		Sprite:      tr.Sprite,
		Rarity:      tr.Rarity,
		Shiny:       tr.Shiny,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s Service) revealPokemon(t Task) (TaskReward, error) {
	err := s.repo.RevealPokemon(t.ID)

	if err != nil {
		return TaskReward{}, err
	}

	err = s.handlesCollectionEntry(t.UserID, t.Reward)

	if err != nil {
		return TaskReward{}, err
	}

	err = s.updateUserStatisticOnReveal(t.UserID)

	if err != nil {
		return TaskReward{}, err
	}

	return TaskReward{}, nil
}

func (s Service) GetTaskReward(ID int, userID string) (TaskReward, error) {
	return s.repo.GetTaskReward(ID)
}

//Handles Collection Entry/Pokedex

func (s Service) handlesCollectionEntry(userID string, tr TaskReward) error {
	ID, err := s.existCollectionEntry(userID, tr.PokemonID, tr.Shiny)

	if err != nil {
		return err
	}

	if ID == 0 {
		err := s.createCollectionEntry(userID, tr)
		if err != nil {
			return err
		}
		return nil
	} else {
		err := s.updateCollectionEntry(tr.PokemonID, tr.Shiny, userID)
		if err != nil {
			return err
		}
		return nil
	}
}

func (s Service) createCollectionEntry(userID string, tr TaskReward) error {
	err := s.repo.CreateCollectionEntry(CollectionEntry{
		UserID:      userID,
		PokemonID:   tr.PokemonID,
		PokemonName: tr.PokemonName,
		Rarity:      tr.Rarity,
		Shiny:       tr.Shiny,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s Service) updateCollectionEntry(pokemonID int, shiny bool, userID string) error {
	err := s.repo.UpdateCollectionEntry(pokemonID, shiny, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) ListCollection(userID string) ([]CollectionEntry, error) {
	return s.repo.ListCollection(userID)
}

func (s Service) existCollectionEntry(userID string, pokemonID int, shiny bool) (int, error) {
	return s.repo.ExistCollectionEntry(userID, pokemonID, shiny)
}

// Handles Statistics
func (s Service) getDataForStatistic(userID string) (int, int, int, []time.Time, error) {
	return s.repo.GetDataForStatistic(userID)
}

func (s Service) GetStatistic(userID string) (UserStatistic, error) {
	return s.repo.GetStatistic(userID)
}

func (s Service) CreateUserStatistic(us UserStatistic) error {
	return s.repo.CreateUserStatistic(us)
}

func (s Service) updateUserStatisticOnReveal(userID string) error {
	longestStreak, shinyCaughtTotal, uniquePokemonTotal, dates, err := s.getDataForStatistic(userID)
	if err != nil {
		return err
	}
	today := time.Now().Truncate(24 * time.Hour)

	// If no task today, start checking yesterday.
	expected := today

	if len(dates) == 0 {
		return fmt.Errorf("update user statistic: no completion dates for user %q", userID)
	}

	// If today isn't completed, the streak starts by expecting yesterday.
	if !dates[0].Equal(today) {
		expected = today.AddDate(0, 0, -1)
	}

	streak := 0

	for _, d := range dates {
		if d.Equal(expected) {
			streak++
			expected = expected.AddDate(0, 0, -1)
		} else {
			break
		}
	}
	currentStreak := streak
	if longestStreak < currentStreak {
		longestStreak = currentStreak
	}

	err = s.repo.UpdateUserStatisticOnClose(userID, shinyCaughtTotal, uniquePokemonTotal, currentStreak, longestStreak)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) updateUserStatisticOnCreate(userID string) error {
	return s.repo.UpdateUserStatisticOnCreate(userID)
}

func (s Service) updateUserStatisticOnDelete(userID string) error {
	return s.repo.UpdateUserStatisticOnDelete(userID)
}

//General use Functions

func (s Service) GetHealth() (Health, error) {
	var h Health
	h.Status = "OK"
	h.Version = "0.2.0"
	h.Database = "connected"
	return h, nil
}

func getPokemonFromAPIs() (TaskReward, error) {
	if maxPokemonSpecies < 1 {
		return TaskReward{}, ErrPokemonSpeciesUnavailable
	}
	ID := rand.IntN(maxPokemonSpecies) + 1
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tr TaskReward
	pokemon, err := getPokemonFromPokemonAPI(ctx, ID)
	if err != nil {
		return TaskReward{}, err
	}

	rarity, err := getPokemonFromSpicieAPI(ctx, ID)
	if err != nil {
		return TaskReward{}, err
	}

	tr.PokemonID = ID
	tr.PokemonName = pokemon.Name
	tr.Rarity = checkRarity(rarity)
	tr.Shiny = IsShiny()
	if tr.Shiny {
		tr.Sprite = pokemon.Sprites.FrontShiny
	} else {
		tr.Sprite = pokemon.Sprites.FrontDefault
	}

	return tr, nil
}

func getPokemonFromPokemonAPI(ctx context.Context, ID int) (Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + strconv.Itoa(ID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Pokemon{}, fmt.Errorf("pokemon api pokemon: create request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return Pokemon{}, fmt.Errorf("pokemon api pokemon: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Pokemon{}, fmt.Errorf("pokemon api pokemon: bad status: %s", resp.Status)
	}

	var pokemon Pokemon
	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		return Pokemon{}, fmt.Errorf("pokemon api pokemon: decode: %w", err)
	}

	return pokemon, nil
}

func getPokemonFromSpicieAPI(ctx context.Context, ID int) (Rarity, error) {
	url := "https://pokeapi.co/api/v2/pokemon-species/" + strconv.Itoa(ID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Rarity{}, fmt.Errorf("pokemon api species: create request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return Rarity{}, fmt.Errorf("pokemon api species: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Rarity{}, fmt.Errorf("pokemon api species: bad status: %s", resp.Status)
	}

	var rarity Rarity
	if err := json.NewDecoder(resp.Body).Decode(&rarity); err != nil {
		return Rarity{}, fmt.Errorf("pokemon api species: decode: %w", err)
	}

	return rarity, nil
}

func IsShiny() bool {
	return rand.IntN(50) == 0
}

func checkRarity(rarity Rarity) int {
	if rarity.IsMythical {
		return 5
	} else if rarity.IsLegendary {
		return 4
	} else {
		return 1
	}
}

func LoadPokemonSpeciesCount(ctx context.Context) error {
	url := "https://pokeapi.co/api/v2/pokemon-species?limit=1"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("load species count: get: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("load species count: do request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("load species count: bad status: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("load species count: decode: %w", err)
	}

	maxPokemonSpecies = result.Count

	if maxPokemonSpecies < 1 {
		return ErrPokemonSpeciesUnavailable
	}

	return nil
}
