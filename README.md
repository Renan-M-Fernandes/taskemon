# Taskemon

> A self-hosted, Pokémon-inspired gamified task management platform.

Taskemon transforms everyday tasks into collectible adventures.

Every task you create secretly generates a Pokémon encounter. Complete the task to reveal the hidden Pokémon, add it to your collection, and build your Pokédex over time.

This repository contains the **official Taskemon backend**, written in Go.

---

## Features

- REST API
- SQLite database
- Task CRUD
- Hidden Pokémon rewards
- Pokémon collection (Pokédex)
- User statistics
- Home Assistant integration
- Designed for future SPA and CLI clients

---

## Roadmap

### Backend

- [x] Task CRUD
- [x] Hidden reward generation
- [x] Reward reveal on task completion
- [x] Pokédex collection
- [x] Statistics
- [ ] Authentication
- [ ] Configuration file
- [ ] Docker support
- [ ] Swagger / OpenAPI
- [ ] Unit tests

### Frontend

- [ ] Web SPA
- [ ] Home Assistant dashboard
- [ ] Pokédex page
- [ ] Statistics dashboard
- [ ] Reward animations

---

## Architecture

```
                ┌──────────────┐
                │ Home Assistant│
                └──────┬───────┘
                       │
                ┌──────▼───────┐
                │ Web SPA      │
                └──────┬───────┘
                       │
                ┌──────▼───────┐
                │ CLI          │
                └──────┬───────┘
                       │
                REST API
                       │
             ┌─────────▼─────────┐
             │   Taskemon Server │
             └─────────┬─────────┘
                       │
                   SQLite
```

The backend is designed so multiple clients can interact with the same API.

---

## Current Database

- Tasks
- Hidden task rewards
- Pokémon collection
- User statistics

---

## API

Current endpoints:

| Method | Endpoint |
|--------|----------|
| GET | `/api/v1/health` |
| GET | `/api/v1/tasks/{userID}` |
| POST | `/api/v1/tasks/{userID}` |
| GET | `/api/v1/tasks/{userID}/{taskID}` |
| PATCH | `/api/v1/tasks/{userID}/{taskID}` |
| DELETE | `/api/v1/tasks/{userID}/{taskID}` |
| POST | `/api/v1/tasks/{userID}/{taskID}/complete` |
| GET | `/api/v1/users/{userID}/collection` |
| GET | `/api/v1/users/{userID}/stats` |

Authentication is planned for a future release.

---

## Running

Clone the repository

```bash
git clone https://github.com/Renan-M-Fernandes/taskemon.git
```

Install dependencies

```bash
go mod tidy
```

Run

```bash
go run ./cmd/taskemon
```

The server starts on port **8080** by default.

---

## Tech Stack

- Go
- SQLite
- REST
- Home Assistant
- PokéAPI

---

## Project Status

Taskemon is currently under active development.

The backend is considered an early preview and the API may change before a stable v1 release.

---

## Inspiration

Taskemon began as a small thermal-printer task manager created by
[@jamesponddotco](https://github.com/jamesponddotco).

This project expands that original idea into a self-hosted gamified task management platform with a REST API, future web frontend, Home Assistant integration and Pokémon-inspired progression.

---

## License

This project is distributed under the **EUPL-1.2** license.

See the LICENSE files for details.