# Architecture Decisions

This file records durable decisions. Keep entries short and update them when decisions change.

## ADR-001: Use Go For Backend

Status: accepted

Decision:

Use Go for backend API, import commands, and preprocessing commands.

Rationale:

- Fast runtime and low operational overhead.
- Simple deployment as compiled binaries.
- Good standard library for HTTP, concurrency, CLI commands, and tests.
- Explicit code style that is easy for humans and AI agents to maintain.
- Strong fit for deterministic rule parsing and relationship preprocessing.

Consequences:

- Frontend may still use TypeScript.
- SQL and schema boundaries should stay explicit.

## ADR-002: Use PostgreSQL As Primary Database

Status: accepted

Decision:

Use PostgreSQL as the application database.

Rationale:

- Search relationships are naturally relational.
- PostgreSQL supports strong indexing for lookup and joins.
- `pg_trgm` supports tolerant card-name search.
- JSONB can store raw upstream payloads.

Consequences:

- No separate search engine in the MVP.
- Query performance should be solved first with indexes and precomputed tables.

## ADR-003: Precompute Search Relationships

Status: accepted

Decision:

Normal user searches must read precomputed relationship rows instead of scanning card text.

Rationale:

- User queries must be fast.
- Rule parsing can be expensive and should run outside normal requests.
- Precomputed rows allow review, status, and versioning.

Consequences:

- Sync/preprocessing must be part of the core backend workflow.
- Relationship rows need action, source zone, destination zone, preprocessing version, and status.

## ADR-004: Use CLI Workflows For MVP Sync And Preprocessing

Status: accepted

Decision:

Use manual CLI commands for card synchronization and search relationship preprocessing in the MVP. Deployment uses the current prepared database state and does not automatically run sync or preprocessing.

Rationale:

- Lower complexity than queues, workers, admin UI, or schedulers.
- Matches the MVP requirement for a technical/admin workflow.
- Easy to run locally and in deployment jobs.

Consequences:

- No Redis or queue at project start.
- A future admin UI can wrap the same service logic.
- Sync and preprocessing may run locally or through an explicit external workflow outside deploy.

## ADR-005: Choose Initial Data Source

Status: proposed

Decision:

Choose a single initial bulk card-data ingestion source before heavy importer implementation.

Rationale:

- The importer and payload mapping should be designed around a concrete source.
- Criteria matching depends on source field quality and consistency.
- The project still needs to compare available Yu-Gi-Oh! data sources.

Consequences:

- Keep source-specific code isolated in the importer.
- YGOPRODeck API v7 remains a candidate source to evaluate.
- Keep raw payloads for audit/debugging.
- Treat official Konami card text as canonical when wording conflicts appear.

## ADR-006: Use Extracted Effects As Preprocessing Input

Status: accepted

Decision:

Parse card text into supported extracted effect records, and make relationship preprocessing consume those extracted effects as its candidate input.

Rationale:

- Sync should remain idempotent and preserve cards for future reevaluation.
- Preprocessing should not waste work reparsing card text during relationship generation.
- A card that cannot search other cards may still be a valid target card.

Consequences:

- Extracted effects need action, source zone, destination zone, the original card text fragment, parser version, status, and normalized condition data.
- Public searcher queries rely on accepted precomputed relationships, not on card text.
- Re-extraction must be able to add, update, or remove effect rows when card text changes, parser versions change, or supported action/source/destination flows expand.

## ADR-007: Let The Selected Source Define The Initial Card Universe

Status: accepted

Decision:

The MVP imports and exposes the available cards returned by the selected API or local test dataset. Do not add special TCG, OCG, Master Duel, upcoming-card, Rush Duel, Speed Duel, Speed Spell, skill-card, or variant filtering until a scope change requires it.

Rationale:

- The initial data source is still unresolved.
- The MVP should avoid premature filtering rules that depend on source-specific fields.
- The core rule engine can be validated against a small curated dataset first.

Consequences:

- The selected source or fixture dataset is responsible for providing the cards considered by the app.
- Variant-specific exclusions can be added later as explicit scope changes.
- Import mapping should preserve source metadata when available, but the MVP does not need out-of-scope detection fields before a source is selected.

## ADR-008: Persist Extracted Effects Separately From Relationships

Status: accepted

Decision:

Persist extracted card effects as structured records, separate from precomputed target relationships. Each card can have multiple effect records, each scoped by action, source zone, destination zone, parser version, status, and normalized condition data.

Rationale:

- The system needs to explain and reevaluate why a card can search another card.
- The MVP supports Deck-to-hand add effects, but the schema should be ready for future filters such as adding from the GY, sending from Deck to the GY, Special Summoning from the Deck, banishing from the Deck, or destroying from another zone to the GY.
- Action is an independent persisted dimension. Initial action values include `add`, `send`, `banish`, `special_summon`, and `destroy`.
- Each extracted effect should keep the original text fragment that produced it, with location metadata when practical, so parser output can be debugged and reviewed later.
- JSONB can store audit-friendly rule expressions while Go code keeps typed parser and matcher structs.
- Structured conditions make future review, admin tools, and rule version migrations easier.

Consequences:

- Rule parsing must produce typed effects that can serialize to normalized JSON.
- Rule parsing should use PSCT punctuation to separate activation condition, activation text, and resolving effect before matching supported effect patterns.
- Relationship preprocessing should derive target rows from extracted effects, not directly from raw card text.
- Non-default action/source/destination flows may be stored for audit or future use, but should not generate public relationships until explicitly supported.

## ADR-009: Use Latest Raw Payload Only For MVP

Status: accepted

Decision:

Store the latest raw upstream payload for each imported card, not full historical source snapshots, unless a later requirement needs history.

Rationale:

- Normal queries should use normalized fields and precomputed relationships.
- Historical snapshots are extra storage and complexity that do not improve MVP lookup performance.
- Latest raw payloads are enough for audit/debugging during early implementation.

Consequences:

- Import upserts replace the stored raw payload.
- Any future historical audit table should be separate from public query paths.

## ADR-010: Keep MVP Public Responses Minimal

Status: accepted

Decision:

Use `snake_case` API fields. Public searcher results return matching cards only and do not include matched reasoning in the MVP. User reports may be anonymous but must have lightweight spam protection.

Rationale:

- `snake_case` matches the current API default.
- Reasoning can be added later without blocking the core search experience.
- Anonymous reports reduce friction, but public endpoints need abuse controls.

Consequences:

- API response contracts should avoid exposing internal rule explanations for now.
- Report endpoints need validation and basic spam mitigation.
