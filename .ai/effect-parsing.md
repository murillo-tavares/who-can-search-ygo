# Effect Parsing Guide

This guide defines how deterministic effect parsing should interpret card text into `card_effects`.

The parser model must support multiple effect codes over time. `add_deck_to_hand` is only the first supported profile.

The parser is intentionally conservative. If it can detect a supported action but cannot build a correct selector from supported fields, it should write an unresolved effect or skip writing the effect according to the extraction workflow.

## Core Concepts

An effect parser has two layers:

- a generic text segmentation and selector-building layer;
- one or more effect profiles, each tied to an `effect_code`.

An effect profile defines:

- supported action verb or action family;
- source zone requirements, when relevant;
- destination zone requirements, when relevant;
- whether matching evaluates target cards, source cards, or another entity;
- allowed target phrase templates;
- profile-specific exclusions;
- profile-specific unresolved cases.

The same selector expression model should be reused across profiles whenever possible.

## Parser Output Contract

A deterministic parser should produce the same logical record shape as external extraction:

- `effect_code`;
- exact `source_text`;
- optional `condition_text`;
- optional `cost_text`;
- exact isolated `action_text`;
- optional `restriction_text`;
- `selector_status`;
- `selector_json` when resolved.

The parser must preserve official card wording. It must not summarize or paraphrase source text.

If selector generation is uncertain, use:

```json
{
  "selector_status": "unresolved",
  "selector_json": null
}
```

## Generic Text Segmentation

Use PSCT punctuation conservatively:

- text before `:` is usually `condition_text` when it is timing or condition wording;
- text before `;` is usually `cost_text` when it is an explicit cost;
- `source_text` is the full segment containing the supported action;
- `action_text` is the isolated supported action plus any earlier wording needed to identify the target;
- `restriction_text` is post-action limitation text such as `also, you cannot...` or `for the rest of this turn...` when it is not part of the selector.

Do not split aggressively when punctuation is ambiguous. Keep `source_text` exact and mark uncertain fields null.

If a sentence contains multiple actions, extract the supported profile action plus text needed to preserve official meaning. Other actions only belong in `action_text` when they identify the target or are inseparable from the supported action.

## Generic Selector Templates

These templates are safe when all named values are explicit and supported by current selector fields.

### Exact Named Card

Text shape:

```txt
<action> "Dark Magician" ...
```

Selector:

```json
{ "type": "comparison", "field": "name", "op": "=", "value": "Dark Magician" }
```

### Archetype Or Named Family

Text shape:

```txt
<action> "K9" monster ...
```

Selector:

```json
{
  "type": "logical",
  "op": "and",
  "args": [
    { "type": "comparison", "field": "archetype", "op": "=", "value": "K9" },
    { "type": "comparison", "field": "card_type", "op": "=", "value": "Monster" }
  ]
}
```

Use this for `"<Name>" card`, `"<Name>" monster`, `"<Name>" Spell`, `"<Name>" Trap`, and `"<Name>" Spell/Trap` when the quoted value is a family selector rather than one exact card.

### Monster Filters

Supported static filters include:

- `monster_categories`, such as `Ritual`, `Pendulum`, `Tuner`, `Union`, `Spirit`, `Toon`, and `Normal`;
- monster `race`;
- `attribute`;
- exact or bounded `level`;
- exact or bounded `rank`;
- exact or bounded `link_rating`;
- exact or bounded `atk`;
- exact or bounded `def`;
- exact or bounded `combined_atk_def`.

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

For monster categories, use `contains`:

```json
{
  "type": "logical",
  "op": "and",
  "args": [
    { "type": "comparison", "field": "card_type", "op": "=", "value": "Monster" },
    { "type": "comparison", "field": "monster_categories", "op": "contains", "value": "Ritual" }
  ]
}
```

For combined ATK and DEF requirements, use derived selector field `combined_atk_def`:

```json
{
  "type": "comparison",
  "field": "combined_atk_def",
  "op": "=",
  "value": 2000
}
```

### Spell And Trap Filters

For `Spell/Trap`, use:

```json
{ "type": "comparison", "field": "card_type", "op": "in", "value": ["Spell", "Trap"] }
```

For specific Spell/Trap subtypes, combine `card_type` and `spell_trap_type`.

Example:

```json
{
  "type": "logical",
  "op": "and",
  "args": [
    { "type": "comparison", "field": "card_type", "op": "=", "value": "Spell" },
    { "type": "comparison", "field": "spell_trap_type", "op": "=", "value": "Field" }
  ]
}
```

### Mentions Or Specifically Lists

Use `mentions contains "<Card Name>"` for official wording like:

- `mentions "<Card Name>"`;
- `specifically lists the card "<Card Name>"`;
- `lists "<Card Name>" in its text`.

Use the same field for explicit family or archetype mention requirements when the wording refers to a named family instead of one exact card, such as:

- `mentions a "Destiny HERO" monster's card name`;
- `specifically lists "Iron Core of Koa'ki Meiru" in its text`;
- `mentions "Fallen of Albaz"`.

