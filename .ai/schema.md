# Schema Context

This file defines the MVP logical database schema.

The MVP goal is:

> Given a selected target card and an effect code, return source cards whose active stored selector matches that target card.

The schema is optimized for:

- local card catalog storage;
- external AI extraction writes;
- per-card AI processing metadata;
- reprocessing by extraction version;
- real-time selector matching;
- simple public queries.

The MVP does not store precomputed selector-to-target relationships.

## Schema Principles

- Store all imported cards in a normalized catalog.
- Store the latest raw upstream card payload.
- Store external AI extraction output separately from card data.
- Track completed AI processing on each card even when no effect is found.
- Use JSONB for selector criteria and raw upstream card payloads.
- Keep controlled vocabulary values as text validated in application code.
- Use migrations as the executable source of truth once implemented.

## `cards`

Stores normalized Yu-Gi-Oh! card data imported from an external source.

| Field | Type | Description |
| --- | --- | --- |
| `id` | uuid primary key | Internal card ID. |
| `upstream_source` | text not null | Source system identifier. |
| `upstream_id` | text not null | Upstream card identifier. |
| `passcode` | text null | Official passcode when available. |
| `konami_id` | text null | Official Konami database ID when available. |
| `name` | text not null | Official English card name. |
| `normalized_name` | text not null | Search-normalized card name. |
| `aliases` | jsonb not null default `[]` | Structured official alternate or treated-as names. |
| `normalized_aliases` | text[] not null default `{}` | Normalized exact-name aliases used for lookup and selector matching. |
| `description` | text not null | Official English card text. |
| `card_type` | text null | Monster, Spell, Trap, etc. |
| `frame_type` | text null | Normal, Effect, Fusion, Synchro, etc. |
| `race` | text null | Warrior, Spellcaster, Quick-Play, etc. |
| `attribute` | text null | DARK, LIGHT, etc. |
| `monster_categories` | text[] not null default `{}` | Tuner, Spirit, Toon, etc. |
| `spell_trap_type` | text null | Normal, Continuous, Quick-Play, etc. |
| `atk` | integer null | ATK value. |
| `def` | integer null | DEF value. |
| `level` | integer null | Monster Level. |
| `rank` | integer null | Xyz Rank. |
| `link_rating` | integer null | Link Rating. |
| `archetype` | text null | Archetype/family name when available. |
| `mentions` | text[] not null default `{}` | Official card names or explicit family/archetype labels specifically listed in this card's text. |
| `text_features` | text[] not null default `{}` | Normalized feature flags derived from official card text, such as `coin_toss_effect`, `places_counter_effect`, `dice_roll_effect`, or `cannot_be_normal_summoned_set`. |
| `image_url` | text null | Upstream image URL used by the frontend. |
| `ai_processing` | jsonb not null default `{}` | Latest completed AI processing metadata by effect code. |
| `raw_payload` | jsonb not null | Latest raw upstream payload. |
| `imported_at` | timestamptz not null | First import timestamp. |
| `updated_at` | timestamptz not null | Last update timestamp. |

Indexes and constraints:

- unique `upstream_source, upstream_id`;
- trigram index on `normalized_name`;
- GIN index on `normalized_aliases`;
- GIN index on `monster_categories`;
- GIN index on `mentions`;
- GIN index on `text_features`;
- indexes on commonly matched selector fields such as `card_type`, `race`, `attribute`, `level`, `rank`, `link_rating`, `atk`, `def`, and `archetype`.

## Card Alias Objects

`cards.aliases` stores official alternate or treated-as values on the card row.

Alias object fields:

| Field | Type | Description |
| --- | --- | --- |
| `alias` | text not null | Official English alias or treated-as value. |
| `normalized_alias` | text not null | Search-normalized alias value. |
| `alias_kind` | text not null | Meaning of this alias for matching. |
| `applies_in_zone_codes` | text[] not null default `[]` | Zones where the alias applies. Empty means every zone. |
| `condition_json` | jsonb not null default `{}` | Structured condition for applying the alias. |
| `source` | text not null | Origin of the alias. |

Supported MVP values:

- `alias_kind`: `exact_name`, `archetype_membership`.
- `source`: `card_text`, `upstream`.
- `applies_in_zone_codes[]`: `deck`, `hand`, `field`, `gy`, `banishment`.
- `condition_json.type`: `always`.

