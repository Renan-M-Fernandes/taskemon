package task

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
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

//Handles the tasks

func (s Service) CreateTask(t Task) error {
	if t.Title == "" {
		return ErrEmptyTitle
	}
	if t.Tag == "" {
		t.Tag = "Misc"
	}

	t, err := s.repo.CreateTask(Task{
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		DueAt:       t.DueAt,
		Tag:         t.Tag,
	})
	if err != nil {
		return err
	}

	err = s.CreateTaskReward(t)
	if err != nil {
		return err
	}

	us, err := s.repo.ExistStatistic(t.UserID)
	if err != nil {
		return err
	}

	if us.UserID == "" {
		us.UserID = t.UserID
		us.TasksOpened++
		err := s.repo.CreateUserStatistic(us)
		if err != nil {
			return err
		}
	} else {
		err = s.updateUserStatisticOnCreate(t.UserID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Service) DeleteTask(t Task) error {
	t, err := s.repo.ExistTask(t)

	if err != nil {
		return err
	}

	if t.Completed {
		return ErrTaskCompletedDelete
	}

	err = s.repo.DeleteTaskReward(t.Reward.ID)

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

func (s Service) CompleteTask(t Task) error {
	t, err := s.repo.GetTask(t)
	if err != nil {
		return err
	}

	if t.Completed {
		return ErrTaskAlreadyCompleted
	}

	err = s.repo.CompleteTask(Task{
		ID: t.ID,
	})

	if err != nil {
		return err
	}

	err = s.RevealPokemon(t)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) UpdateTask(dueAtSent bool, update TaskUpdate) error {
	t := Task{
		ID:     update.ID,
		UserID: update.UserID,
	}

	t, err := s.GetTask(t)
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
	}

	err = s.repo.UpdateTask(t)

	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetTask(t Task) (Task, error) {
	return s.repo.GetTask(t)
}

func (s Service) ListTasksByUser(t Task) ([]Task, error) {
	return s.repo.ListTasksByUser(t)
}

//Handles the rewards

func (s Service) CreateTaskReward(t Task) error {
	tr, err := getPokemonFromAPI()
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

func (s Service) RevealPokemon(t Task) error {
	tr := t.Reward

	err := s.repo.RevealPokemon(TaskReward{
		ID: tr.ID,
	})

	if err != nil {
		return err
	}

	err = s.handlesCollectionEntry(t.UserID, tr)

	if err != nil {
		return err
	}

	err = s.updateUserStatisticOnReveal(t.UserID)

	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetTaskReward(t Task) (TaskReward, error) {
	return s.repo.GetTaskReward(Task{
		ID: t.ID,
	})
}

//Handles Collection Entry/Pokedex

func (s Service) handlesCollectionEntry(userID string, tr TaskReward) error {
	ID, err := s.existCollectionEntry(userID, tr.PokemonID)

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

func (s Service) existCollectionEntry(userID string, pokemonID int) (int, error) {
	return s.repo.ExistCollectionEntry(userID, pokemonID)
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
		return fmt.Errorf("update user statistic: date len: %w", err)
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
	h.Version = "1.0.0"
	h.Database = "connected"
	return h, nil
}

func getPokemonFromAPI() (TaskReward, error) {
	id := rand.IntN(maxPokemonSpecies-1) + 1
	url := "https://pokeapi.co/api/v2/pokemon/" + strconv.Itoa(id)
	urlSpecies := "https://pokeapi.co/api/v2/pokemon-species/" + strconv.Itoa(id)
	//Get data from pokemon api
	resp, err := http.Get(url)
	if err != nil {
		return TaskReward{}, fmt.Errorf("pokemon api: get: %w", err)
	}
	var pokemon Pokemon

	err = json.NewDecoder(resp.Body).Decode(&pokemon)
	if err != nil {
		return TaskReward{}, fmt.Errorf("pokemon api: decode: %w", err)
	}

	//Get data from spicies api
	resp, err = http.Get(urlSpecies)
	if err != nil {
		return TaskReward{}, fmt.Errorf("pokemon api species: get: %w", err)
	}
	var rarity Rarity

	err = json.NewDecoder(resp.Body).Decode(&rarity)
	if err != nil {
		return TaskReward{}, fmt.Errorf("pokemon api species: decode: %w", err)
	}
	var tr TaskReward

	tr.PokemonID = id
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

func LoadPokemonSpeciesCount() error {
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon-species?limit=1")
	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic(err)
	}

	maxPokemonSpecies = result.Count
	return err
}
