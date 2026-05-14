# Rule Fixture Contract

Rule parser fixtures define examples that a human can review before implementation.

Each fixture provides one full effect text and the selectors expected from that text.

Canonical fixture location:

```txt
backend/internal/rules/testdata/fixtures/
```

Use one JSON file with a `cases` array for the M1 fixture list:

```txt
backend/internal/rules/testdata/fixtures/add_from_deck_to_hand.json
```

This keeps a large manually curated list easy to scan and edit in one place.

M1 fixture tests are intentionally scoped to `add` from Deck to hand. The action is implicit; fixtures only assert which selectors should be extracted for that action.

## Goals

- Make rule behavior explicit before database or API work.
- Let AI agents add parser support one fixture at a time.
- Keep unsupported text invisible to public results by expecting empty output.

## Fixture Input

Each case should contain:

- `name`: stable case name for test output and quick identification.
- `effect`: full effect text used by the parser.
- `expected_selectors`: selectors for supported `add` from Deck to hand actions.

Optional fixture metadata:

- `source`: source or card reference used to create the fixture.

## Expected Output Rules

- If the parser fully supports the relevant effect, `expected_selectors` contains the selectors for that effect.
- If the parser does not support the text, `expected_selectors` is an empty array.
- Do not model partial parser output.
- Do not use status fields in fixtures.
- Keep expected output focused on selector criteria, not the database schema.
- Omit database-only or generated fields such as IDs, timestamps, versions, and hashes.

## Suggested Shape

```json
{
  "cases": [
    {
      "name": "reinforcement_of_the_army_level_race",
      "effect": "Add 1 Level 4 or lower Warrior monster from your Deck to your hand.",
      "expected_selectors": [
        {
          "criteria": [
            { "field": "card_type", "op": "eq", "value": "Monster" },
            { "field": "race", "op": "eq", "value": "Warrior" },
            { "field": "level", "op": "lte", "value": 4 }
          ]
        }
      ]
    }
  ]
}
```

## Comparison Rules

- Compare normalized parser selectors against `expected_selectors`.
- Preserve selector order.
- Use `name` as the subtest name.
- Do not require database IDs.
- Do not require timestamps.
- Do not require parser or selector versions.
- Do not require hashes in manual fixtures.
- Hashes should be tested separately from normalized parser output.