`exact_name` means the card is treated as another exact card name. Example: `Harpie Lady 1` is always treated as `Harpie Lady`.

`archetype_membership` means the card is treated as a member of a named card family or archetype without becoming that exact card name. Example: a card whose text says it is always treated as an `Archfiend` card.

Example:

```json
{
  "alias": "Harpie Lady",
  "normalized_alias": "harpie lady",
  "alias_kind": "exact_name",
  "applies_in_zone_codes": [],
  "condition_json": {
    "type": "always"
  },
  "source": "card_text"
}
```

Example for archetype membership:

```json
{
  "alias": "Archfiend",
  "normalized_alias": "archfiend",
  "alias_kind": "archetype_membership",
  "applies_in_zone_codes": [],
  "condition_json": {
    "type": "always"
  },
  "source": "card_text"
}
```

Matching notes:

- `normalized_aliases` is generated from `aliases`; do not edit it independently.
- Exact-name selectors should match `name` and `exact_name` aliases that apply in the relevant zone.
- Archetype/family selectors should match native archetype data and `archetype_membership` aliases that apply in the relevant zone.
- An always-applied alias with empty `applies_in_zone_codes` applies in `deck`.
- Conditional aliases that do not apply in `deck` must not affect `add_deck_to_hand` matching.

## Effect Codes

Supported effect filters are application constants, not database rows.

Initial supported value:

- `add_deck_to_hand`

Use `effect_code` text columns in extraction tables. Validate supported values in application code.

## `cards.ai_processing`

Tracks latest completed AI processing per card.

Extraction versions use this format:

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

- missing key means the card has not completed that process;
- matching `version` means the card was processed by the current extractor version;
- different `version` means the card should be reprocessed;
- `result_count = 0` means processing completed and no supported effect was extracted.

The MVP does not store extraction run history.

## `card_effects`

Stores supported effects extracted from source cards.

| Field | Type | Description |
| --- | --- | --- |
| `id` | uuid primary key | Internal extracted effect ID. |
| `source_card_id` | uuid not null references `cards(id)` | Card containing the effect. |
| `effect_code` | text not null | Supported effect code. |
| `extraction_version` | text not null | Extraction version that produced this effect. |
| `source_text` | text not null | Full official card text segment that caused extraction. |
| `condition_text` | text null | Exact condition or timing text such as `If ...`, `When ...`, or `During ...`, without trailing punctuation. |
| `cost_text` | text null | Exact cost text such as `Discard 1 card`, without trailing punctuation. |
| `action_text` | text not null | Isolated supported action text, such as `Add ... from your Deck to your hand`. |
| `restriction_text` | text null | Exact restriction or limitation text tied to the extracted effect. |
| `selector_status` | text not null | Whether the extracted effect is public-matchable, intentionally ignored, or unresolved. |
| `selector_json` | jsonb null | Canonical selector criteria when resolved. |
| `selector_hash` | text not null | Stable hash of canonical selector. |
| `ai_confidence` | numeric null | AI confidence score when provided. |
| `is_active` | boolean not null default true | Whether this extracted effect is used in public results. |
| `created_at` | timestamptz not null | Creation timestamp. |
| `updated_at` | timestamptz not null | Last update timestamp. |

Supported `selector_status` values:

- `resolved`;
- `ignored`;
- `unresolved`.

Meanings:

- `resolved`: selector is deterministic and should be used by public matching.
- `ignored`: parser recognized the effect and the handling rule is known, but it should not be shown in public matching, such as selectors that would match any card, any Monster, any Spell, or any Trap.
- `unresolved`: parser recognized a supported action but does not yet know how to represent it safely; these records need later review.

Indexes and constraints:

- index on `source_card_id`;
- index on `effect_code, is_active, selector_status`;
- index on `selector_hash`;
- optional unique index on `source_card_id, effect_code, extraction_version, selector_hash, source_text`.

Text segment rules:

- `source_text` is the full source segment that caused extraction.
- `condition_text` records `if`, `when`, `during`, or similar timing/condition text.
- `cost_text` records explicit costs such as `discard`, `pay`, `tribute`, `banish`, `send`, or `destroy`.
- `action_text` is the isolated supported action.
- `restriction_text` records post-action restrictions or locks, such as `also, for the rest of this turn...`.
- Selector exclusions such as `except "Card Name"` belong in `action_text` and `selector_json`, not `restriction_text`.
- These fields preserve official wording and must not be summaries or paraphrases.
- If a segment cannot be split confidently, keep `source_text` exact and leave uncertain split fields null.

