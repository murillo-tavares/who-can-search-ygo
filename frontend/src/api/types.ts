export type CardSummary = {
  id: string;
  name: string;
  image_url?: string | null;
};

export type Card = CardSummary & {
  normalized_name: string;
  description?: string;
  card_type: string | null;
  frame_type?: string | null;
  race: string | null;
  attribute: string | null;
  monster_categories: string[];
  spell_trap_type: string | null;
  atk: number | null;
  def: number | null;
  level: number | null;
  rank: number | null;
  link_rating: number | null;
  archetype: string | null;
  mentions: string[];
  text_features: string[];
};

export type CardSearchResponse = {
  results: CardSummary[];
};

export type SearcherResult = {
  effect_id: string;
  source_card: CardSummary;
  source_text: string;
  condition_text: string | null;
  cost_text: string | null;
  action_text: string;
  restriction_text: string | null;
};

export type SearchersResponse = {
  target_card: {
    id: string;
    name: string;
  };
  effect_code: "add_deck_to_hand";
  results: SearcherResult[];
};
