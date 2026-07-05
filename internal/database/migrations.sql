PRAGMA foreign_keys = ON;

----------------------------------------------------
-- Tasks
----------------------------------------------------

CREATE TABLE IF NOT EXISTS tasks (

    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         TEXT NOT NULL,
    title           TEXT NOT NULL,
    description     TEXT NOT NULL,
    completed       INTEGER NOT NULL DEFAULT 0,
    due_at          DATETIME,
    tag             TEXT NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at    DATETIME

);


CREATE INDEX IF NOT EXISTS idx_tasks_completed
ON tasks(completed);

CREATE INDEX IF NOT EXISTS idx_tasks_user
ON tasks(user_id);

CREATE INDEX IF NOT EXISTS idx_tasks_due_at
ON tasks(due_at);



----------------------------------------------------
-- Hidden rewards
----------------------------------------------------

CREATE TABLE IF NOT EXISTS task_rewards (

    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id             INTEGER NOT NULL UNIQUE,
    pokemon_id          INTEGER NOT NULL,
    pokemon_name        TEXT NOT NULL,
    sprite              TEXT, 
    rarity              INTEGER NOT NULL DEFAULT 1, --1 Common, 2 Uncommon, 3 Rare, 4 Legendary, 5 Mythical
    shiny               INTEGER NOT NULL DEFAULT 0,
    revealed            INTEGER NOT NULL DEFAULT 0,
    generated_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revealed_at         DATETIME,

    FOREIGN KEY(task_id)
        REFERENCES tasks(id)
        ON DELETE CASCADE

);


CREATE INDEX IF NOT EXISTS idx_rewards_revealed
ON task_rewards(revealed);



----------------------------------------------------
-- User collection
----------------------------------------------------

CREATE TABLE IF NOT EXISTS collection_entries (

    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id             TEXT NOT NULL,
    pokemon_id          INTEGER NOT NULL,
    pokemon_name        TEXT NOT NULL,
    count               INTEGER NOT NULL DEFAULT 1,
    rarity              INTEGER NOT NULL DEFAULT 1,
    shiny               INTEGER NOT NULL DEFAULT 0,    
    first_caught_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_caught_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, pokemon_id, shiny)

);



----------------------------------------------------
-- Statistics
----------------------------------------------------

CREATE TABLE IF NOT EXISTS user_statistics (

        user_id                TEXT PRIMARY KEY,
        tasks_completed        INTEGER NOT NULL DEFAULT 0,
        tasks_opened           INTEGER NOT NULL DEFAULT 0, 
        tasks_deleted          INTEGER NOT NULL DEFAULT 0,    
        pokemon_caught         INTEGER NOT NULL DEFAULT 0,
        shiny_caught           INTEGER NOT NULL DEFAULT 0,      
        unique_pokemon         INTEGER NOT NULL DEFAULT 0,
        current_streak         INTEGER NOT NULL DEFAULT 0,
        longest_streak         INTEGER NOT NULL DEFAULT 0
);