Example:

```json
{
  "effect_code": "add_deck_to_hand",
  "source_text": "Add 1 Level 4 or lower Warrior monster from your Deck to your hand.",
  "condition_text": null,
  "cost_text": null,
  "action_text": "Add 1 Level 4 or lower Warrior monster from your Deck to your hand.",
  "restriction_text": null,
  "selector_status": "resolved",
  "selector_json": {
    "type": "logical",
    "op": "and",
    "args": [
      { "type": "comparison", "field": "card_type", "op": "=", "value": "Monster" },
      { "type": "comparison", "field": "race", "op": "=", "value": "Warrior" },
      { "type": "comparison", "field": "level", "op": "<=", "value": 4 }
    ]
  },
  "is_active": true
}
```

Example with mixed alternatives:

```json
{
  "effect_code": "add_deck_to_hand",
  "source_text": "When this card is activated: Look at the top 3 cards of your Deck, then you can reveal 1 \"Dark Magician\" or 1 Spell/Trap that specifically lists the card \"Dark Magician\" in its text, among them, and add it to your hand, also place the remaining cards on top of your Deck in any order.",
  "condition_text": "When this card is activated",
  "cost_text": null,
  "action_text": "Look at the top 3 cards of your Deck, then you can reveal 1 \"Dark Magician\" or 1 Spell/Trap that specifically lists the card \"Dark Magician\" in its text, among them, and add it to your hand, also place the remaining cards on top of your Deck in any order.",
  "restriction_text": null,
  "selector_status": "resolved",
  "selector_json": {
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
  },
  "is_active": true
}
```

Selector expression contract:

- `comparison` nodes require `field`, `op`, and `value`.
- `logical` nodes require `op` and `args`.
- supported logical operators: `and`, `or`;
- supported comparison operators: `=`, `!=`, `in`, `not_in`, `contains`, `not_contains`, `<`, `<=`, `>`, `>=`;
- unknown fields, unknown operators, invalid value types, or malformed expression nodes must fail closed.

Supported selector fields include `monster_categories`. Use `contains` or `not_contains` for category filters such as `Ritual`, `Pendulum`, `Tuner`, `Union`, `Spirit`, `Toon`, `Normal`, and similar normalized monster category values.

Supported selector fields include `mentions`. Use `contains` or `not_contains` for exact card-name mentions and explicit family/archetype mention requirements, such as `Dark Magician`, `Destiny HERO`, or `Fallen of Albaz`.

Supported selector fields include `text_features`. Use `contains` or `not_contains` for normalized card text feature flags such as `coin_toss_effect`, `places_counter_effect`, `dice_roll_effect`, or `cannot_be_normal_summoned_set`.

Supported derived selector fields include `combined_atk_def`. It is computed from `atk + def` by matcher/parser code, not stored as a `cards` column unless performance later requires it. Use numeric comparison operators for requirements such as `combined ATK & DEF equal 2000`.

When AI detects a supported effect but cannot define an exact supported selector, store the effect with:

```json
{
  "selector_status": "unresolved",
  "selector_json": null,
  "selector_hash": "unresolved"
}
```

Unresolved effects are audit/debug records only and must not appear in public search results.

When the parser can define the handling rule but the effect should not be exposed, store:

```json
{
  "selector_status": "ignored",
  "selector_json": null,
  "selector_hash": "ignored"
}
```

Ignored effects are decided cases, not parser failures.

## Public Query Shape

The public search query should:

1. load target card by ID;
2. load active `card_effects` where `effect_code = 'add_deck_to_hand'` and `selector_status = 'resolved'`;
3. evaluate `selector_json` against the target card in application rule code;
4. return matching source cards with the matched `source_text` and `action_text`.

If performance later requires precomputed relationships, add an ADR and schema migration.

## MVP Scope Notes

The MVP intentionally does not model:

- full PSCT parsing;
- action phases;
- chain timing;
- activation legality;
- full game simulation;
- complete effect decomposition;
- user reports;
- precomputed target relationships.
