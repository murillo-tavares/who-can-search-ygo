# Schema Context

This file defines the initial logical database schema. Migrations are the executable source of truth once implemented, but they should follow this model unless a later decision updates it.

## Schema Principles

- Store imported card data separately from reusable card selectors, extracted effects, and search relationships.
- Parse card text into supported extracted effect records before relationship preprocessing.
- Normal public requests must query precomputed relationships, not card text.
- Keep target-card relationships tied to reusable card selectors. Each selector stores the full criteria set for one valid target group.
- Keep parser versions on extracted effects so results can be audited and regenerated without storing every possible rule combination.
- Store controlled vocabulary values as text and define their allowed enum values in application code.
- Use build-like identifiers for rule versions, such as `rules-1` and `relationships-1`.
- Use JSONB for raw upstream payloads, canonical action configuration, canonical selector criteria, text segments, and audit metadata, not as a replacement for indexed relational fields.

## Core Tables

### `cards`

Description:

Stores normalized card data and the latest raw source payload. Cards can be searchers, targets, both, or neither.

Fields:

| Field | Type | Description |
| --- | --- | --- |
| `id` | uuid primary key | Internal card ID used by application relationships and public API references. |
| `upstream_source` | text not null | Source system identifier, such as `ygoprodeck`, once the initial source is chosen. |
| `upstream_id` | text not null | Card identifier from the upstream source. |
| `passcode` | text null | Yu-Gi-Oh! card passcode when available. |
| `konami_id` | text null | Official Konami database ID when available. |
| `name` | text not null | Official English card name from the selected source. |
| `normalized_name` | text not null | Search-normalized card name for lookup and matching. |
| `aliases` | jsonb not null default `[]` | Structured official alternate or treated-as names for this card. Stores alias metadata directly on the card row. |
| `normalized_aliases` | text[] not null default `{}` | Search-normalized exact-name alias values derived from `aliases` for card-name lookup and exact-name rule matching. |
| `card_type` | text null | Broad card category, such as Monster, Spell, or Trap, when available. |
| `frame_type` | text null | Source frame type, such as normal, effect, spell, trap, synchro, xyz, link, or pendulum variants. |
| `description` | text not null | Official English card text used by the parser. |
| `race` | text null | Monster race/type, such as Spellcaster or Warrior. Null for Spell and Trap cards. |
| `monster_categories` | text[] not null default `{}` | Special monster classifications, such as `tuner`, `gemini`, `spirit`, `union`, `toon`, and `flip`. Empty for cards without special monster categories. |
| `spell_trap_type` | text null | Normalized Spell/Trap type code, such as `normal`, `quick_play`, `continuous`, `equip`, `field`, `ritual`, or `counter`. Null for Monster cards. |
| `attribute` | text null | Monster Attribute, such as DARK, LIGHT, EARTH, WATER, FIRE, WIND, or DIVINE. |
| `atk` | integer null | Monster ATK value when numeric and available. |
| `def` | integer null | Monster DEF value when numeric and available. |
| `level` | integer null | Monster Level when available. |
| `rank` | integer null | Xyz Monster Rank when available. |
| `link_rating` | integer null | Link Monster Link Rating when available. |
| `archetype` | text null | Source-provided archetype or family name when reliable. |
| `mentions` | text[] not null default `{}` | Normalized official card names mentioned in this card's text, used by selectors such as Spell/Trap that mentions `Dark Magician`. |
| `raw_payload` | jsonb not null | Latest raw upstream payload for audit and debugging. |
| `imported_at` | timestamptz not null | First import timestamp for this upstream card. |
| `updated_at` | timestamptz not null | Last import/update timestamp for this card row. |

Supported values:

