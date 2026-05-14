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
- Relationship rows match reusable selectors to target cards. Action scope stays on extracted effect actions.

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

## ADR-006: Use Extracted Effects As Preprocessing Input

Status: accepted

Decision:

Parse card text into supported extracted effect records, and make relationship preprocessing consume those extracted effects as its candidate input.

Rationale:

- Sync should remain idempotent and preserve cards for future reevaluation.
- Preprocessing should not waste work reparsing card text during relationship generation.
- A card that cannot search other cards may still be a valid target card.

Consequences:

- Extracted effects need parser version, review status, text segments, action tags, and `actions_json` objects that reference reusable selectors.
- Public searcher queries rely on accepted precomputed relationships, not on card text.
- Re-extraction must be able to add, update, or remove effect rows when card text changes, parser versions change, or supported action families expand.

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

Persist parsed effects separately from selector-target relationships. Store text segments and ordered parsed actions in `extracted_effects`; store reusable target criteria in `card_selectors`.

Rationale:

- Parser output needs review, versioning, and regeneration.
- Public search should reuse precomputed selector-target rows.
- JSONB stores audit shape; Go structs remain the rule-engine model.

Consequences:

- Rule parsing must produce typed effects that serialize to `text_segments_json` and `actions_json`.
- Rule parsing should use PSCT punctuation before matching supported effect patterns.
- Relationship preprocessing should derive target rows from accepted selectors referenced by accepted `actions_json` entries with `action_kind = move_card`.

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

## ADR-011: Keep Application Enums In Code

Status: accepted

Decision:

Store controlled vocabulary fields as text columns. Define allowed enum values in application code rather than creating database lookup tables for finite application vocabularies, including actions, zones, match kinds, statuses, and report types.

Rationale:

- These values are part of application behavior, parser support, and workflow logic.
- Application enums keep validation close to the business logic that uses each value.
- Avoids extra lookup tables that do not add value for the MVP.

Consequences:

- The application layer must validate enum values before writes.
- Migrations do not need lookup tables for actions, zones, match kinds, statuses, report types, or similar small finite workflow values.
- If an admin workflow later needs editable labels, ordering, or metadata, this can be revisited.

## ADR-012: Use Build-Like Rule Version Identifiers

Status: accepted

Decision:

Use explicit build-like identifiers for parser and relationship-preprocessor versions, starting with `rules-1` and `relationships-1`.

Rationale:

- Rule output needs traceability more than package-style semantic versioning.
- Parser and relationship preprocessing can evolve independently.
- Short named identifiers are easy to store, compare, and read in audit data.

Consequences:

- `parser_version` values should look like `rules-1`, `rules-2`, and so on.
- `relationship_version` values should look like `relationships-1`, `relationships-2`, and so on.
- Any change that can alter extracted effects or generated relationships should bump the relevant identifier.
