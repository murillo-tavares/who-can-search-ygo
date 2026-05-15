# AGENTS.md

This is the entry point for AI agents working on Who Can Search YGO.

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

Who Can Search YGO helps users find Yu-Gi-Oh! cards whose own effects can affect a selected target card.

The initial MVP answers:

```txt
Which cards can add this selected card from the Deck to the hand?
```

The only initially supported effect code is:

```txt
Add ... from your Deck to your hand
```

The application stores cards and AI-extracted effect selectors in its own database. Public search queries use those stored selectors in real time.

## Current Defaults

- Backend: Go.
- Frontend: React, Vite, TypeScript.
- Database: PostgreSQL with `pg_trgm`.
- Card data source: unresolved.
- Card images: hotlinked from the selected card data source.
- AI extraction: use extraction data available in the database.
- MVP public effect filter: `add_deck_to_hand`.

## Non-Negotiables

- Keep backend and frontend as separate subprojects.
- Keep card import idempotent.
- Keep stored AI extraction auditable and versioned.
- Track which AI-supported processes each card has completed, including cards with no extracted result.
- Normal public requests must not call AI.
- Normal public requests must not scan every card text.
- Normal public requests may evaluate stored selectors against the selected target card in real time.
- Use English for documentation, code, identifiers, database names, API fields, comments, commits, UI text, and internal terminology.

## Where To Change Context

- Product scope changes: update `.ai/product.md`.
- Stack or architecture changes: update `.ai/architecture.md` and `.ai/decisions.md`.
- Rule, selector, or data model changes: update `.ai/data-and-rules.md` and `.ai/schema.md`.
- Coding/test workflow changes: update `.ai/engineering-rules.md`.
- Unresolved decisions: update `.ai/open-questions.md`.
