# AGENTS.md

This is the entry point for AI agents working on Who Can Search YGO.

Keep this file short. Detailed context lives in `.ai/` so agents can load only what they need.

## Read Order

Before making changes, read:

1. `.ai/product.md`
2. `.ai/architecture.md`
3. `.ai/data-and-rules.md`
4. `.ai/schema.md`
5. `.ai/engineering-rules.md`
6. `.ai/decisions.md`
7. `.ai/open-questions.md`

## Project Mission

Who Can Search YGO identifies which Yu-Gi-Oh! cards can add a selected target card from the Deck to the hand.

```txt
Who can get this Yu-Gi-Oh! card from the Deck to the hand?
```

The primary/default experience answers which cards can add a target card from the Deck to the hand, beginning with official English wording patterns like:

```txt
Add ... from your Deck to your hand
```

Cards may have multiple extracted effect records. Effects store parsed actions in `actions_json`; target matching is reused through `card_selectors` and `search_relationships`.

Rule parsing should read card text and persist supported extracted effects. Relationship preprocessing should consume those extracted effects instead of discovering candidates by scanning card text.

Normal user queries must use precomputed relationships. Do not scan all card text during request handling.

## Current Defaults

- Backend: Go.
- Frontend: React, Vite, TypeScript.
- Database: PostgreSQL with `pg_trgm`.
- Data source: unresolved; see `.ai/open-questions.md`.
- Sync/preprocessing: manual CLI commands for MVP.
- Public relationship results: show only `accepted` relationships by default.

## Non-Negotiables

- Keep backend and frontend as separate subprojects.
- Keep rule parsing isolated, deterministic, testable, and versioned.
- Keep extracted effects separate from precomputed target relationships.
- Keep imports and preprocessing idempotent.
- Preserve unrelated user changes.
- Use English for documentation, code, identifiers, database names, API fields, comments, commits, UI text, and internal terminology.

## Where To Change Context

- Product scope changes: update `.ai/product.md`.
- Stack or architecture changes: update `.ai/architecture.md` and `.ai/decisions.md`.
- Rule engine or data model changes: update `.ai/data-and-rules.md`.
- Coding/test workflow changes: update `.ai/engineering-rules.md`.
- Unresolved decisions: update `.ai/open-questions.md`.
