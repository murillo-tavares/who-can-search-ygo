# Architecture Context

## Stack

Backend:

- Go stable 1.x.
- `net/http` with a small router such as `chi` if needed.
- `pgx/v5` for PostgreSQL access.
- `sqlc` for type-safe SQL generation.
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
- JSONB for raw upstream card payloads and selector criteria.

Local infrastructure:

- Docker Compose.
- PostgreSQL service from the beginning.

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
    migrate/
  internal/
    app/
    config/
    httpapi/
    importer/
    postgres/
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

Keep backend and frontend separate. The backend must not depend on frontend tooling.

Use `.ai/schema.md` as the logical schema contract before migrations exist. Once migrations are implemented, migrations become the executable source of truth and should stay aligned with `.ai/schema.md`.

## Backend Boundaries

- HTTP handlers parse requests and write responses.
- Services own application use cases.
- Repository/database packages own persistence.
- Import mapping lives in `backend/internal/importer`.
- Selector normalization and matching live in `backend/internal/rules`.
- AI provider orchestration is not part of the MVP application.

Do not put business rules directly in handlers.

## Initial API

Public endpoints:

- `GET /healthz`
- `GET /cards?query=...`
- `GET /cards/{id}`
- `GET /cards/{id}/searchers`

Internal or CLI-oriented operations may be added later for imports and extraction writes, but public endpoints should stay small.

Response standards:

- JSON only.
- Use `snake_case` API fields.
- Errors include a stable code and human-readable message.
- Public endpoints must not expose raw upstream payloads.

## Background Work

Use CLI commands for MVP operational work:

- `backend/cmd/migrate`, if not using `goose` directly

The application works with the card and extraction data currently available in the database.

The schema includes:

- extraction version;
- effect code;
- per-card completion metadata;
- extracted effects;
- selectors;

Do not introduce queues, schedulers, Redis, or distributed workers until there is a real need.

## Query Model

Normal public lookups should be index-friendly and selector-driven:

- Card name search queries the `cards` table.
- Searcher results read active extracted effects with resolved stored selectors.
- Searcher results evaluate each stored selector against the selected target card at request time.
- Searcher results must not scan all card text.
- Searcher results must not call AI.
- The MVP does not store precomputed selector-to-card relationships.

Target search flow:

1. Load the selected target card by ID.
2. Load active extracted effects with `selector_status = resolved` for `add_deck_to_hand`.
3. Evaluate each effect selector against the target card using deterministic rule code.
4. Return source cards whose selectors match, including the matched `card_effects.effect_text`.

This is acceptable for the MVP because only stored active selectors are evaluated, not all card text.

If real-time selector evaluation becomes too slow, add a new decision before introducing precomputed relationships.