- `upstream_source`: unresolved until the initial ingestion source is selected; `ygoprodeck` is the current candidate, and `fixture` may be used for curated local test data.
- `card_type`: source-driven broad category values, normalized before rules depend on them; expected MVP values are `Monster`, `Spell`, and `Trap`.
- `frame_type`, `race`, and `attribute`: source-driven values, normalized by importer code before rules depend on them.
- `spell_trap_type`: `normal`, `quick_play`, `continuous`, `equip`, `field`, `ritual`, `counter`.

Constraints and indexes:

- Unique `upstream_source, upstream_id`.
- Trigram index on `normalized_name`.
- Index on `passcode` when available.
- GIN index on `normalized_aliases`.
- Optional GIN index on `aliases` only when review or import tooling needs metadata queries over alias objects.
- GIN index on `mentions`.
- GIN index on `monster_categories`.

### Card Alias Objects

Description:

`cards.aliases` stores structured official alternate or treated-as values on the card row. It is not only for alternate card names. It can also preserve official text saying a card is always treated as part of a named card family, such as an "Archfiend" card.

Alias object fields:

| Field | Type | Description |
| --- | --- | --- |
| `alias` | text not null | Official English alias or treated-as value. |
| `normalized_alias` | text not null | Search-normalized alias value. |
| `alias_kind` | text not null | Meaning of this alias for matching. |
| `applies_in_zone_codes` | text[] not null default `[]` | Zones where the alias applies. Empty array means it applies in every zone. |
| `condition_json` | jsonb not null default `{}` | Structured condition for applying the alias. |
| `source` | text not null | Origin of the alias. |

Supported values:

- `alias_kind`: `exact_name`, `archetype_membership`.
- `source`: `upstream`, `card_text`, `manual`.
- `applies_in_zone_codes[]`: `hand`, `field`, `gy`, `banishment`.
- `condition_json.type`: `always`, `while_in_zones`.

Alias kind meanings:

- `exact_name`: the card is treated as another exact card name. Example: `Harpie Lady 1` is always treated as `Harpie Lady`.
- `archetype_membership`: the card is treated as a member of a named card family or archetype without becoming that exact card name. Example: `Axe of Despair` is always treated as an `Archfiend` card.

Example for `Harpie Lady 1`, whose text says this card's name is always treated as `Harpie Lady`:

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

Example for `Axe of Despair`, whose text says this card is always treated as an `Archfiend` card:

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

Example for `Proto-Cyber Dragon`, whose text says this card's name becomes `Cyber Dragon` while on the field:

```json
{
  "alias": "Cyber Dragon",
  "normalized_alias": "cyber dragon",
  "alias_kind": "exact_name",
  "applies_in_zone_codes": ["field"],
  "condition_json": {
    "type": "while_in_zones",
    "zone_codes": ["field"]
  },
  "source": "card_text"
}
```

Notes:

- `cards.aliases` replaces a separate alias table. Alias data should be imported or derived directly onto the owning card row.
- `normalized_aliases` is a denormalized lookup field generated only from aliases that should participate in card-name lookup, such as `alias_kind = exact_name`; do not hand-edit it independently.
- Always-applied exact-name aliases can participate in Deck-to-hand target matching once alias matching is implemented.
- Archetype-membership aliases can participate in archetype/card-family criteria matching, but should not make a card match an exact-name selector.
- Conditional aliases, such as names that apply only while on the field, should be preserved in `aliases` but ignored by MVP Deck-to-hand matching unless the condition applies to the modeled zone flow.

Example records:

```json
{
  "id": "00000000-0000-4000-8000-000000000001",
  "upstream_source": "ygoprodeck",
  "upstream_id": "46986414",
  "passcode": "46986414",
  "konami_id": null,
  "name": "Dark Magician",
  "normalized_name": "dark magician",
  "aliases": [],
  "normalized_aliases": [],
  "card_type": "Monster",
  "frame_type": "normal",
  "description": "The ultimate wizard in terms of attack and defense.",
  "race": "Spellcaster",
  "monster_categories": [],
  "spell_trap_type": null,
  "attribute": "DARK",
  "atk": 2500,
  "def": 2100,
  "level": 7,
  "rank": null,
  "link_rating": null,
  "archetype": "Dark Magician",
  "mentions": [],
  "raw_payload": {
    "id": 46986414,
    "name": "Dark Magician"
  },
  "imported_at": "2026-04-30T12:00:00Z",
  "updated_at": "2026-04-30T12:00:00Z"
}
```