For family wording, the value should be the family or exact phrase named by the effect, for example `mentions contains "Destiny HERO"`.

Example:

```json
{
  "type": "logical",
  "op": "and",
  "args": [
    { "type": "comparison", "field": "card_type", "op": "in", "value": ["Spell", "Trap"] },
    { "type": "comparison", "field": "mentions", "op": "contains", "value": "Dark Magician" }
  ]
}
```

### Text Feature Filters

Use `text_features contains "<feature_code>"` when a target requirement refers to a normalized text feature.

Examples:

- `has a coin tossing effect` -> `coin_toss_effect`;
- `has an effect that places a counter` -> `places_counter_effect`.
- `requires a die roll` -> `dice_roll_effect`;
- `cannot be Normal Summoned/Set` -> `cannot_be_normal_summoned_set`.

Selector:

```json
{
  "type": "comparison",
  "field": "text_features",
  "op": "contains",
  "value": "coin_toss_effect"
}
```

Feature codes should be stable snake_case strings. Add new feature codes only when a parser or import process can populate them consistently on cards.

### Alternatives

Use `or` only when each branch is independently parseable.

Example:

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

### Exclusions

Exclusions belong in both `action_text` and `selector_json`. Do not put selector exclusions in `restriction_text`.

Example:

```json
{
  "type": "logical",
  "op": "and",
  "args": [
    { "type": "comparison", "field": "archetype", "op": "=", "value": "Abyss Actor" },
    { "type": "comparison", "field": "name", "op": "!=", "value": "Abyss Actor - Wild Hope" }
  ]
}
```

## Generic Selector Simplifications

These cases can still produce resolved selectors when a non-selector constraint or dynamic comparison can be intentionally ignored and at least one useful static selector filter remains.

### Ignorable Runtime Comparison Filters

Examples:

- `same name as the banished Spell Card`;
- `same Attribute`;
- `different name than that revealed Spell`;
- `different names from each other`;
- `different from cards you control`;
- `declared card`;
- `that monster`.

Rule:

When the referenced object depends on runtime player choice and cannot be statically defined from the effect text, ignore the dynamic comparison filter and preserve any static filters that remain.

Example:

```txt
Add up to 2 Spells/Traps with different names from each other that mention "Dark Time Wizard" from your Deck to your hand.
```

The `different names from each other` filter is ignored. The selector should still require `card_type in ["Spell", "Trap"]` and `mentions contains "Dark Time Wizard"`.

If ignoring the dynamic comparison leaves no static selector filter, fail closed to unresolved.

Pronouns such as `that monster`, `it`, or `that card` may still be resolved when an earlier phrase in the same effect provides a static selector. Preserve the static selector and ignore only the runtime identity.

Examples:

- if the effect previously identifies an `Artmage` monster and then says `add that monster from your Deck to your hand`, use `archetype = "Artmage"` and `card_type = "Monster"`;
- if an opponent declares a `Gunkan` monster's card name, use `archetype = "Gunkan"` and `card_type = "Monster"` and ignore which exact name was declared.
- if the effect sends a `Level 4 or lower Normal Monster` from the Deck and then adds a card with the same name as that sent card, use `card_type = "Monster"`, `monster_categories contains "Normal"`, and `level <= 4`.

Runtime numeric limits may be reduced only when the game gives a clear hard maximum. For example, `whose Level is less than or equal to the number of "Ghostrick" monsters you control` can use `level <= 6` plus `archetype = "Ghostrick"` and `card_type = "Monster"`, because a player normally controls at most 6 monsters. If no safe maximum is known, leave the runtime numeric filter unresolved or ignore it only when other useful static filters remain.

### Broad Generic Selectors

Some supported actions produce selectors that are technically valid but too broad for useful public results.

Examples:

- `add 1 card from your Deck to your hand`;
- `add 1 monster from your Deck to your hand`;
- `add 1 Spell from your Deck to your hand`;
- `add 1 Field Spell from your Deck to your hand`.

Rule:

Do not represent a 100% generic `card` selector. It would match every target and pollute public search results.

Only selectors that reduce to any card, any Monster, any Spell, or any Trap are considered too generic by default.

Store decided generic cases as `selector_status = ignored`, not `unresolved`. These are known parser outcomes and should not require later human review.

Do not discard useful broad selectors when they still have a meaningful family, category, text feature, stat, or named-card filter. For example:

- `any "Artmage" monster` -> `archetype = "Artmage"` and `card_type = "Monster"`;
- `any "Gunkan" monster` -> `archetype = "Gunkan"` and `card_type = "Monster"`;
- `any "Tenyi" monster` -> `archetype = "Tenyi"` and `card_type = "Monster"`;
- `any "Six Samurai" monster` -> `archetype = "Six Samurai"` and `card_type = "Monster"`.

Do not classify narrower selectors as generic spam just because they are broad. Any Field Spell, any archetype monster, any monster category, any text feature, any stat range, or any mention requirement still has a meaningful selector and may be resolved.

