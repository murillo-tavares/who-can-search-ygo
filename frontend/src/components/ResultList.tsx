import { useQueries, useQuery } from "@tanstack/react-query";
import { getCard, getCardSearchers } from "../api/client";
import type { Card } from "../api/types";
import { EffectText } from "./EffectText";

type ResultListProps = {
  targetID?: string;
};

export function ResultList({ targetID }: ResultListProps) {
  const searchersQuery = useQuery({
    queryKey: ["cards", targetID, "searchers"],
    queryFn: () => getCardSearchers(targetID!),
    enabled: Boolean(targetID),
  });

  const sourceDetailQueries = useQueries({
    queries: (searchersQuery.data?.results ?? []).map((result) => ({
      queryKey: ["cards", result.source_card.id],
      queryFn: () => getCard(result.source_card.id),
      staleTime: 60_000,
    })),
  });

  const sourceDetails = new Map<string, Card>();
  sourceDetailQueries.forEach((query) => {
    if (query.data) {
      sourceDetails.set(query.data.id, query.data);
    }
  });

  return (
    <section className="results-panel" aria-label="Searcher results" hidden={!targetID}>
      <div className="section-heading">
        <div>
          <p className="eyebrow">Effect filter</p>
          <h2>Add from Deck to hand</h2>
        </div>
        {searchersQuery.data ? <span>{searchersQuery.data.results.length} matches</span> : null}
      </div>

      {!targetID ? <p className="empty-state">Select a target card to see matching searchers.</p> : null}
      {searchersQuery.isLoading ? <p className="empty-state">Loading matching effects...</p> : null}
      {searchersQuery.isError ? <p className="error-text">{searchersQuery.error.message}</p> : null}
      {searchersQuery.data?.results.length === 0 ? (
        <p className="empty-state">No active resolved searchers found for this target.</p>
      ) : null}

      {searchersQuery.data?.results.length ? (
        <ol className="result-list">
          {searchersQuery.data.results.map((result) => {
            const sourceCard = sourceDetails.get(result.source_card.id);

            return (
              <li className="result-row" key={result.effect_id}>
                <div className="source-card-art-placeholder" aria-hidden="true" />
                <div className="source-card-box">
                  <div className="source-card-header">
                    <h3>{result.source_card.name}</h3>
                    {sourceCard ? <span>{buildHeaderBadge(sourceCard)}</span> : null}
                  </div>
                  {sourceCard ? <SourceCardMeta card={sourceCard} /> : <p className="muted">Loading card details...</p>}
                  <EffectText hasTypeLine={sourceCard?.card_type === "Monster"} result={result} fullText={sourceCard?.description} />
                  {sourceCard?.card_type === "Monster" ? (
                    <div className="source-card-stats">
                      <span>ATK/{sourceCard.atk ?? "?"}</span>
                      <span>DEF/{sourceCard.def ?? "?"}</span>
                    </div>
                  ) : null}
                </div>
              </li>
            );
          })}
        </ol>
      ) : null}
    </section>
  );
}

function SourceCardMeta({ card }: { card: Card }) {
  if (card.card_type === "Monster") {
    return (
      <>
        <div className="source-card-level" aria-label={buildMonsterLevelLabel(card)}>
          {buildMonsterLevel(card)}
        </div>
        <div className="source-card-type">{buildMonsterTypeLine(card)}</div>
      </>
    );
  }

  if (card.card_type === "Spell" || card.card_type === "Trap") {
    return <div className="source-card-subtype">{card.spell_trap_type ?? "Normal"}</div>;
  }

  return <div className="source-card-type">{[card.card_type, card.frame_type].filter(Boolean).join(" / ")}</div>;
}

function buildHeaderBadge(card: Card) {
  if (card.card_type === "Monster") {
    return card.attribute ?? "Monster";
  }
  return card.card_type ?? "Card";
}

function buildMonsterLevel(card: Card) {
  if (card.link_rating !== null) {
    return `LINK-${card.link_rating}`;
  }
  const count = card.level ?? card.rank ?? card.link_rating ?? 0;
  if (count <= 0) {
    return "";
  }
  return Array.from({ length: count }, () => "★").join(" ");
}

function buildMonsterLevelLabel(card: Card) {
  if (card.level !== null) {
    return `Level ${card.level}`;
  }
  if (card.rank !== null) {
    return `Rank ${card.rank}`;
  }
  if (card.link_rating !== null) {
    return `Link Rating ${card.link_rating}`;
  }
  return "No Level";
}

function buildMonsterTypeLine(card: Card) {
  const parts = [card.race, ...card.monster_categories].filter(Boolean);
  return parts.length ? `[${parts.join(" / ")}]` : "[Monster]";
}
