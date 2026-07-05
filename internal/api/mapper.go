package api

import (
	"github.com/Renan-M-Fernandes/taskemon/internal/task"
)

func ToTaskResponse(t task.Task) TaskResponse {
	response := TaskResponse{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		DueAt:       t.DueAt,
		Tag:         t.Tag,
		CreatedAt:   t.CreatedAt,
		CompletedAt: t.CompletedAt,
		Reward: RewardResponse{
			Revealed: t.Reward.Revealed,
			Rarity:   t.Reward.Rarity,
			Shiny:    t.Reward.Shiny,
		},
	}

	if t.Reward.Revealed {
		response.Reward.PokemonID = &t.Reward.PokemonID
		response.Reward.PokemonName = &t.Reward.PokemonName
		response.Reward.PokemonSprite = &t.Reward.Sprite
		response.Reward.RevealedAt = t.Reward.RevealedAt
	}

	return response
}

func ToTaskResponseSlice(tasks []task.Task) []TaskResponse {
	response := make([]TaskResponse, 0, len(tasks))

	for _, t := range tasks {
		response = append(response, ToTaskResponse(t))
	}

	return response
}

func ToCollectionResponse(c task.CollectionEntry) CollectionEntryResponse {
	return CollectionEntryResponse{
		PokemonID:     c.PokemonID,
		PokemonName:   c.PokemonName,
		Count:         c.Count,
		Rarity:        c.Rarity,
		Shiny:         c.Shiny,
		FirstCaughtAt: c.FirstCaughtAt,
		LastCaughtAt:  c.LastCaughtAt,
	}
}

func ToCollectionResponseSlice(collection []task.CollectionEntry) []CollectionEntryResponse {
	response := make([]CollectionEntryResponse, 0, len(collection))

	for _, c := range collection {
		response = append(response, ToCollectionResponse(c))
	}

	return response
}

func ToUserStatisticResponse(s task.UserStatistic) UserStatisticResponse {
	return UserStatisticResponse{
		UserID:         s.UserID,
		TasksCompleted: s.TasksCompleted,
		TasksOpened:    s.TasksOpened,
		TasksDeleted:   s.TasksDeleted,
		PokemonCaught:  s.PokemonCaught,
		ShinyCaught:    s.ShinyCaught,
		UniquePokemon:  s.UniquePokemon,
		CurrentStreak:  s.CurrentStreak,
		LongestStreak:  s.LongestStreak,
	}
}
