<!--
SPDX-FileCopyrightText: 2025 James Pond <james@cipher.host>
SPDX-FileCopyrightText: 2026 Renan Mateus Fernandes

SPDX-License-Identifier: CC0-1.0
-->

# Taskemon

> Self-hosted gamified task manager backend with hidden Pokémon rewards.

Taskemon turns everyday tasks into collectible encounters. Each task gets a hidden Pokémon reward when it is created. Complete the task to reveal the Pokémon, add it to your collection, and build your Pokédex over time.

This repository contains the **Taskemon backend server**, written in Go. It exposes a REST API backed by SQLite and is designed to support multiple clients, including a future web SPA, CLI, and Home Assistant dashboard.

---

## Project status

Taskemon is in early development. The backend API, database schema, and reward system are usable for development, but they may change before a stable `v1.0.0` release.

Current focus:

- backend correctness
- service, repository, API, CORS, and migration tests
- stable API response shapes
- SQLite persistence
- hidden reward and collection behavior

Planned clients such as the web SPA, CLI, and Home Assistant dashboard are not considered stable yet.

---

## Features

- REST API backend
- SQLite task persistence
- Task creation, listing, updating, completion, and deletion
- Open and completed task views
- Hidden Pokémon reward generated per task
- Reward reveal on task completion
- Pokémon collection tracking
- Normal and shiny Pokémon collection entries
- User statistics
- API response DTOs that avoid exposing internal models directly
- Backend tests for service, repository, API handlers, mappers, CORS, and migrations
- Designed for future SPA, CLI, and Home Assistant clients

---

## How it works

1. Create a task.
2. Taskemon secretly generates a Pokémon reward.
3. The task remains open without revealing the Pokémon.
4. Complete the task.
5. The hidden Pokémon is revealed.
6. The Pokémon is added to the user's collection.
7. User statistics are updated.

The Pokémon reward is intentionally hidden until completion so the reward acts as motivation instead of decoration.

---

## Architecture

```text
┌────────────────┐   ┌──────────────┐   ┌──────────────┐
│ Home Assistant │   │   Web SPA    │   │     CLI      │
└───────┬────────┘   └──────┬───────┘   └──────┬───────┘
        │                   │                  │
        └───────────────────┼──────────────────┘
                            │
                         REST API
                            │
                    ┌───────▼────────┐
                    │ Taskemon Server │
                    └───────┬────────┘
                            │
                 ┌──────────┴──────────┐
                 │                     │
              SQLite                PokéAPI
```

The backend is the source of truth. Clients should interact with Taskemon through the REST API.

---

## Requirements

- Go `1.25.0` or newer
- Network access to PokéAPI during startup and reward generation
- SQLite support through [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite)

Taskemon currently loads the Pokémon species count from PokéAPI during startup. If PokéAPI is unavailable, the server may fail to start.

---

## Running locally

Clone the repository:

```bash
git clone https://github.com/Renan-M-Fernandes/taskemon.git
cd taskemon
```

Install dependencies:

```bash
go mod tidy
```

Run the server:

```bash
go run ./cmd/taskemon
```

The server starts on port `8080` by default.

Health check:

```bash
curl http://localhost:8080/api/v1/health
```

---

## API endpoints

Authentication is planned for a future release. For now, `userID` is passed in the URL.

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/health` | Health check |
| `GET` | `/api/v1/tasks/{userID}` | List all tasks for a user |
| `GET` | `/api/v1/tasks/open/{userID}` | List open tasks for a user |
| `GET` | `/api/v1/tasks/completed/{userID}` | List completed tasks for a user |
| `POST` | `/api/v1/tasks/{userID}` | Create a task |
| `GET` | `/api/v1/tasks/{userID}/{taskID}` | Get one task |
| `PATCH` | `/api/v1/tasks/{userID}/{taskID}` | Update a task |
| `DELETE` | `/api/v1/tasks/{userID}/{taskID}` | Delete an open task |
| `POST` | `/api/v1/tasks/{userID}/{taskID}/complete` | Complete a task and reveal its reward |
| `GET` | `/api/v1/users/{userID}/collection` | List a user's Pokémon collection |
| `GET` | `/api/v1/users/{userID}/stats` | Get user statistics |

Example task creation:

```bash
curl -X POST http://localhost:8080/api/v1/tasks/ash \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Clean desk",
    "description": "Clear the desk before coding",
    "tag": "home"
  }'
```

---

## Database

Taskemon currently uses SQLite with the following core data areas:

- tasks
- hidden task rewards
- Pokémon collection entries
- user statistics

A local `taskemon.db` file is created when running the server.

---

## Development

Run formatting, dependency cleanup, tests, and checks before committing:

```bash
gofmt -w .
go mod tidy
go test ./...
go vet ./...
staticcheck ./...
go build ./...
```

If `staticcheck` is not installed, install it with:

```bash
go install honnef.co/go/tools/cmd/staticcheck@latest
```

Useful development checks:

```bash
git status --short
git diff
git diff --cached
```

---

## Roadmap

### Backend

- [x] REST API backend
- [x] SQLite persistence
- [x] Task CRUD
- [x] Hidden reward generation
- [x] Reward reveal on task completion
- [x] Pokémon collection
- [x] User statistics
- [x] API response DTOs
- [x] Backend test foundation
- [ ] Configuration file / environment configuration
- [ ] Authentication / API keys
- [ ] Transactions for multi-step operations
- [ ] Docker support
- [ ] OpenAPI documentation
- [ ] PokéAPI provider abstraction for deterministic tests

### Clients

- [ ] Web SPA
- [ ] CLI
- [ ] Home Assistant dashboard / integration
- [ ] Pokédex page
- [ ] Statistics dashboard
- [ ] Reward reveal animation

### Future ideas

- Thermal printer support
- Achievements
- Streak improvements
- Daily quests
- GitHub issue import
- Home Assistant automation hooks

---

## Versioning

Taskemon is pre-`v1.0.0`. Minor versions may still include breaking changes while the API and schema are evolving.

See [`CHANGELOG.md`](CHANGELOG.md) for release notes.

---

## Acknowledgements

Taskemon started as a fork of the original thermal-printer task manager by [James Pond](https://github.com/jamesponddotco/taskemon).

This version expands the original idea into a self-hosted gamified task manager backend with hidden Pokémon rewards, REST API support, SQLite persistence, backend tests, and future Home Assistant, SPA, and CLI clients.

---

## License

This project is distributed under the **EUPL-1.2** license.

See [`LICENSE.md`](LICENSE.md) and the [`LICENSES`](LICENSES) directory for details.