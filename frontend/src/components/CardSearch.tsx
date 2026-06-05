import { useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { searchCards } from "../api/client";
import type { CardSummary } from "../api/types";

type CardSearchProps = {
  onClear: () => void;
  onSelect: (card: CardSummary) => void;
  selectedName?: string;
};

export function CardSearch({ onClear, onSelect, selectedName }: CardSearchProps) {
  const [query, setQuery] = useState("");
  const [isAutocompleteOpen, setIsAutocompleteOpen] = useState(false);
  const debouncedQuery = useDebouncedValue(query.trim(), 250);

  const searchQuery = useQuery({
    queryKey: ["cards", "search", debouncedQuery],
    queryFn: () => searchCards(debouncedQuery),
    enabled: debouncedQuery.length >= 2,
  });

  function handleSelect(card: CardSummary) {
    setQuery(card.name);
    setIsAutocompleteOpen(false);
    onSelect(card);
  }

  function handleQueryChange(value: string) {
    setQuery(value);
    setIsAutocompleteOpen(value.trim().length >= 2);
    if (value === "") {
      setIsAutocompleteOpen(false);
      onClear();
      return;
    }
    if (selectedName && value !== selectedName) {
      onClear();
    }
  }

  function handleFocus() {
    if (query.trim().length >= 2 && query !== selectedName) {
      setIsAutocompleteOpen(true);
    }
  }

  return (
    <section className="filter-search" aria-label="Card search">
      <div className="filter-line">
        <span>Who can</span>
        <strong>add</strong>
        <input
          id="card-search"
          type="search"
          value={selectedName && query === "" ? selectedName : query}
          onChange={(event) => handleQueryChange(event.target.value)}
          onFocus={handleFocus}
          placeholder="target card"
          autoComplete="off"
          aria-label="Target card"
        />
        <span>from</span>
        <strong>Deck</strong>
        <span>to</span>
        <strong>hand</strong>
        <span>?</span>
      </div>
      <div className="autocomplete-results" aria-live="polite">
        {isAutocompleteOpen ? (
          <>
        {query.trim().length > 0 && query.trim().length < 2 ? (
          <p className="muted">Type at least 2 characters.</p>
        ) : null}
        {searchQuery.isLoading ? <p className="muted">Searching cards...</p> : null}
        {searchQuery.isError ? <p className="error-text">{searchQuery.error.message}</p> : null}
        {searchQuery.data?.length === 0 ? <p className="muted">No cards found.</p> : null}
        {searchQuery.data && searchQuery.data.length > 0 ? (
          <ul>
            {searchQuery.data.map((card) => (
              <li key={card.id}>
                <button type="button" onClick={() => handleSelect(card)}>
                  <span className="thumb-placeholder" aria-hidden="true" />
                  <span>{card.name}</span>
                </button>
              </li>
            ))}
          </ul>
        ) : null}
          </>
        ) : null}
      </div>
    </section>
  );
}

function useDebouncedValue<T>(value: T, delayMs: number) {
  const [debouncedValue, setDebouncedValue] = useState(value);

  useEffect(() => {
    const timeoutID = window.setTimeout(() => setDebouncedValue(value), delayMs);
    return () => window.clearTimeout(timeoutID);
  }, [delayMs, value]);

  return debouncedValue;
}
