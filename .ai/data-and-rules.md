# Data And Rule Context

## Product Direction

The MVP focuses on:

- an imported Yu-Gi-Oh! card catalog;
- supported effect codes;
- AI-extracted selectors;
- per-card AI processing metadata by effect code and extraction version;
- deterministic selector matching at request time.

The initial supported effect code is:

- `add_deck_to_hand`

Generic effect parsing guidance and the initial `add_deck_to_hand` profile live in `.ai/effect-parsing.md`.

Future effect codes may include:

- add from GY to hand;
- send from Deck to GY;
- special summon from Deck;
- other explicitly supported categories.

Future effect codes require explicit scope changes.

## External Card Data Source

The initial card data source is unresolved.

Candidate:

- YGOPRODeck API v7.

Guidelines:

- Import cards into the local database.
- Keep normalized fields needed for selector matching.
- Keep the latest raw upstream payload for audit/debugging.
- Use upstream image URLs directly.
- Do not store images locally in the MVP.
- Keep source-specific mapping isolated.

## Card Data

Store enough normalized card data to evaluate selectors against a target card.

Important fields:

- internal ID;
- upstream source and ID;
- passcode;
- Konami ID, when available;
- name;
- normalized name;
- official aliases or treated-as names, when available;
- official English card text;
- card type;
- frame type;
- race or monster type;
- monster categories;
- Spell/Trap type;
- attribute;
- ATK;
- DEF;
- derived values such as combined ATK and DEF;
- Level;
- Rank;
- Link Rating;
- archetype;
- specifically mentioned official card names;
- normalized card text feature flags;
- upstream image URL;
- raw upstream payload;
- import timestamps.

Imported cards remain available as target cards even when they have no extracted effects.

## Effect Codes

Effect codes define what the application can search for.

Initial effect code:

```txt
add_deck_to_hand
```

For the MVP, the frontend may hardcode `add_deck_to_hand`.

## External AI Extraction Model

AI extraction is external to the MVP application.

The external process should:

1. Read cards from the database.
2. Decide whether a card has a supported effect code.
3. Generate normalized selector criteria for supported effects.
4. Update processing metadata on every processed card.
5. Write resolved effects when selector output is valid and useful for public matching.
6. Optionally write unresolved effects when a supported effect is detected but the selector cannot be represented exactly.
7. Optionally write ignored effects when the parser knows the rule but the effect should not be exposed publicly.

The application must track completed processing even when no effect is extracted.

This supports:

- resuming interrupted extraction;
- identifying cards already processed;
- identifying cards with no supported effect;
- reprocessing using a new extraction version;
- identifying new cards that need extraction.

## Per-Card Processing Metadata

AI processing completion is stored on `cards.ai_processing`.

This metadata records the latest completed processing per effect code.

Use extraction versions in this format:

```txt
<effect-code-kebab>-v<number>
```

Initial version:

```txt
add-deck-to-hand-v1
```

Example:

```json
{
  "add_deck_to_hand": {
    "version": "add-deck-to-hand-v1",
    "processed_at": "2026-05-15T12:00:00Z",
    "result_count": 1
  }
}
```

Rules:

- missing effect code means not processed;
- matching version means processed for the current extractor;
- different version means reprocess is needed;
- `result_count = 0` means the card was checked and no supported effect was extracted.

The application does not need to know how many extraction runs happened.

## Selectors

Selectors describe valid target cards.

Selectors are generated externally by AI, then stored and matched deterministically by the application.

Selectors should use a generic typed expression AST, similar to a safe subset of SQL `where` semantics.

Do not store selectors as raw SQL strings. Store structured JSON so the application can validate fields, operators, values, and boolean nesting before matching.

Example selector:

```json
{
  "type": "logical",
  "op": "and",
  "args": [
    { "type": "comparison", "field": "card_type", "op": "=", "value": "Monster" },
    { "type": "comparison", "field": "race", "op": "=", "value": "Warrior" },
    { "type": "comparison", "field": "level", "op": "<=", "value": 4 }
  ]
}
```

Example with mixed alternatives:

```txt
name = 'Dark Magician' OR (card_type IN ('Spell', 'Trap') AND mentions CONTAINS 'Dark Magician')
```

Canonical JSON:

```json
{
  "type": "logical",
  "op": "or",
  "args": [
    { "type": "comparison", "field": "name", "op": "=", "value": "Dark Magician" },
    {
      "type": "logical",
      "op": "and",
      "args": [
        { "type": "comparison", "field": "card_type", "op": "in", "value": ["Spell", "Trap"] },
        { "type": "comparison", "field": "mentions", "op": "contains", "value": "Dark Magician" }
      ]
    }
  ]
}
```

Initial selector support should remain conservative.

Supported fields:

- `name`;
- `card_type`;
- `race`;
- `attribute`;
- `level`;
- `rank`;
- `link_rating`;
- `atk`;
- `def`;
- `combined_atk_def`;
- `monster_categories`;
- `archetype`;
- `mentions`;
- `text_features`;
- `spell_trap_type`.

Text feature selector matching:

- `text_features contains "<feature_code>"` should match normalized card text feature flags stored on the target card.
- Feature codes describe card text properties that are useful as selector filters, such as cards that have a coin tossing effect or cards that place counters.
- Feature codes should be stable snake_case strings.
- Initial examples include `coin_toss_effect`, `places_counter_effect`, `dice_roll_effect`, and `cannot_be_normal_summoned_set`.

Derived numeric selector matching:

- `combined_atk_def` is a derived selector field computed by matcher/parser code as `atk + def` when both values are known.
- Use it for target requirements such as `whose combined ATK & DEF equal 2000`.
- Do not store it as a physical database column unless query performance later requires it.
- If either `atk` or `def` is unknown, comparisons against `combined_atk_def` should fail closed.

