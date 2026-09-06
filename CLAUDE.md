# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Working with the author

- The author writes the code. Default to explaining, pointing at the relevant spot, and suggesting an approach — do not make edits unless explicitly asked.
- The author is new to Go. When something breaks (compile errors, runtime panics, `go vet` complaints), teach: what the error means, why Go works this way, and how to fix it themselves. Don't just hand back corrected code. Keep explanations brief.
- This is a small side project, not enterprise software. Prefer the simplest thing that works. Push back on abstraction, layering, or dependencies that aren't earning their place.

## What this is

A minimal HTTP JSON API for tracking carrot-eating events ("Gulrot-tracker"). Single Go binary, SQLite storage, no framework — just `net/http` with Go 1.22+ method-prefixed route patterns.

## Commands

```bash
go run .                 # run locally on :8080 (loads .env, creates ./data/)
go build -o vetles-roetter .
go vet ./...              # only static checking configured; there is no linter setup
docker compose up --build # run containerized; needs CARROT_WRITE_KEY in the environment
```

There are no tests. `api-testing.sh` is an ad-hoc curl script and is currently stale — it hits `/` but the routes are under `/api/`.

## Configuration

- `CARROT_WRITE_KEY` — required for `POST` requests, compared in constant time against the raw `Authorization` header value (no `Bearer` prefix). Loaded from `.env` locally via `godotenv`; `.env` is gitignored.
- SQLite file lives at `data/carrotvault.sqlite`. The `data/` dir is created at startup and is gitignored; in Docker it is a named volume.

## Architecture

Three files, `package main`, wired together in `main.go`:

- **`carrot_repository.go`** — `CarrotRepository` wraps `*sql.DB`. Owns the schema (`initDatabase` runs `CREATE TABLE IF NOT EXISTS` on every boot — this is the only migration mechanism) and all SQL. The `Carrot` struct's json tags define the API response shape.
- **`carrot_handler.go`** — `CarrotHandler` holds a `CarrotRepository` by value and does HTTP encoding/decoding. Also defines `requireAuth`, a `http.HandlerFunc` middleware wrapper.
- **`main.go`** — opens the DB, constructs repository → handler, registers routes on the default `ServeMux`, starts the server.

Routes: `GET /api/` (list, public), `POST /api/` (add one carrot with server-set timestamp, auth-gated). `addCarrot` inserts `DEFAULT VALUES ... RETURNING`, so a carrot has no client-supplied fields.

Uses the pure-Go `modernc.org/sqlite` driver (`CGO_ENABLED=0`), which is why the Docker build produces a static distroless image.

## Conventions

- New resources follow the `*_repository.go` + `*_handler.go` pair, repository constructed first and passed into the handler in `main.go`.
- Handlers currently do not `return` after writing an error in some paths (e.g. repository errors fall through to also encode a response) — flag it to the author, don't silently fix.
