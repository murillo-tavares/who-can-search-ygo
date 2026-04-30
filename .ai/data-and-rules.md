# Data And Rule Engine Context

## External Data Source

The initial ingestion source is still unresolved.

YGOPRODeck API v7 is a candidate source to evaluate.

Guidelines:

- If using YGOPRODeck, use `cardinfo.php` for bulk card data.
- Store imported data locally.
- Respect documented source rate limits.
- Hotlink card images initially to simplify the MVP.
- Revisit image downloading, caching, or re-hosting in a later version if performance, reliability, or source terms require it.
- Treat the official Konami Yu-Gi-Oh! card database as canonical when wording conflicts are found.

Initial ingestion assumption:

- Use the available cards returned by the selected API/dataset once the source is resolved.
- The selected source or local mock dataset defines the card universe for the MVP.
- Do not add special handling for TCG versus OCG, Master Duel, upcoming cards, Rush Duel, Speed Duel, Speed Spell variants, skill cards, or other variants until scope changes require it.
- Store format availability when available so ingestion scope can expand later.

## Card Data

Store enough data to evaluate whether a target card satisfies search criteria:

- internal ID;
- upstream passcode/API ID;
- Konami ID when available;
- name;
- normalized name;
- card type;
- frame type;
- official English card text;
- race/monster type;
- attribute;
- ATK;
- DEF;
- level or rank;
- pendulum scale;
- link rating;
- link markers;
- archetype;
- card images and metadata;
- card sets and printings;
- banlist information when available;
- format availability when available;
- release dates when available;
- alias or "treated as" data when available;
- normalized alias or "treated as" metadata when rules require it;
- whether the card has an effect when available;
- raw source payload;
- import timestamps.

All imported cards should remain in the database for idempotent sync and future reevaluation. Relationship preprocessing should start from extracted effect rows that match the supported action/source/destination scope.

Keep the latest raw upstream payload for audit/debugging unless a future requirement needs historical snapshots. Query performance should come from normalized fields, indexes, and precomputed relationship tables, not from querying raw payload history.

## Extracted Effects

Persist extracted effects separately from target relationships. A single card can have zero, one, or many extracted effects.

The parser should read card text and register effects supported by the current application rules. When supported rule coverage expands, rerun extraction for affected cards to create newly supported effect records.

Each effect should record:

- searcher card ID;
- action, such as `add`, `send`, `banish`, `special_summon`, or `destroy`;
- source zone, such as `deck`, `gy`, `banishment`, or `extra_deck`;
- destination zone, such as `hand`, `field`, `gy`, or `banishment`;
- search type, such as `exact_name`, `criteria`, `archetype`, or `manual`;
- original card text fragment that produced the extracted effect;
- original text location metadata when practical, such as start and end offsets within the card text;
- normalized condition JSON;
- parser rule ID or pattern ID;
- parser/rule version;
- status, such as `accepted`, `needs_review`, or `invalid`;
- reason when review or invalid status applies;
- processed timestamp.

For the MVP, only effects with action `add`, source zone `deck`, and destination zone `hand` can produce public accepted search relationships. Other recognized action/source/destination flows may be stored for audit or future filters, but they should not generate public relationships until that scope is explicitly added.

Condition JSON should be a small rule expression, not arbitrary raw text. Prefer predictable operators and fields, for example:

```json
{
  "all": [
    { "field": "name", "op": "eq", "value": "Dark Magician" }
  ]
}
```

```json
{
  "all": [
    { "field": "card_category", "op": "eq", "value": "monster" },
    { "field": "atk", "op": "lte", "value": 1500 }
  ]
}
```

Core rule logic should use typed structs. JSON is the persistence and audit format for extracted criteria, not a replacement for typed parser and matcher code. The original text fragment is for debugging, review, and rule-quality audits; matching should use structured fields and normalized conditions.

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

The first supported extraction should match against the resolving effect segment, not against activation condition or cost text. For example, if the resolving segment contains `Add ... from your Deck to your hand`, create an extracted effect with action `add`, source `deck`, and destination `hand`.

Use periods (`.`), line breaks, and bullet markers as candidate effect-block boundaries, but treat this as a heuristic. A sentence may be a continuation or restriction for the previous effect, especially when it uses phrases such as `this effect`, `that target`, or conditions that must remain true through resolution. Do not assume every period creates an independent effect.

When scanning candidate blocks:

- parse each candidate block using the PSCT punctuation rules above;
- extract supported effects from the resolving effect segment;
- skip candidate blocks whose resolving effect segment does not match a supported action/source/destination pattern;
- preserve the original text fragment for every extracted effect.

## Search Relationships

Persist precomputed relationships derived from supported extracted effects with:

- searcher card ID;
- target card ID;
- extracted effect ID;
- action;
- source zone;
- destination zone;
- match kind, such as `exact_name`, `criteria`, or `manual`;
- matched criteria JSON;
- status, such as `accepted`, `needs_review`, or `invalid`;
- preprocessing version;
- processed timestamp.

Only `accepted` relationships should be shown publicly.

Public searcher queries should filter relationships by status, action, source zone, and destination zone. The default public query uses status `accepted`, action `add`, source `deck`, and destination `hand`. Request handling should not inspect card text.

Relationship preprocessing should consider only extracted effects matching the relationship scope being generated, such as action `add`, source `deck`, and destination `hand` for the initial public search. It should not discover candidate searchers by scanning raw card text.

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

Represent extracted search effects as typed structs. Do not pass loosely shaped maps through core rule logic unless the data is genuinely dynamic.

## Initial Rule Support

Support card text that can add from Deck to hand through:

- exact quoted card names;
- monster type/race criteria;
- Attribute criteria;
- ATK and/or DEF thresholds;
- Level criteria;
- Rank criteria;
- Link Rating criteria;
- card category criteria such as Monster, Spell, or Trap;
- archetype/family criteria when source data supports it reliably.
- alias and "treated as" relationships when they affect whether a target card satisfies a search rule.

Important early examples:

- `Illusion Magic` searching `Dark Magician`.
- `Sangan` searching monsters with `1500 or less ATK`.

These examples should guide early tests, but they do not have to be mandatory fixtures with that exact file structure. Prefer a small curated local dataset or mock source payloads while the importer and rule engine are still being built. Avoid using the full real API dataset for routine development tests until the implementation is stable enough to make failures easy to understand.

## Parser Boundaries

Do not attempt full natural-language understanding in the MVP.

When a card text contains unsupported wording, ambiguous scope, or criteria the engine cannot safely model:

- store a supported extracted effect as `needs_review` when the parser can identify the effect family but cannot safely accept the normalized criteria; or
- create no extracted effect when the wording is outside currently supported rules.

Do not silently create accepted relationships from ambiguous text.
