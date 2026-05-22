-- name: GetCard :one
SELECT
    id::text,
    upstream_source,
    upstream_id,
    passcode,
    konami_id,
    name,
    normalized_name,
    aliases,
    normalized_aliases,
    description,
    card_type,
    frame_type,
    race,
    attribute,
    monster_categories,
    spell_trap_type,
    atk,
    def,
    level,
    rank,
    link_rating,
    archetype,
    mentions,
    text_features,
    image_url,
    ai_processing,
    raw_payload,
    imported_at::text,
    updated_at::text
FROM cards
WHERE id = $1;

-- name: SearchCards :many
SELECT id::text, name, image_url
FROM cards
WHERE normalized_name ILIKE '%' || $1 || '%'
   OR normalized_name % $1
   OR $1 = ANY(normalized_aliases)
ORDER BY
    CASE WHEN normalized_name = $1 THEN 0 ELSE 1 END,
    similarity(normalized_name, $1) DESC,
    name ASC
LIMIT $2;
