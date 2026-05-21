# Who Can Search YGO

Who Can Search YGO helps players find Yu-Gi-Oh! cards whose own effects can affect a selected target card.

The first MVP question is:

```txt
Which cards can add this selected card from the Deck to the hand?
```

The only supported public effect filter for this first delivery is `add_deck_to_hand`.

## Current Stack

- Backend: Go, `net/http`, `pgx/v5`, `goose`
- Database: PostgreSQL with `pg_trgm`
- API docs: OpenAPI served by the backend with Scalar
- Local infrastructure: Docker Compose

## Quick Start

Create local environment values:

```powershell
Copy-Item .env.example .env
```

Or on Unix-like shells:

```sh
cp .env.example .env
```

Start the local backend stack:

```sh
docker compose up --build backend
```

This starts PostgreSQL, runs migrations, syncs `backend/testdata/fixtures`, then starts the API.

Useful local URLs:

- API docs: http://localhost:18080/docs
- OpenAPI JSON: http://localhost:18080/openapi.json
- Health check: http://localhost:18080/healthz
- PostgreSQL: `localhost:15432`

## Backend Commands

Run commands from `backend/`:

```sh
make test
make migrate
make seed-fixtures
make run
```

The API expects the database schema and data to already exist. Migrations and fixture sync are separate operational commands, and Docker Compose runs them as one-shot services before the backend starts.

## Fixture Data

Local fixture data lives in:

```txt
backend/testdata/fixtures
```

The fixture sync command upserts those objects into PostgreSQL. This keeps local API behavior database-backed while still reusing the same data used by tests.

## API Shape

Public endpoints:

- `GET /cards?query=...`
- `GET /cards/{id}`
- `GET /cards/{id}/searchers`
- `GET /openapi.json`
- `GET /docs`
- `GET /healthz`

Public search requests use stored active selectors from the database. They must not call AI and must not scan every card text.
