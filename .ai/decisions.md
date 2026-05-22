# Architecture Decisions

This file records durable decisions.

Keep entries short and update them when decisions change.

## ADR-001: Use Go For Backend

Status: accepted

Decision:

Use Go for the backend API, import commands, migrations wrapper, and deterministic selector matching.

Rationale:

- Fast runtime and low operational overhead.
- Simple compiled deployment.
- Explicit code style that is easy for humans and AI agents to maintain.
- Strong fit for deterministic rule evaluation.

## ADR-002: Use React For Frontend

Status: accepted

Decision:

Use React, Vite, and TypeScript for the frontend.

Rationale:

- Fast MVP iteration.
- Strong ecosystem for search UI and async API state.
- TypeScript improves API contract safety.

## ADR-003: Use PostgreSQL As Primary Database

Status: accepted

Decision:

Use PostgreSQL as the application database.

Rationale:

- Card data and extracted effects are relational.
- `pg_trgm` supports tolerant card-name search.
- JSONB supports raw upstream card payloads and selector criteria.

## ADR-004: Use Stored AI Extraction Data

Status: accepted

Decision:

The application works with extraction data currently available in the database.

The application stores:

- cards;
- supported effect codes;
- per-card AI processing metadata;
- extracted effects and selectors.

Rationale:

- Keeps MVP application focused on card search behavior.
- Avoids coupling the application to a specific extraction runner.
- Still preserves auditability and reprocessing support.

## ADR-005: Track Completed AI Processing On Cards

Status: accepted

Decision:

Track latest completed AI processing in `cards.ai_processing`.

Store effect code, extraction version, processed timestamp, and result count. Do not store extraction runs.

Rationale:

- Supports resume after interruption.
- Shows which cards were already processed.
- Allows reprocessing with new extraction versions.
- Identifies newly imported cards needing extraction.
- Avoids tables that only record operational history not needed by the MVP.

Consequences:

- Missing metadata means processing is still needed.
- Version mismatch means reprocessing is needed.
- Historical run counts are intentionally unavailable.

## ADR-006: Use Active Resolved Effects For Public Results

Status: accepted

Decision:

Public search uses active extracted effects with `selector_status = resolved`.

Unresolved effects may be stored for audit/debugging when a supported effect is detected but an exact selector cannot be defined.

Ignored effects may be stored when the parser knows the rule but the effect should not be exposed publicly, such as generic selectors that reduce to any card, any Monster, any Spell, or any Trap.

Rationale:

- Incorrect records can be disabled with `is_active = false`.
- Ignored and unresolved records do not appear in public search results.

## ADR-007: Store Image URLs Only

Status: accepted

Decision:

Store upstream image URLs and hotlink images in the MVP.

Rationale:

- Avoids image storage and processing complexity.
- Keeps card import focused on metadata and selector fields.

## ADR-008: Use Structured Selector Expressions

Status: accepted

Decision:

Represent selectors as structured JSON expression trees with comparison and logical nodes.

Example semantics:

```txt
name = 'Dark Magician' OR (card_type IN ('Spell', 'Trap') AND mentions CONTAINS 'Dark Magician')
```

Rationale:

- Handles mixed `and`/`or` selector logic without special-case fields.
- Maps naturally to player-readable conditions and SQL-like semantics.
- Remains safer than storing raw SQL.
- Keeps selector matching deterministic and validateable.

Consequences:

- Selector validation must reject unknown fields, operators, and invalid value types.
- Selector normalization and hashing must canonicalize expression trees.
- Boolean nesting has no fixed depth limit, but malformed expression nodes must fail closed.

## ADR-009: Use Stable Effect Codes

Status: accepted

Decision:

Use stable application-defined `effect_code` values.

Rationale:

- MVP has one supported filter: `add_deck_to_hand`.
- Text codes keep processing metadata and effects simple.

Consequences:

- Application code validates supported effect codes.

## ADR-010: Use Kebab-Case Extraction Versions

Status: accepted

Decision:

Use extraction versions in this format:

```txt
<effect-code-kebab>-v<number>
```

Initial version:

```txt
add-deck-to-hand-v1
```

Rationale:

- Human-readable.
- Tied to the effect code being extracted.
- Simple to increment when extraction logic changes.

## ADR-011: Run Database Migrations As A Separate Process

Status: accepted

Decision:

Database migrations run through the dedicated `backend/cmd/migrate` command. Docker Compose runs this command as a one-shot service before the API starts. The API does not apply migrations during startup.

Rationale:

- Preserves a versioned migration history through `goose`.
- Keeps API startup focused on serving traffic.

Consequences:

- The API may fail at runtime if it starts against an unmigrated database.
- Local Compose startup fails before the API starts if migrations fail.
- A PostgreSQL advisory lock protects concurrent startup attempts from running migrations at the same time.

## ADR-012: Serve OpenAPI Documentation From The API

Status: accepted

Decision:

Expose the OpenAPI document at `GET /openapi.json` and render interactive API documentation at `GET /docs` with Scalar.

Rationale:

- Keeps endpoint documentation close to the backend implementation.
- Gives humans an interactive way to inspect and try public endpoints.
- Avoids adding a separate documentation service for the MVP.

## ADR-013: Use Environment Files For Local Configuration

Status: accepted

Decision:

Keep runtime configuration in environment variables and `.env` files.

The backend loads local `.env` files during startup. The backend Makefile does not set runtime environment variables inline.

Docker Compose can run both PostgreSQL and the backend, using `.env` values with safe local defaults where possible.

Rationale:

- Keeps local configuration easy to change without editing Make targets.
- Makes Docker and non-Docker runs use the same configuration names.

Consequences:

- Developers should create a local `.env` from `.env.example` when using `make run`.
- Docker Compose uses an internal `DATABASE_URL` pointing at the `postgres` service.

## ADR-014: Sync Local Fixtures Into PostgreSQL

Status: accepted

Decision:

For local development, synchronize fixture JSON data into PostgreSQL through the dedicated `backend/cmd/seed-fixtures` command.

Docker Compose runs fixture synchronization as a one-shot service after migrations and before the API starts. Sync operations are idempotent and reuse the same fixture files used by backend tests.

Rationale:

- Keeps local API behavior database-backed.
- Makes local testing closer to production query behavior.
- Avoids maintaining separate local fixture data and database seed data.

Consequences:

- Fixture JSON changes affect both local seed data and fixture-backed tests.
