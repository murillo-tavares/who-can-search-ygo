# Data And Rule Engine Context

## External Data Source

The initial ingestion source is still unresolved.

YGOPRODeck API v7 is a candidate source to evaluate.

Guidelines:

- If using YGOPRODeck, use `cardinfo.php` for bulk card data.
- Store imported data locally.
- Respect documented source rate limits.
- Hotlink card images initially to simplify the MVP.

Initial ingestion assumption:

- Use the available cards returned by the selected API/dataset once the source is resolved.
- The selected source or local fixture dataset defines the MVP card universe.
- Do not add special variant handling until scope changes require it.

## Card Data

Store enough data to evaluate whether a target card satisfies search criteria:

- internal ID;
- upstream source and ID;
- passcode;
- Konami ID when available;
- name;
- normalized name;
- aliases and normalized exact-name aliases;
- card type;
- frame type;
- official English card text;
- race/monster type;
- monster categories;
- Spell/Trap type;
- attribute;
- ATK;
- DEF;
- level;
- rank;
- link rating;
- archetype;
- mentioned official card names;
- raw source payload;
- import timestamps.

All imported cards should remain in the database for idempotent sync and reevaluation. Relationship preprocessing should use persisted `card_selectors`, not raw card text.

Keep the latest raw upstream payload for audit/debugging. Query performance should come from normalized fields, indexes, and precomputed relationships.

## Extracted Effects

Persist extracted effects separately from target relationships. A single card can have zero, one, or many extracted effects.

The parser should read card text and register effects supported by the current application rules. When supported rule coverage expands, rerun extraction for affected cards to create newly supported effect records.

Each extracted effect should record:

- searcher card ID;
- parser version;
- effect type;
- action tags;
- PSCT-derived text segments and offsets;
- ordered parsed actions in `actions_json`;
- effect hash;
- extraction timestamp.

For the MVP, only persisted effects with a resolution action where `action_kind = move_card`, `verb = add`, `config.from = deck`, and `config.to = hand` can produce public searcher results.

Actions should model what the effect actually does. MVP support is `move_card` with:

- `verb = add`;
- `config.from = deck`;
- `config.to = hand`;
- quantity;
- target selectors referencing `card_selectors`.

Quantities should preserve card text semantics:

```json
{ "kind": "exactly", "value": 1 }
```

```json
{ "kind": "up_to", "max": 2, "unit": "copies" }
```

Target criteria must live in reusable `card_selectors`, not inside raw action text. Criteria use predictable fields and operators from `.ai/schema.md`.

Core rule logic should use typed structs. JSONB is persistence/audit format, not the core rule model.

## Card Text Segmentation

Use official Problem-Solving Card Text punctuation as the first parser guide.

References:

- Konami, "Problem-Solving Card Text, Part 3: Conditions, Activations, and Effects": https://www.yugioh-card.com/en/play/psct/psct-3/
- Konami, "Problem-Solving Card Text, Part 4: The Clues on Your Cards": https://www.yugioh-card.com/en/play/psct/psct-4/

For PSCT-style effects:

- Text before `:` is activation condition, timing, or usage condition.
- Text after `:` and before `;` is activation text, such as cost and targeting.
- Text after `;` is the resolving effect.
- If there is `:` but no `;`, text after `:` is the resolving effect.
- If there is `;` but no `:`, text before `;` is activation text and text after `;` is the resolving effect.

The first supported extraction should match against the resolving effect segment, not activation condition or cost text. If the resolving segment contains a fully supported `Add ... from your Deck to your hand` pattern, create an extracted effect with action tag `add`, a `move_card` resolution action, quantity, and target selector references.

Use periods (`.`), line breaks, and bullet markers as candidate effect-block boundaries, but treat this as a heuristic. A sentence may be a continuation or restriction for the previous effect, especially when it uses phrases such as `this effect`, `that target`, or conditions that must remain true through resolution. Do not assume every period creates an independent effect.

When scanning candidate blocks:

- parse each candidate block using the PSCT punctuation rules above;
- extract supported effects from the resolving effect segment into typed structs;
- skip candidate blocks whose resolving effect segment does not match a supported parser pattern;

## Search Relationships

Persist selector-target relationships with:

- card selector ID;
- target card ID;
- relationship version;
- match kind;
- relationship hash;
- processed timestamp.

Persisted relationships are valid for public searcher results.

Public searcher queries should filter relationships by target card, then join extracted effects whose `actions_json` references the relationship selector and matches the MVP action scope. Request handling should not inspect card text.

Relationship preprocessing should build relationships from persisted `card_selectors`. It should not discover candidate searchers by scanning raw card text.

## Rule Engine Requirements

The rule engine must be:

- isolated from HTTP and database code;
- deterministic;
- testable with fixtures;
- versioned;
- conservative when uncertain.

Implementation location:

```txt
backend/internal/rules
```

Represent extracted effects as typed structs. Do not pass loosely shaped maps through core rule logic unless the data is genuinely dynamic.

## Initial Rule Support

Support card text that can add from Deck to hand through:

- exact quoted card names;
- monster type/race criteria;
- Attribute criteria;
- ATK and/or DEF thresholds;
- Level criteria;
- Rank criteria;
- Link Rating criteria;
- card type criteria such as Monster, Spell, or Trap;
- archetype/family criteria when source data supports it reliably;
- alias and "treated as" relationships when they affect whether a target card satisfies a search rule.

Important early examples:

- `Illusion Magic` searching `Dark Magician`.
- `Sangan` searching monsters with `1500 or less ATK`.

Use these examples to guide early parser and matcher tests. Prefer curated fixtures or mock source payloads until full-dataset failures are easy to debug.

## Parser Boundaries

Do not attempt full natural-language understanding in the MVP.

When a card text contains unsupported wording, ambiguous scope, or criteria the engine cannot fully and safely model:

- create no extracted effect;
- create no selector; and
- create no relationship.

The MVP parser should not persist partial extracted effects. If it cannot parse 100% of the supported action and selector semantics, it should ignore that candidate block.
