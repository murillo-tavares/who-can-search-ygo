# M1: Backend Rule Parser Foundation

Status: planned

## Goal

Start the backend with a small, deterministic rule parser foundation driven by manually curated card-text fixtures before database, API, or frontend work begins.

## Why This Comes First

The project's core value depends on correctly extracting supported `Add ... from your Deck to your hand` effects. The first milestone should let a human provide full effect text and expected selectors, then let tests prove the application generates the same selectors.

## Deliverables

- Backend Go module under `backend/`.
- `backend/internal/rules` package.
- Typed rule model for parser input and output:
  - card input;
  - extracted effect;
  - text segments;
  - action;
  - quantity;
  - card selector;
- Fixture file contract for rule parser examples.
- Manual rule fixture file where each case contains:
  - stable case name;
  - full effect text in `effect`;
  - expected selectors for supported `add` from Deck to hand actions.
- Fixture file at `backend/internal/rules/testdata/fixtures/add_from_deck_to_hand.json`, with a `cases` array for the full manual test list.
- Fixture-driven table tests that run the parser for every case and compare generated selectors against `expected_selectors`.
- Name and text normalization helpers needed by the first fixtures.
- PSCT segmentation helper for condition, activation, and resolution text.
- First supported exact-name parser case for `Add ... from your Deck to your hand`.
- At least one unsupported fixture that produces no selector.

## Out Of Scope

- PostgreSQL migrations.
- `sqlc` query generation.
- HTTP API endpoints.
- Importing the full external dataset.
- Relationship preprocessing.
- Frontend work.
- Full natural-language parsing.
- Competitive legality, costs, restrictions, rulings, or combo discovery.

## Initial Fixture Cases

| Fixture | Purpose | Expected Result |
| --- | --- | --- |
| `illusion_magic_exact_name` | Exact quoted card name target. | Extracted effect with exact-name selector for `Dark Magician`. |
| `reinforcement_of_the_army_level_race` | Criteria selector with Level and race. | Extracted effect with Monster, Warrior, Level <= 4 selector once criteria parsing starts. |
| `sangan_atk_threshold` | Criteria selector with ATK threshold. | Extracted effect with Monster and ATK <= 1500 selector once criteria parsing starts. |
| `beckoning_light` | Add from GY to hand, outside Deck-to-hand scope. | Empty selector list. |

## Done When

- `go test ./...` passes inside `backend/`.
- Fixture tests load all manual examples and fail when generated selectors differ from `expected_selectors`.
- Fixture tests cover at least one supported exact-name add case and one unsupported case with an empty selector list.
- Parser output uses typed structs, not ad hoc maps.
- Parser only accepts effects whose resolution action is `move_card`, `add`, from `deck` to `hand`.
- Documentation explains how to add a new parser fixture with full effect text and expected selectors.

## Implementation Notes

- Keep rule parsing isolated in `backend/internal/rules`.
- Keep fixture tests independent from PostgreSQL and network access.
- Use `rules-1` as the first parser version.
- Do not persist JSON in M1. JSON fixture files may be used for tests, but core parser logic should operate on typed Go structs.
- Prefer conservative behavior: uncertain text should create no extracted effect, no selector, and no relationship.
- Generated output should be normalized before comparison so fixture diffs reflect rule behavior, not incidental ordering or formatting.

## Open Questions For This Milestone

- No open questions right now. Add more cases to `add_from_deck_to_hand.json` as parser coverage expands.
