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
- Avoid repeated logic; extract helpers when duplication starts to hide intent or creates multiple places to fix the same behavior.
- Prefer typed data structures over ad hoc maps.
- Keep functions small enough to test directly.
- Keep each function or method focused on one responsibility.
- When a function grows beyond a small, easily scannable block, especially around 5 to 10+ meaningful lines, actively check whether parsing, validation, mapping, persistence, or orchestration should be split into clearer helper functions.
- Prefer simple composition of small named functions over one method that mixes several levels of detail.
- Use comments only when they clarify non-obvious logic.
- Keep documentation, identifiers, comments, commits, database names, and UI text in English.

## Go Backend Rules

- Keep business logic out of HTTP handlers.
- Use `context.Context` for request-scoped operations.
- Return typed errors or wrap errors with useful context.
- Use table-driven tests for parser and matcher logic.
- Use `sqlc` query files as the source for generated database code.
- Use migrations as the source of truth for schema.
- Prefer standard library packages unless a dependency clearly reduces complexity.
- Do not introduce a heavy ORM.
- Do not use JavaScript or TypeScript in backend services.

## SQL And Database Rules

- Use explicit SQL.
- Add indexes with the query patterns they support.
- Use transactions for multi-step writes.
- Make import upserts idempotent.
- Keep raw upstream payloads for audit/debugging.
- Do not scan all card text during normal user search requests.

## Frontend Rules

- Build the actual search experience as the first screen.
- Keep UI dense, readable, and responsive.
- Make the experience immersive and polished through purposeful visual hierarchy, motion, and transitions.
- Use animations to clarify state changes or navigation, not to delay core search workflows.
- Creative UI treatment is welcome, but search, selection, result scanning, and reporting must remain fast and obvious.
- Use TypeScript for all frontend source.
- Keep API access in a small API layer.
- Use TanStack Query for server state.
- Test important user flows:
  - card search;
  - card selection;
  - searcher list states;
  - report submission.

## Testing Priorities

Backend mandatory tests:

- Source payload mapping.
- Import idempotency.
- Name normalization.
- Rule parsing fixtures.
- Criteria matching.
- Relationship preprocessing.
- User report creation.

Frontend mandatory tests:

- Search input behavior.
- Target card selection.
- Loading, empty, results, and error states.
- Report form submission.