```json
{
  "id": "00000000-0000-4000-8000-000000000005",
  "upstream_source": "ygoprodeck",
  "upstream_id": "32807846",
  "passcode": "32807846",
  "konami_id": null,
  "name": "Reinforcement of the Army",
  "normalized_name": "reinforcement of the army",
  "aliases": [],
  "normalized_aliases": [],
  "card_type": "Spell",
  "frame_type": "spell",
  "description": "Add 1 Level 4 or lower Warrior monster from your Deck to your hand.",
  "race": null,
  "monster_categories": [],
  "spell_trap_type": "normal",
  "attribute": null,
  "atk": null,
  "def": null,
  "level": null,
  "rank": null,
  "link_rating": null,
  "archetype": null,
  "mentions": [],
  "raw_payload": {
    "id": 32807846,
    "name": "Reinforcement of the Army"
  },
  "imported_at": "2026-04-30T12:00:00Z",
  "updated_at": "2026-04-30T12:00:00Z"
}
```

## Extracted Effects

### `card_selectors`

Description:

Stores reusable target-card criteria. Selectors are deduplicated by stable hash so multiple extracted effects can share one target group.

Only persist selectors that are useful for supported matching or review workflows. The MVP should persist selectors that can produce target-card relationships. Do not store every recognized non-matching fragment unless a later workflow needs it.

Fields:

| Field | Type | Description |
| --- | --- | --- |
| `id` | uuid primary key | Internal reusable selector ID. |
| `selector_version` | text not null | Version of the canonical selector criteria shape, such as `selector-1`. |
| `criteria_json` | jsonb not null | Canonical card criteria list. All criteria in the list must match. |
| `status` | text not null | Selector review/support status. |
| `selector_hash` | text not null | Stable hash generated from selector version and canonical criteria JSON. |
| `created_at` | timestamptz not null | Timestamp when the selector row was first created. |
| `updated_at` | timestamptz not null | Timestamp when the selector row was last updated. |

Supported values:

- `status`: `accepted`, `needs_review`, `invalid`.
- `criteria_json[].field`: `name`, `card_type`, `race`, `spell_trap_type`, `attribute`, `atk`, `def`, `level`, `rank`, `link_rating`, `archetype`, `mentions`.
- `criteria_json[].op`: `eq`, `in`, `lte`, `gte`.

Constraints and indexes:

- Unique `selector_hash`.
- Index on `status`.
- Optional GIN index on `criteria_json` only when preprocessing or review tooling needs it.

Notes:

- The selector hash must ignore source card identity and text offsets. Two cards with the same normalized selector criteria set should reuse the same selector row.
- The hash should be generated from canonical JSON with stable key ordering and normalized enum/value casing.
- The full card-text occurrence lives on `extracted_effects`, not on `card_selectors`.
- Do not model nested boolean composition in the MVP schema. If one action has multiple alternative valid target groups, reference multiple `card_selectors` from the action.

Example record:

```json
{
  "id": "00000000-0000-4000-8000-000000000040",
  "selector_version": "selector-1",
  "criteria_json": [
    { "field": "card_type", "op": "eq", "value": "Monster" },
    { "field": "race", "op": "eq", "value": "Warrior" },
    { "field": "level", "op": "lte", "value": 4 }
  ],
  "status": "accepted",
  "selector_hash": "sha256:selector-1:monster-warrior-level-lte-4",
  "created_at": "2026-04-30T12:00:00Z",
  "updated_at": "2026-04-30T12:00:00Z"
}
```

