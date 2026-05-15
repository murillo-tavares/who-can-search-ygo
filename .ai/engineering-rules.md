# Engineering Rules

## Agent Workflow

Before making changes:

1. Read `AGENTS.md`.
2. Read the relevant files under `.ai/`.
3. Inspect the current repository state.
4. Preserve unrelated user changes.
5. Prefer small, verifiable changes.

When changing behavior:

- Add or update tests.
- Run the narrowest meaningful test command first.
- Run broader checks when practical.
- Update `.ai/decisions.md` if an architectural decision changes.
- Update the relevant context file if product scope changes.

## General Code Rules

- Keep code explicit and boring.
- Avoid unnecessary abstractions in the MVP.
- Prefer typed data structures over ad hoc maps.
- Keep functions focused and testable.
- Use comments only when they clarify non-obvious logic.
- Keep documentation, identifiers, comments, commits, database names, and UI text in English.
- Preserve unrelated user changes.

## Go Backend Rules

- Keep business logic out of HTTP handlers.
- Use `context.Context` for request-scoped operations.
- Return typed errors or wrap errors with useful context.
- Use table-driven tests for selector matching.
- Use `sqlc` query files as the source for generated database code.
- Use migrations as the source of truth for schema.
- Prefer standard library packages unless a dependency clearly reduces complexity.
- Do not introduce a heavy ORM.

## SQL And Database Rules

- Use explicit SQL.
- Add indexes for real query patterns.
- Use transactions for multi-step writes.
- Make imports idempotent.
- Keep raw upstream payloads for audit/debugging.
- Track completed AI processing on `cards.ai_processing` even when AI extracts no supported effect.
- Do not scan all card text during public user requests.
- Do not use AI inference during public user requests.
- Do not add precomputed relationship tables without a new decision.

## Selector Rules

- Selector matching must be deterministic.
- Selector JSON must be normalized before hashing.
- Selector JSON must be a structured expression AST, not raw SQL.
- Public matching uses stored active selectors with `selector_status = resolved`.
- Public matching evaluates selectors against one selected target card at request time.
- Unsupported selector fields or operators must fail closed.
- Text-listing requirements such as `lists the card "Dark Magician"` should be represented with `mentions contains "Dark Magician"` when source data supports it.
- Ambiguous AI output should use `selector_status = unresolved` or be skipped entirely.

## AI Extraction Rules

- The first AI extraction process is external to this application.
- AI extraction is assistive, not authoritative.
- External extraction should write public-matchable effects only when `selector_status = resolved`.
- AI should generate structured selectors, not direct final result lists.
- Public requests must never call AI models.
- Prefer conservative extraction over speculative interpretation.

## Frontend Rules

- Build the actual search experience as the first screen.
- Keep UI readable, responsive, and direct.
- Use TypeScript for all frontend source.
- Keep API access in a small API layer.
- Use TanStack Query for server state.
- Public UI should only display active extracted effects with resolved selectors.
- Public UI should not expose selector internals unless internal tooling is explicitly added.

Test important user flows:

- card search;
- target card selection;
- effect-kind filter display;
- source card result states;
- empty and error states.

## Testing Priorities

Backend:

- source payload mapping;
- import idempotency;
- name normalization;
- selector normalization;
- selector hashing;
- selector matching;
- AI processing metadata queries;
- public searcher query behavior.

Frontend:

- search input behavior;
- target card selection;
- loading, empty, results, and error states;
- effect-kind selection.

Tests must not depend on external AI APIs.
