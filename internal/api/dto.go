package api

import "time"

type TaskRequest struct {
	ID          int            `json:"id"`
	UserID      string         `json:"userID"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Completed   bool           `json:"completed"`
	DueAt       *time.Time     `json:"dueAt"`
	Tag         string         `json:"tag"`
	CreatedAt   time.Time      `json:"createdAt"`
	CompletedAt *time.Time     `json:"completedAt"`
	Reward      RewardResponse `json:"reward"`
}

type TaskRequestUpdate struct {
	ID          int        `json:"id"`
	UserID      string     `json:"userID"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	DueAt       *time.Time `json:"dueAt"`
	Tag         *string    `json:"tag"`
}

type TaskResponse struct {
	ID          int            `json:"id"`
	UserID      string         `json:"userID"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Completed   bool           `json:"completed"`
	DueAt       *time.Time     `json:"dueAt"`
	Tag         string         `json:"tag"`
	CreatedAt   time.Time      `json:"createdAt"`
	CompletedAt *time.Time     `json:"completedAt"`
	Reward      RewardResponse `json:"reward"`
}

type RewardResponse struct {
	Revealed      bool       `json:"revealed"`
	Rarity        int        `json:"rarity"`
	Shiny         bool       `json:"shiny"`
	PokemonID     *int       `json:"pokemonID,omitempty"`
	PokemonName   *string    `json:"pokemonName,omitempty"`
	PokemonSprite *string    `json:"pokemonSprite,omitempty"`
	RevealedAt    *time.Time `json:"revealedAt,omitempty"`
}
