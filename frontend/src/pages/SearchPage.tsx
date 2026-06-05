import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { getCard } from "../api/client";
import type { CardSummary } from "../api/types";
import { CardFace } from "../components/CardFace";
import { CardSearch } from "../components/CardSearch";
import { CardSearchHero } from "../components/CardSearchHero";
import { ResultList } from "../components/ResultList";

export function SearchPage() {
  const [selectedCard, setSelectedCard] = useState<CardSummary | null>(null);

  const targetQuery = useQuery({
    queryKey: ["cards", selectedCard?.id],
    queryFn: () => getCard(selectedCard!.id),
    enabled: Boolean(selectedCard),
  });
  const targetCard = selectedCard ? targetQuery.data : undefined;

  return (
    <main className="app-shell">
      <CardSearchHero isSelected={Boolean(targetCard)}>
        <CardFace card={targetCard} />
        <CardSearch onClear={() => setSelectedCard(null)} onSelect={setSelectedCard} selectedName={targetCard?.name} />
        {targetQuery.isError ? <p className="error-text">{targetQuery.error.message}</p> : null}
      </CardSearchHero>

      <div className="results-layout">
        <ResultList targetID={targetCard?.id} />
      </div>
    </main>
  );
}
