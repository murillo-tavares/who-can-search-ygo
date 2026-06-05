import { describe, expect, it } from "vitest";
import { splitEffectText } from "./EffectText";
import type { SearcherResult } from "../api/types";

describe("splitEffectText", () => {
  it("marks condition, cost, and action in source order", () => {
    const result: SearcherResult = {
      effect_id: "effect-1",
      source_card: { id: "card-1", name: "Example" },
      source_text:
        'During your Main Phase: You can discard 1 card; add 1 "Example" monster from your Deck to your hand.',
      condition_text: "During your Main Phase",
      cost_text: "You can discard 1 card",
      action_text: 'add 1 "Example" monster from your Deck to your hand.',
      restriction_text: null,
    };

    expect(splitEffectText(result)).toEqual([
      { text: "During your Main Phase", kind: "condition" },
      { text: ": " },
      { text: "You can discard 1 card", kind: "cost" },
      { text: "; " },
      { text: 'add 1 "Example" monster from your Deck to your hand.', kind: "action" },
    ]);
  });
});
