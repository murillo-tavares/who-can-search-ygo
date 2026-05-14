# Product Context

## Product

Who Can Search YGO is a web application for discovering which Yu-Gi-Oh! cards can add a selected target card from the Deck to the hand.

The core user question is:

```txt
Who can add this card from the Deck to the hand?
```

The MVP answers the classic search question: which cards can add a target card from the Deck to the hand.

The first supported wording family is:

```txt
Add ... from your Deck to your hand
```

The data model keeps parsed effect actions separate from target relationships so normal requests never parse card text.

## Audience

The initial audience is Yu-Gi-Oh! players, deck builders, content creators, and competitive researchers.

The product should work for casual and advanced players, but all terminology, card text, and card data must follow official English Yu-Gi-Oh! wording from the selected source.

## MVP Scope

The MVP must allow users to:

- Search for a card by name.
- Choose a target card from search results.
- List cards with accepted precomputed relationships that can get the target card from the Deck to the hand.
- Extract and store supported effect records from card text so relationship preprocessing can match targets from structured effect data.
- Support exact-name search effects.
- Support common criteria-based effects:
  - card category, such as Monster, Spell, or Trap;
  - monster type/race;
  - Attribute;
  - ATK and DEF thresholds;
  - Level, Rank, and Link Rating;
  - archetype when reliably represented by source data.
- Maintain an internal card database populated from an external card API.
- Precompute search relationships so normal user queries are fast.
- Let users report incorrect or missing search results.

## Out Of Scope For MVP

- Public result filters outside the default Deck-to-hand add flow.
- Relationship generation outside accepted `add` from `deck` to `hand` actions.
- Indirect combo discovery.
- Full natural-language understanding of every card text nuance.
- Cost, restriction, ruling, or competitive legality analysis.
- User login.
- Favorites.
- Deck builder.
- Admin UI.
- Card universes outside the selected initial API/dataset unless explicitly added through a scope change.

## Product Principles

- The first screen should be the useful search experience, not a marketing landing page.
- Public results should favor correctness over quantity.
- Public searcher results must show only `accepted` relationships.
- Normal user requests must never scan every card text.
- Keep the initial Deck-to-hand focus narrow enough that rule quality can be tested well.
- Start rule and importer validation from fixtures, mock payloads, or curated test cards. Use the full external dataset only after importer and rule failures are easy to debug.
- The UI should feel immersive and expressive while staying useful: prefer fluid interactions, thoughtful transitions, and creative visual treatment when they improve the search experience.

## Language Standard

Everything in the project must be written in English:

- documentation;
- code identifiers;
- file and folder names;
- API routes and payload fields;
- database tables and columns;
- comments;
- commit messages;
- UI text;
- internal technical terminology.

Card names and card text must use official English wording from the selected source.
