package task

import "time"

type Task struct {
	ID          int        `json:"id"`
	UserID      string     `json:"userID"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	DueAt       *time.Time `json:"dueAt"`
	Tag         string     `json:"tag"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt"`
	Reward      TaskReward
}

type TaskReward struct {
	ID          int
	TaskID      int
	PokemonID   int
	PokemonName string
	Sprite      string
	Rarity      int
	Shiny       bool
	Revealed    bool
	GeneratedAt time.Time
	RevealedAt  *time.Time
}

type CollectionEntry struct {
	ID            int
	UserID        string
	PokemonID     int
	PokemonName   string
	Count         int
	Rarity        int
	Shiny         bool
	FirstCaughtAt time.Time
	LastCaughtAt  time.Time
}

type UserStatistic struct {
	UserID         string
	TasksCompleted int
	PokemonCaught  int
	TasksOpened    int
	TasksDeleted   int
	ShinyCaught    int
	UniquePokemon  int
	CurrentStreak  int
	LongestStreak  int
}

type Health struct {
	Status   string `json:"status"`
	Version  string `json:"version"`
	Database string `json:"database"`
}

type TaskUpdate struct {
	ID          int        `json:"id"`
	UserID      string     `json:"userID"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description"`
	DueAt       *time.Time `json:"dueAt,omitempty"`
	Tag         *string    `json:"tag,omitempty"`
}

type Pokemon struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Sprites Sprites `json:"sprites"`
}

type Sprites struct {
	FrontDefault string `json:"front_default"`
	FrontShiny   string `json:"front_shiny"`
}

type Rarity struct {
	IsLegendary bool `json:"is_legendary"`
	IsMythical  bool `json:"is_mythical"`
}
