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

func ToTaskResponseSlice(task []task.Task) []TaskResponse {
	var response []TaskResponse

	for _, t := range task {
		res := TaskResponse{
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
			res.Reward.PokemonID = &t.Reward.PokemonID
			res.Reward.PokemonName = &t.Reward.PokemonName
			res.Reward.PokemonSprite = &t.Reward.Sprite
			res.Reward.RevealedAt = t.Reward.RevealedAt
		}

		response = append(response, res)
	}
	return response
}
