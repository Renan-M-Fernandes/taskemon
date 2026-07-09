<!--
SPDX-FileCopyrightText: 2025 James Pond <james@cipher.host>
SPDX-FileCopyrightText: 2026 Renan Mateus Fernandes

SPDX-License-Identifier: CC0-1.0
-->

# Taskemon

> Self-hosted gamified task manager backend with hidden Pokémon rewards and optional thermal task-card printing.

Taskemon turns everyday tasks into collectible encounters. Each task gets a hidden Pokémon reward when it is created. Complete the task to reveal the Pokémon, add it to your collection, and build your Pokédex over time.

This repository contains the **Taskemon backend server**, written in Go. It exposes a REST API backed by SQLite and is designed to support multiple clients, including a future web SPA, CLI, and Home Assistant dashboard.

---

## Project status

Taskemon is in early development. The backend API, database schema, printer behavior, and reward system are usable for development, but they may change before a stable `v1.0.0` release.

Current focus:

- backend correctness
- service, repository, API, CORS, config, printer, and migration tests
- stable API response shapes
- SQLite persistence
- hidden reward and collection behavior
- optional ESC/POS thermal printer support
- simple configuration for local self-hosted setups

Planned clients such as the web SPA, CLI, and Home Assistant dashboard are not considered stable yet.

---

## Features

- REST API backend
- SQLite task persistence
- Task creation, listing, updating, completion, deletion, and printing
- Open and completed task views
- Hidden Pokémon reward generated per task
- Reward reveal on task completion
- Pokémon collection tracking
- Normal and shiny Pokémon collection entries
- User statistics
- API response DTOs that avoid exposing internal models directly
- Optional ESC/POS thermal printer support
- Native thermal ticket layout with title, description, QR code, and tag footer
- Printer QR modes for temporary Pokémon reveal and future task-completion flows
- Configurable printer transport, USB IDs, QR behavior, cut command, line widths, and layout options
- Linux helper script for optional USB printer udev setup
- Backend tests for service, repository, API handlers, mappers, CORS, config, printer behavior, and migrations
- Designed for future SPA, CLI, and Home Assistant clients

---

## How it works

1. Create a task.
2. Taskemon secretly generates a Pokémon reward.
3. The task remains open without revealing the Pokémon.
4. Optionally print a physical task card.
5. Complete the task.
6. The hidden Pokémon is revealed.
7. The Pokémon is added to the user's collection.
8. User statistics are updated.

The Pokémon reward is intentionally hidden until completion so the reward acts as motivation instead of decoration.

During early development, the printer can use a temporary QR mode that points directly to the generated Pokémon page. This is useful before the official web UI exists. Later, the QR flow is expected to open a Taskemon page where the user can complete the task and reveal the reward.

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
        ┌───────────────────┼───────────────────┐
        │                   │                   │
     SQLite              PokéAPI          ESC/POS Printer
```

The backend is the source of truth. Clients should interact with Taskemon through the REST API.

Printer support is an optional output integration. Taskemon can run normally without a printer.

---

## Requirements

- Go `1.25.0` or newer
- Network access to PokéAPI during startup and reward generation
- SQLite support through [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite)
- Optional: ESC/POS-compatible thermal printer for physical task cards
- Optional on Linux: USB access to the printer device when using direct USB printing

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

Optional: create a local config file:

```bash
cp config.example.json config.json
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

## Configuration

Taskemon can run with defaults, but local setups should use `config.json` when enabling printer support or changing server/database settings.

Example printer section:

```json
{
  "printer": {
    "enabled": true,
    "transport": "usb",
    "vendorID": "0x0418",
    "productID": "0x5011",
    "endpoint": 1,
    "qrMode": "pokemon_placeholder",
    "baseURL": "http://localhost:8080",
    "charsPerLine": 64,
    "feedLinesBeforeCut": 4,
    "qrSize": 4,
    "qrCorrection": "H",
    "cutCommand": "gs_v_0"
  }
}
```

Common printer values:

| Option | Description |
|--------|-------------|
| `enabled` | Enables or disables printer support |
| `transport` | Printer transport, currently focused on `usb` |
| `vendorID` | USB vendor ID, for example `0x0418` |
| `productID` | USB product ID, for example `0x5011` |
| `endpoint` | USB output endpoint, usually `1` |
| `qrMode` | QR behavior for printed cards |
| `baseURL` | Base URL used by Taskemon QR flows |
| `charsPerLine` | Text layout width for native ESC/POS printing |
| `feedLinesBeforeCut` | Blank lines before cutting |
| `qrSize` | Requested native QR module size, if supported by the printer |
| `qrCorrection` | QR error correction level, usually `L`, `M`, `Q`, or `H` |
| `cutCommand` | Cut command variant for ESC/POS printers |

### Printer QR modes

| Mode | Description |
|------|-------------|
| `pokemon_placeholder` | Temporary development mode. The printed QR points to the generated Pokémon page so the reward can be viewed before the official UI exists. |
| `task_completion` | Future-facing mode. The printed QR points to Taskemon so a UI can complete the task and reveal the reward. |

---

## Thermal printer setup

Taskemon prints through direct ESC/POS commands. The current implementation focuses on USB thermal printers using `gousb`.

Printer support is optional. If printing is disabled, the backend task API and reward system still work normally.

### Linux USB printer access

Taskemon tries to automatically detach the Linux `usblp` kernel driver when opening a USB printer. On many systems this is enough and no manual setup is required.

Some Linux setups still block direct USB access because the current user does not have permission to open the USB printer device.

sudo modprobe -r usblp usually works

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
| `POST` | `/api/v1/tasks/{userID}/{taskID}/print` | Print a physical task card |
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

Example print request:

```bash
curl -X POST http://localhost:8080/api/v1/tasks/ash/1/print
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
- [x] Configuration file support
- [x] Optional thermal printer support
- [x] Printer API endpoint
- [ ] Simple authentication / API keys
- [ ] Stable authenticated route shape
- [ ] Transactions for multi-step operations
- [ ] Docker support
- [ ] OpenAPI documentation
- [ ] PokéAPI provider abstraction for deterministic tests

### Printer

- [x] Native ESC/POS text ticket layout
- [x] USB printer transport
- [x] Configurable QR behavior
- [x] Configurable cut/feed behavior
- [ ] Linux udev helper script
- [ ] Printer doctor command
- [ ] TCP/network printer transport
- [ ] Optional checklist/list printing
- [ ] More portable printer compatibility notes

### Clients

- [ ] Temporary development SPA
- [ ] Official web SPA
- [ ] CLI
- [ ] Home Assistant dashboard / integration
- [ ] Pokédex page
- [ ] Statistics dashboard
- [ ] Reward reveal animation

### Future ideas

- Achievements
- Streak improvements
- Daily quests
- GitHub issue import
- Home Assistant automation hooks
- Checklist tasks
- Recurring tasks

---

## Versioning

Taskemon is pre-`v1.0.0`. Minor versions may still include breaking changes while the API and schema are evolving.

See [`CHANGELOG.md`](CHANGELOG.md) for release notes.

---

## Acknowledgements

Taskemon started as a fork of the original thermal-printer task manager by [James Pond](https://github.com/jamesponddotco/taskemon).

This version expands the original idea into a self-hosted gamified task manager backend with hidden Pokémon rewards, REST API support, SQLite persistence, backend tests, optional thermal printing, and future Home Assistant, SPA, and CLI clients.

---

## License

This project is distributed under the **EUPL-1.2** license.

See [`LICENSE.md`](LICENSE.md) and the [`LICENSES`](LICENSES) directory for details.