Example selector for Spell/Trap cards that mention a card name:

```json
{
  "id": "00000000-0000-4000-8000-000000000041",
  "selector_version": "selector-1",
  "criteria_json": [
    { "field": "card_type", "op": "in", "value": ["Spell", "Trap"] },
    { "field": "mentions", "op": "eq", "value": "Dark Magician" }
  ],
  "status": "accepted",
  "selector_hash": "sha256:selector-1:spell-trap-mentions-dark-magician",
  "created_at": "2026-04-30T12:00:00Z",
  "updated_at": "2026-04-30T12:00:00Z"
}
```

### `extracted_effects`

Description:

Stores parsed effect occurrences from card text. A card can have zero, one, or many extracted effects.

This table keeps effect identity, parser metadata, review status, original text segments, and parsed action specifications. Reusable target criteria sets live in `card_selectors`; actions reference selectors by ID and hash inside `actions_json`.

Fields:

| Field | Type | Description |
| --- | --- | --- |
| `id` | uuid primary key | Internal extracted effect ID. |
| `card_id` | uuid not null references `cards(id)` | Card that contains the parsed effect; this is the searcher card. |
| `parser_version` | text not null | Overall parser version used when the effect was extracted, using a build-like identifier such as `rules-1`. |
| `effect_type` | text not null | Parsed effect type, such as `activated_effect`, `trigger_effect`, or `continuous_effect`. |
| `action_tags` | text[] not null | Application-defined action tags found in this effect. |
| `text_segments_json` | jsonb not null | Full effect text, PSCT-derived text segments, offsets, and optional linking metadata. |
| `actions_json` | jsonb not null default `[]` | Ordered parsed action specifications, including phase, action kind, verb, action-specific config, and referenced card selector IDs/hashes. |
| `status` | text not null | Review status for this extracted effect. |
| `status_reason` | text null | Reason when status is `needs_review` or `invalid`. |
| `effect_hash` | text not null | Stable hash generated from parser version plus normalized text segments, action specifications, and selector hashes. |
| `extracted_at` | timestamptz not null | Timestamp when extraction produced this row. |

Supported values:

- `status`: `accepted`, `needs_review`, `invalid`.
- `effect_type`: `activated_effect`, `trigger_effect`, `continuous_effect`.
- `action_tags`: `add`.

Constraints and indexes:

- Unique `card_id, parser_version, effect_hash`.
- Index on `card_id`.
- Index on `parser_version, status`.
- GIN index on `actions_json` to support public searcher queries that match selector references inside accepted extracted effects.
- Optional GIN index on `action_tags` or `text_segments_json` only when preprocessing, public query filters, or review tooling needs it.

Notes:

- Core rule logic should use typed structs and serialize parsed action occurrences to `actions_json`.
- Store full effect text and PSCT-derived text segments inside `text_segments_json`.
- Matching should use linked `card_selectors`, not raw text fields.
- For the MVP, public search relationships should be generated only for selectors referenced by accepted extracted effects whose `actions_json` contains an action with `action_kind = move_card`, `verb = add`, `config.from = deck`, and `config.to = hand`.

### Extracted Effect Actions

Description:

`actions_json` stores ordered parsed action objects for one extracted effect. Action objects preserve what the effect does after text segmentation while keeping reusable target-card criteria in `card_selectors`.

Action object fields:

| Field | Type | Description |
| --- | --- | --- |
| `phase` | text not null | Effect phase that contains this action. |
| `action_index` | integer not null | Zero-based action order inside the extracted effect. |
| `action_kind` | text not null | Structured action family. |
| `verb` | text not null | Card-text verb represented by this action. |
| `config` | jsonb not null default `{}` | Action-specific configuration such as zones and quantity. |
| `selectors` | jsonb not null default `[]` | Linked reusable selectors used by this action. |
| `text_fragment` | jsonb not null | Source text fragment and offsets that produced this action. |

