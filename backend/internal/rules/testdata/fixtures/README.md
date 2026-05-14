# Rule Parser Fixtures

Put manual parser examples in this directory.

M1 uses one JSON file with many cases:

```txt
add_from_deck_to_hand.json
```

Each case should contain only:

- `name`: stable case name for test output;
- `effect`: full card effect text;
- `expected_selectors`: selectors expected for supported `add` from Deck to hand actions.

Fixtures should omit database-only or generated fields:

- IDs;
- timestamps;
- versions;
- hashes.

Fixtures should avoid duplicated human-only metadata. M1 assumes every non-empty selector result belongs to an `add` from Deck to hand extracted effect.

The test runner should load every case in `add_from_deck_to_hand.json` and compare normalized parser selectors against `expected_selectors`.
