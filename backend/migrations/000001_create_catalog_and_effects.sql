-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE cards (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    upstream_source text NOT NULL,
    upstream_id text NOT NULL,
    passcode text NULL,
    konami_id text NULL,
    name text NOT NULL,
    normalized_name text NOT NULL,
    aliases jsonb NOT NULL DEFAULT '[]'::jsonb,
    normalized_aliases text[] NOT NULL DEFAULT '{}',
    description text NOT NULL,
    card_type text NULL,
    frame_type text NULL,
    race text NULL,
    attribute text NULL,
    monster_categories text[] NOT NULL DEFAULT '{}',
    spell_trap_type text NULL,
    atk integer NULL,
    def integer NULL,
    level integer NULL,
    rank integer NULL,
    link_rating integer NULL,
    archetype text NULL,
    mentions text[] NOT NULL DEFAULT '{}',
    image_url text NULL,
    ai_processing jsonb NOT NULL DEFAULT '{}'::jsonb,
    raw_payload jsonb NOT NULL,
    imported_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT cards_upstream_unique UNIQUE (upstream_source, upstream_id)
);

CREATE INDEX cards_normalized_name_trgm_idx ON cards USING gin (normalized_name gin_trgm_ops);
CREATE INDEX cards_normalized_aliases_gin_idx ON cards USING gin (normalized_aliases);
CREATE INDEX cards_mentions_gin_idx ON cards USING gin (mentions);
CREATE INDEX cards_card_type_idx ON cards (card_type);
CREATE INDEX cards_race_idx ON cards (race);
CREATE INDEX cards_attribute_idx ON cards (attribute);
CREATE INDEX cards_level_idx ON cards (level);
CREATE INDEX cards_rank_idx ON cards (rank);
CREATE INDEX cards_link_rating_idx ON cards (link_rating);
CREATE INDEX cards_atk_idx ON cards (atk);
CREATE INDEX cards_def_idx ON cards (def);
CREATE INDEX cards_archetype_idx ON cards (archetype);

CREATE TABLE card_effects (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    source_card_id uuid NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    effect_code text NOT NULL,
    extraction_version text NOT NULL,
    source_text text NOT NULL,
    condition_text text NULL,
    cost_text text NULL,
    action_text text NOT NULL,
    restriction_text text NULL,
    selector_status text NOT NULL,
    selector_json jsonb NULL,
    selector_hash text NOT NULL,
    ai_confidence numeric NULL,
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT card_effects_selector_status_check CHECK (selector_status IN ('resolved', 'unresolved')),
    CONSTRAINT card_effects_resolved_selector_check CHECK (
        (selector_status = 'resolved' AND selector_json IS NOT NULL)
        OR (selector_status = 'unresolved' AND selector_json IS NULL)
    )
);

CREATE INDEX card_effects_source_card_id_idx ON card_effects (source_card_id);
CREATE INDEX card_effects_public_lookup_idx ON card_effects (effect_code, is_active, selector_status);
CREATE INDEX card_effects_selector_hash_idx ON card_effects (selector_hash);
CREATE UNIQUE INDEX card_effects_dedupe_idx
    ON card_effects (source_card_id, effect_code, extraction_version, selector_hash, source_text);

-- +goose Down
DROP TABLE IF EXISTS card_effects;
DROP TABLE IF EXISTS cards;
DROP EXTENSION IF EXISTS pg_trgm;
DROP EXTENSION IF EXISTS pgcrypto;