### Derived Aggregate Filters

Some aggregate-looking requirements can be reduced to one-target selector filters.

Combined ATK and DEF:

```txt
Add 1 Fiend monster whose combined ATK & DEF equal 2000 from your Deck to your hand.
```

Use `combined_atk_def = 2000`, plus any other static filters such as `race = "Fiend"`.

Total Level sets:

```txt
Add up to 3 Rock monsters from your Deck to your hand, whose total Levels equal 8.
```

For one-target matching, reduce this to `level <= 6` because the other 2 selected monsters have minimum Level 1. Preserve static filters, so this example becomes `card_type = "Monster"`, `race = "Rock"`, and `level <= 6`.

Formula:

```txt
max_target_level = total_level - (target_count - 1)
```

If `target_count`, `total_level`, or a valid positive `max_target_level` cannot be derived, fail closed to unresolved.

## Effect Profile: `add_deck_to_hand`

This profile answers:

```txt
Which cards can add this selected card from the Deck to the hand?
```

### Profile Semantics

- `effect_code`: `add_deck_to_hand`
- action family: add to hand
- required source zone: Deck
- required destination zone: hand
- public matching target: selected target card
- relevant target zone for alias matching: `deck`

### Common Broad Shapes

Common broad shapes are simple:

- `You can add <target> from your Deck to your hand.`
- `add <target> from your Deck to your hand.`
- the same shapes with `except "<Card Name>"`.

These shapes should be handled first, but parser correctness depends on selector semantics, not only action detection.

### Candidate Detection

Treat a card text segment as a candidate when it contains an action that adds one or more cards from the Deck to the hand.

Common direct forms:

- `Add 1 ... from your Deck to your hand.`
- `You can add 1 ... from your Deck to your hand.`
- `...; add 1 ... from your Deck to your hand.`
- `..., and if you do, add 1 ... from your Deck to your hand.`
- `..., then you can add 1 ... from your Deck to your hand.`

Common still-in-scope forms:

- `add it to your hand` when the same clause first identifies the revealed target as being from the Deck;
- `add that card/monster from your Deck to your hand` when `that card/monster` has a resolvable static selector;
- `look at/reveal/excavate ... from your Deck, then add ... to your hand` when the target selector is explicit;
- effects that can add from the Main Deck plus other zones, such as `Deck or Extra Deck`, because the Main Deck is included;
- `add ... from your Deck to your hand, except "<Card Name>"`.

Out-of-scope for this profile:

- adding from the GY, banishment, hand, field, or Extra Deck only;
- drawing cards;
- setting cards from the Deck;
- sending, summoning, equipping, or placing cards without adding them to the hand.

### Profile Examples

Exact named target:

```txt
Add 1 "Dark Magician" from your Deck to your hand.
```

Family target:

```txt
You can add 1 "K9" monster from your Deck to your hand.
```

Static typed target:

```txt
Add 1 Level 4 or lower Warrior monster from your Deck to your hand.
```

Monster category target:

```txt
Add 1 Ritual Monster from your Deck to your hand.
```

Mentions target:

```txt
Add 1 Spell/Trap that mentions "Dark Magician" from your Deck to your hand.
```

Mixed alternatives:

```txt
Add 1 "Dark Magician" or 1 Spell/Trap that mentions "Dark Magician" from your Deck to your hand.
```

Exclusion:

```txt
Add 1 "Abyss Actor" card from your Deck to your hand, except "Abyss Actor - Wild Hope".
```

### Profile-Specific Open Cases

The following observed `add_deck_to_hand` cases need explicit decisions before deterministic parsing should resolve them.

Keep them unresolved for now:

- indirect target identity where no static selector remains, such as Fusion Materials listed on a revealed Extra Deck monster or Ritual Monsters listed on a revealed/sent Ritual Spell;
- runtime numeric requirements that depend on game state, such as `whose sum of ATK and DEF equals your LP`;
- target requirements based on text features that have not been added to `text_features`;
- generic selectors that reduce to any card, any Monster, any Spell, or any Trap should be ignored, not unresolved.

## Implementation Guidance

Recommended parser shape:

1. Choose an effect profile by action and zone requirements.
2. Detect candidate action segments for that profile.
3. Split condition, cost, action, and restriction conservatively.
4. Build selectors only from explicit static target phrases.
5. Use recursive `and`/`or` selectors for alternatives.
6. Include exclusions as selector comparisons.
7. Ignore runtime comparison filters when they depend on player choice and static filters remain.
8. Reduce supported aggregate constraints to one-target selector filters when documented.
9. Mark as ignored when parsing succeeds but the resulting selector is any card, any Monster, any Spell, or any Trap.
10. Fail closed to unresolved when no static selector remains, when the target requirement cannot be represented by supported selector fields, or when aggregate constraints or grammar are ambiguous.

The first useful backend parser does not need to replace all external extraction. It can safely own common static templates and leave the rest as unresolved for later human review.
