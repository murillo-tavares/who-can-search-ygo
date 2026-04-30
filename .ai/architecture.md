# Architecture Context

## Stack

Backend:

- Go 1.26.x or newer stable Go 1.x.
- `net/http` with a small router such as `chi` if needed.
- `pgx/v5` for PostgreSQL access.
- `sqlc` for type-safe SQL code generation.
- `goose` for migrations.
- `log/slog` for logging.
- Go `testing` for backend tests.

Frontend:

- React.
- Vite.
- TypeScript.
- TanStack Query for API state.
- React Router only when multiple routes are needed.
- Vitest and Testing Library.

Database:

- PostgreSQL.
- `pg_trgm` for tolerant card-name lookup.
- JSONB for upstream raw payload audit/debugging.
- Indexed extracted effect fields, such as action, source zone, destination zone, parser version, and status.
- Indexed relationship fields for target card, action, source zone, destination zone, and status.

Local infrastructure:

- Docker Compose.
- PostgreSQL service from the beginning.
- Redis only later if CLI/manual preprocessing is no longer enough.

## Repository Shape

Start with this structure:

```txt
.ai/
  architecture.md
  data-and-rules.md
  decisions.md
  engineering-rules.md
  open-questions.md
  product.md
backend/
  cmd/
    api/
    sync-cards/
    preprocess-searches/
  internal/
    app/
    config/
    httpapi/
    importer/
    postgres/
    preprocessing/
    rules/
    service/
  migrations/
  queries/
  testdata/
frontend/
  src/
    api/
    components/
    pages/
    styles/
    test/
docker-compose.yml
Makefile
AGENTS.md
```

Keep backend and frontend separate. Shared generated API types may be introduced later, but the backend must not depend on frontend tooling.

## Backend Boundaries

- HTTP handlers parse requests and write responses.
- Services own application use cases.
- Repositories/database packages own persistence.
- Rule parsing and target matching live in `backend/internal/rules`.
- Import mapping lives in `backend/internal/importer`.
- Preprocessing orchestration lives in `backend/internal/preprocessing`.

Do not put business rules directly in handlers.

## Initial API

Public endpoints:

- `GET /healthz`
- `GET /cards?query=...`
- `GET /cards/{id}`
- `GET /cards/{id}/searchers`
- `POST /reports`

Response standards:

- JSON only.
- Use `snake_case` API fields.
- Errors include a stable code and human-readable message.
- Public endpoints must not expose raw upstream payloads.
- Searcher results only need to indicate matching cards in the MVP; they should not include matched reasoning unless a later version requires it.

## Background Work

Use CLI commands for MVP:

- `backend/cmd/sync-cards`
- `backend/cmd/preprocess-searches`
- `backend/cmd/migrate`, if not using `goose` directly

Do not introduce a queue or scheduler until there is a real need for distributed or recurring background work.

Deployment should use the current prepared database state for the MVP. Card sync and search preprocessing are run separately, locally or through an explicit external workflow, not automatically as part of deploy.

## Query Model

Normal public lookups should be index-friendly and relationship-driven:

- Card name search may query the card table.
- Searcher results must read precomputed relationships, not card text.
- The MVP searcher query defaults to relationships with action `add`, source `deck`, destination `hand`, and `accepted` status.
- Future filters should use persisted action, source zone, and destination zone fields, not reparsed text.
- Relationship preprocessing should read supported extracted effects rather than scanning card text.
- Imported cards remain available as target cards regardless of whether they currently have extracted effects.
- User reports may be anonymous, but the API must include basic spam protection such as rate limiting, validation, or another lightweight abuse-control mechanism.
