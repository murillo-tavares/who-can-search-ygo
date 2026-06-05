import type { Card, CardSearchResponse, SearchersResponse } from "./types";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "/api";

type APIErrorBody = {
  error?: {
    code?: string;
    message?: string;
  };
};

async function requestJSON<T>(path: string): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`);

  if (!response.ok) {
    let message = "Request failed";
    try {
      const body = (await response.json()) as APIErrorBody;
      message = body.error?.message ?? message;
    } catch {
      message = response.statusText || message;
    }
    throw new Error(message);
  }

  return response.json() as Promise<T>;
}

export async function searchCards(query: string, limit = 12) {
  const params = new URLSearchParams({ query, limit: String(limit) });
  const response = await requestJSON<CardSearchResponse>(`/cards?${params}`);
  return response.results;
}

export async function getCard(id: string) {
  return requestJSON<Card>(`/cards/${id}`);
}

export async function getCardSearchers(id: string) {
  return requestJSON<SearchersResponse>(`/cards/${id}/searchers`);
}