Total-Level set constraints:

- Some effects select multiple monsters whose total Levels equal a number.
- For one-target public matching, reduce `up to N <static filters> whose total Levels equal T` to the static filters plus `level <= T - (N - 1)`.
- This assumes the other selected monsters have minimum Level 1.
- Example: `up to 3 Rock monsters whose total Levels equal 8` becomes `card_type = "Monster"`, `race = "Rock"`, and `level <= 6`.
- If the count or total cannot be parsed, or the reduction produces no valid maximum Level, fail closed to unresolved.

Exact-name selector matching:

- `name = "Card Name"` should match the target card's official name.
- It should also match `exact_name` aliases that apply in the relevant zone.
- For `add_deck_to_hand`, the relevant target zone is `deck`.
- MVP alias matching should focus on `condition_json.type = always`, such as cards whose text says they are always treated as another card name.
- Conditional aliases that do not apply in `deck` should be ignored for `add_deck_to_hand`.

Archetype/family selector matching:

- `archetype = "Family Name"` should match native archetype data when available.
- It should also match `archetype_membership` aliases that apply in the relevant zone.
- This covers card text such as `This card is always treated as an "Archfiend" card.`
- Archetype membership must not make the card match an exact `name = "Archfiend"` selector.
- If the selected card data source does not provide an archetype for a card, treat that card as not belonging to the archetype.

Mention selector matching:

- `mentions contains "<value>"` may represent exact card-name mentions or explicit family/archetype mention requirements from card text.
- Use exact names when the effect names one card, such as `mentions "Dark Magician"`.
- Use family values when the effect names a family or archetype, such as `mentions a "Destiny HERO" monster's card name`.
- Parser and import code should keep these values normalized enough for deterministic matching, but should preserve official English names and family labels as selector values.

Supported operators:

- `=`;
- `!=`;
- `in`;
- `not_in`;
- `contains`;
- `not_contains`;
- `<`;
- `<=`;
- `>`;
- `>=`.

Supported expression nodes:

- `comparison`;
- `logical`.

Supported logical operators:

- `and`;
- `or`.

If a selector cannot be represented with supported fields and operators, the external extractor may store the effect with `selector_status = unresolved`, but it must not appear in public results.

Player-choice comparison filters that cannot be statically defined from the effect text may be ignored when building the selector. This includes wording such as `different names from each other`, `different from cards you control`, `same Attribute as the revealed card`, `same name as the banished card`, or `declared card` when the referenced object depends on runtime player choice. Ignore only the dynamic comparison. Preserve any static filters that remain.

Pronoun, declaration, or same-name wording may still be resolved when the same effect provides a static family, type, category, or stat constraint. For example, if a runtime declaration only chooses which `Gunkan` monster is named, the selector can keep `archetype = "Gunkan"` and `card_type = "Monster"`. If an effect sends a `Level 4 or lower Normal Monster` and then adds a card with the same name, the selector can keep `card_type = "Monster"`, `monster_categories contains "Normal"`, and `level <= 4`.

Runtime numeric requirements may be reduced only when there is a clear hard maximum. For example, a requirement based on the number of monsters a player controls can use the normal monster-zone maximum when that produces a conservative bound.

Broad generic selectors:

- Do not expose effects whose selector would match every card, such as `add 1 card from your Deck to your hand`.
- Only selectors that reduce to any card, any Monster, any Spell, or any Trap are considered too generic by default.
- Store decided generic cases as `selector_status = ignored`, not `unresolved`.
- Keep broad selectors when a meaningful filter remains, such as archetype, named card, text feature, monster category, stat, or mention requirement.

Boolean nesting has no fixed depth limit. The matcher should recursively evaluate valid `logical` and `comparison` nodes. Unknown fields, unknown operators, invalid value types, or malformed expression nodes must fail closed.

## Extracted Effects

Each extracted effect stores:

- source card;
- effect code;
- selector;
- selector status;
- extraction version;
- source text, copied exactly from the card text segment that contains the supported action;
- optional condition text, copied exactly from the source segment without trailing punctuation;
- optional cost text, copied exactly from the source segment without trailing punctuation;
- isolated action text for the supported effect code;
- optional restriction text, copied exactly from the source segment when relevant;
- active flag;
- timestamps.

Text segment rules:

- `source_text` is the full source segment that caused extraction.
- `condition_text` records `if`, `when`, `during`, or similar timing/condition text.
- `cost_text` records explicit costs such as `discard`, `pay`, `tribute`, `banish`, `send`, or `destroy`.
- `action_text` is the isolated supported action, such as `Add ... from your Deck to your hand`.
- `restriction_text` records post-action restrictions or locks, such as `also, for the rest of this turn...`.
- Selector exclusions such as `except "Card Name"` belong in `action_text` and `selector_json`, not `restriction_text`.
- These fields should preserve official wording and should not be summaries or paraphrases.
- If a segment cannot be split confidently, keep `source_text` exact and leave uncertain split fields null.

Only active extracted effects with `selector_status = resolved` are used in public results.

The external extractor should write public-matchable effects only when the selector is resolved. Incorrect records may be disabled later by setting `is_active = false`.

## Public Matching

Public matching is deterministic:

1. Load target card.
2. Load active effects for the effect code with `selector_status = resolved`.
3. Evaluate stored selectors against the target card.
4. Return matching source cards.

Public matching must not:

- call AI;
- scan all card text;
- generate new selectors;
- mutate extraction state.

## Conservative Rules

When uncertain:

- update per-card processing metadata;
- store an unresolved effect for audit/debugging, or skip storing the effect entirely;
- do not expose it publicly.

The system should prefer missing relationships over incorrect results.
