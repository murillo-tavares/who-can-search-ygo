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
- JSONB for upstream raw payloads, aliases, text segments, parsed actions, and selector criteria.
- `card_selectors` for reusable target-card criteria.
- `search_relationships` for precomputed selector-to-target matches.
- Indexed status/version/hash fields for effects, selectors, and relationships.

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
  schema.md
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

Use `.ai/schema.md` as the logical schema contract before migrations exist. Once migrations are implemented, they become the executable source of truth and should stay aligned with `.ai/schema.md`.

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
- Searcher results must read precomputed selector-target relationships, not card text.
- The MVP searcher query starts from `accepted` relationships for the target card, then finds accepted extracted effects whose `actions_json` references the relationship selector with `action_kind = move_card`, `verb = add`, `config.from = deck`, and `config.to = hand`.
- Relationship preprocessing should match target cards from `card_selectors`, not scan raw card text.
- Imported cards remain available as target cards regardless of whether they currently have extracted effects.
- User reports may be anonymous, but the API must include basic spam protection such as rate limiting, validation, or another lightweight abuse-control mechanism.