Supported values:

- `phase`: `resolution`.
- `action_kind`: `move_card`.
- `verb`: `add`.
- `config.from`: `deck`.
- `config.to`: `hand`.
- `config.quantity.kind`: `exactly`, `up_to`.
- `selectors[].selector_role`: `target`.

Notes:

- Action specifications are not deduplicated into their own table. The useful reusable unit for target-card relationships is the selector, not the action.
- `config` shape depends on `action_kind`. MVP `move_card` actions use `from`, `to`, and `quantity`.
- `selectors[]` references rows from `card_selectors` by ID and hash.

Example action object:

```json
{
  "phase": "resolution",
  "action_index": 0,
  "action_kind": "move_card",
  "verb": "add",
  "config": {
    "from": "deck",
    "to": "hand",
    "quantity": { "kind": "exactly", "value": 1 }
  },
  "selectors": [
    {
      "card_selector_id": "00000000-0000-4000-8000-000000000040",
      "selector_hash": "sha256:selector-1:monster-warrior-level-lte-4",
      "selector_role": "target",
      "selector_index": 0
    }
  ],
  "text_fragment": {
    "text": "Add 1 Level 4 or lower Warrior monster from your Deck to your hand.",
    "offsets": { "start": 0, "end": 68 }
  }
}
```

Example record:

```json
{
  "id": "00000000-0000-4000-8000-000000000004",
  "card_id": "00000000-0000-4000-8000-000000000005",
  "parser_version": "rules-1",
  "effect_type": "activated_effect",
  "action_tags": ["add"],
  "text_segments_json": {
    "full": "Add 1 Level 4 or lower Warrior monster from your Deck to your hand.",
    "condition": null,
    "activation": null,
    "resolution": "Add 1 Level 4 or lower Warrior monster from your Deck to your hand.",
    "offsets": {
      "full": { "start": 0, "end": 68 },
      "resolution": { "start": 0, "end": 68 }
    }
  },
  "actions_json": [
    {
      "phase": "resolution",
      "action_index": 0,
      "action_kind": "move_card",
      "verb": "add",
      "config": {
        "from": "deck",
        "to": "hand",
        "quantity": { "kind": "exactly", "value": 1 }
      },
      "selectors": [
        {
          "card_selector_id": "00000000-0000-4000-8000-000000000040",
          "selector_hash": "sha256:selector-1:monster-warrior-level-lte-4",
          "selector_role": "target",
          "selector_index": 0
        }
      ],
      "text_fragment": {
        "text": "Add 1 Level 4 or lower Warrior monster from your Deck to your hand.",
        "offsets": { "start": 0, "end": 68 }
      }
    }
  ],
  "status": "accepted",
  "status_reason": null,
  "effect_hash": "sha256:rules-1:rota-text:add-warrior-level-lte-4",
  "extracted_at": "2026-04-30T12:00:00Z"
}
```

## Search Relationships

Description:

Stores precomputed matches from reusable card selectors to target cards. Normal public searcher queries should read this table, then find accepted extracted effects whose `actions_json` references the matched selector and action scope. They should not scan card text or recompute selectors.

`preprocessing_runs` is intentionally omitted from the MVP schema. Relationship generation should be idempotent through `selector_hash`, `relationship_version`, and `relationship_hash`. If later operations need detailed run auditing, add a separate run table without changing the selector-target relationship shape.

Fields:

| Field | Type | Description |
| --- | --- | --- |
| `id` | uuid primary key | Internal relationship ID. |
| `card_selector_id` | uuid not null references `card_selectors(id)` | Reusable card selector that matched the target card. |
| `target_card_id` | uuid not null references `cards(id)` | Card that satisfies the selector. |
| `relationship_version` | text not null | Relationship matcher/preprocessor version used to produce this row, such as `relationships-1`. |
| `match_kind` | text not null | Matching strategy used for this target. |
| `status` | text not null | Review status for public visibility. |
| `status_reason` | text null | Reason when status is `needs_review` or `invalid`. |
| `relationship_hash` | text not null | Stable hash generated from selector hash, target card, and relationship version. |
| `processed_at` | timestamptz not null | Timestamp when this relationship was generated. |

