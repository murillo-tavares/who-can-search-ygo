-- name: ListActiveResolvedEffects :many
SELECT
    ce.id::text,
    ce.source_card_id::text,
    ce.effect_code,
    ce.extraction_version,
    ce.source_text,
    ce.condition_text,
    ce.cost_text,
    ce.action_text,
    ce.restriction_text,
    ce.selector_status,
    ce.selector_json,
    ce.ai_confidence,
    ce.is_active,
    source.id::text AS source_id,
    source.name AS source_name,
    source.image_url AS source_image_url
FROM card_effects ce
JOIN cards source ON source.id = ce.source_card_id
WHERE ce.effect_code = $1
  AND ce.selector_status = 'resolved'
  AND ce.is_active = true
ORDER BY source.name ASC, ce.id ASC;
