# Product Context

## Product

Who Can Search YGO is a web application for finding Yu-Gi-Oh! cards whose own effects can affect a selected target card.

The first product question is:

```txt
Which cards can add this selected card from the Deck to the hand?
```

The initial supported effect code is:

```txt
Add ... from your Deck to your hand
```

Example:

`Reinforcement of the Army` has:

```txt
Add 1 Level 4 or lower Warrior monster from your Deck to your hand.
```

For a target card that is a Level 4 or lower Warrior monster, the application should show `Reinforcement of the Army` as a result.

## Audience

The initial audience is:

- Yu-Gi-Oh! players.
- Deck builders.
- Content creators.
- Competitive researchers.

The product should be useful for casual and advanced players while staying clear, fast, and trustworthy.

All card names, card text, and card terminology must use official English wording from the selected data source.

## MVP Scope

The MVP must allow users to:

- Search cards by name.
- Select one target card.
- Use the fixed `add_deck_to_hand` effect filter.
- See all stored active card effects whose selector matches the selected target card.
- See the matched source effect text so the UI can highlight why the source card appears.
- See useful source card information in results.
- Store an internal card catalog imported from an external card data source.
- Store image URLs from the card data source without downloading images locally.
- Store AI-extracted selectors available in the database.
- Track which AI-supported processes each card has completed and which extraction version completed them.

The application owns:

- card catalog storage;
- selector storage;
- per-card AI processing metadata;
- public search API;
- frontend search experience.

The application works with the card and extraction data currently available in the database. It does not need to know who produced AI extraction data, how that process runs, or how often it runs.

## Out Of Scope For MVP

- User login.
- Favorites.
- Deck builder.
- Public admin UI.
- Live AI calls during public requests.
- Full Yu-Gi-Oh! rules simulation.
- Full PSCT parsing.
- Exact activation legality.
- Costs, rulings, restrictions, or chain interactions.
- Indirect combo discovery.
- Support for effect filters other than `add_deck_to_hand`.
- Downloading, transforming, or storing card images locally.

## Product Principles

- The first screen should be the card search experience, not a marketing page.
- Public results should favor correctness over quantity.
- AI extraction is assistive, not authoritative.
- Public results should use stored active extracted effects with resolved selectors.
- Public requests must not call AI.
- Public requests must not scan all card text.
- Public requests may evaluate stored active selectors against the selected target card in real time.
- Missing a relationship is better than showing an incorrect one.
- Keep the first effect code narrow until the data model and extraction workflow are reliable.

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