Supported values:

- `match_kind`: `exact_name`, `criteria`, `archetype`, `manual`.
- `status`: `accepted`, `needs_review`, `invalid`.

Constraints and indexes:

- Unique `card_selector_id, target_card_id, relationship_version`.
- Unique `relationship_hash`.
- Index on `target_card_id, status`.
- Index on `card_selector_id, status`.
- Index on `relationship_version`.

Notes:

- MVP public query filters `search_relationships.status = accepted`, then matches accepted `extracted_effects.actions_json` entries that reference the relationship's `card_selector_id` with `selector_role = target`. The action scope still requires `action_kind = move_card`, `verb = add`, `config.from = deck`, and `config.to = hand`.

Example record:

```json
{
  "id": "00000000-0000-4000-8000-000000000007",
  "card_selector_id": "00000000-0000-4000-8000-000000000040",
  "target_card_id": "00000000-0000-4000-8000-000000000019",
  "relationship_version": "relationships-1",
  "match_kind": "criteria",
  "status": "accepted",
  "status_reason": null,
  "relationship_hash": "sha256:selector-monster-warrior-level-lte-4:target:relationships-1",
  "processed_at": "2026-04-30T12:01:00Z"
}
```

## User Reports

### `user_reports`

Description:

Stores anonymous user reports for incorrect results, missing results, or other feedback. Public report creation needs lightweight spam protection.

Fields:

| Field | Type | Description |
| --- | --- | --- |
| `id` | uuid primary key | Internal report ID. |
| `target_card_id` | uuid null references `cards(id)` | Target card related to the report, when provided. |
| `searcher_card_id` | uuid null references `cards(id)` | Searcher card related to the report, when provided. |
| `relationship_id` | uuid null references `search_relationships(id)` | Existing relationship related to the report, when provided. |
| `report_type` | text not null | Category of user report. |
| `message` | text not null | User-provided report message. |
| `contact` | text null | Optional contact detail if the user wants follow-up. |
| `status` | text not null | Internal review status. |
| `created_at` | timestamptz not null | Timestamp when the report was created. |
| `metadata` | jsonb not null default `{}` | Audit or abuse-control metadata, such as client hints or validation flags. |

Supported values:

- `report_type`: `missing_result`, `incorrect_result`, `other`.
- `status`: `open`, `reviewed`, `resolved`, `rejected`.

Constraints and indexes:

- Index on `target_card_id`.
- Index on `searcher_card_id`.
- Index on `relationship_id`.
- Index on `status, created_at`.

Example records:

```json
{
  "id": "00000000-0000-4000-8000-000000000008",
  "target_card_id": "00000000-0000-4000-8000-000000000019",
  "searcher_card_id": "00000000-0000-4000-8000-000000000005",
  "relationship_id": "00000000-0000-4000-8000-000000000007",
  "report_type": "incorrect_result",
  "message": "This card should not be listed for the selected target.",
  "contact": null,
  "status": "open",
  "created_at": "2026-04-30T12:05:00Z",
  "metadata": {}
}
```

```json
{
  "id": "00000000-0000-4000-8000-000000000015",
  "target_card_id": "00000000-0000-4000-8000-000000000001",
  "searcher_card_id": null,
  "relationship_id": null,
  "report_type": "missing_result",
  "message": "A card that can search this target appears to be missing.",
  "contact": "player@example.com",
  "status": "open",
  "created_at": "2026-04-30T12:10:00Z",
  "metadata": {
    "honeypot": "empty"
  }
}
```

## Open Schema Questions

- No open schema questions right now.
