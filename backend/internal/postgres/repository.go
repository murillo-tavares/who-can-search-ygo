package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"who-can-search-ygo/backend/internal/domain"
	"who-can-search-ygo/backend/internal/service"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetCard(ctx context.Context, id string) (domain.Card, error) {
	const query = `
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
	image_url,
	ai_processing,
	raw_payload,
	imported_at::text,
	updated_at::text
FROM cards
WHERE id = $1
`

	card, err := scanCard(r.pool.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Card{}, service.ErrNotFound
	}
	if err != nil {
		return domain.Card{}, fmt.Errorf("get card: %w", err)
	}
	return card, nil
}

func (r *Repository) SearchCards(ctx context.Context, query string, limit int) ([]domain.CardSummary, error) {
	const sql = `
SELECT id::text, name, image_url
FROM cards
WHERE normalized_name ILIKE '%' || $1 || '%'
   OR normalized_name % $1
   OR $1 = ANY(normalized_aliases)
ORDER BY
  CASE WHEN normalized_name = $1 THEN 0 ELSE 1 END,
  similarity(normalized_name, $1) DESC,
  name ASC
LIMIT $2
`

	rows, err := r.pool.Query(ctx, sql, normalizeName(query), limit)
	if err != nil {
		return nil, fmt.Errorf("search cards: %w", err)
	}
	defer rows.Close()

	var cards []domain.CardSummary
	for rows.Next() {
		var card domain.CardSummary
		var imageURL pgtype.Text
		if err := rows.Scan(&card.ID, &card.Name, &imageURL); err != nil {
			return nil, fmt.Errorf("scan card summary: %w", err)
		}
		card.ImageURL = textPtr(imageURL)
		cards = append(cards, card)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate card summaries: %w", err)
	}
	return cards, nil
}

func (r *Repository) ListActiveResolvedEffects(ctx context.Context, effectCode string) ([]domain.CardEffect, error) {
	const query = `
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
	source.id::text,
	source.name,
	source.image_url
FROM card_effects ce
JOIN cards source ON source.id = ce.source_card_id
WHERE ce.effect_code = $1
  AND ce.selector_status = 'resolved'
  AND ce.is_active = true
ORDER BY source.name ASC, ce.id ASC
`

	rows, err := r.pool.Query(ctx, query, effectCode)
	if err != nil {
		return nil, fmt.Errorf("list active resolved effects: %w", err)
	}
	defer rows.Close()

	var effects []domain.CardEffect
	for rows.Next() {
		var effect domain.CardEffect
		var conditionText, costText, restrictionText, imageURL pgtype.Text
		var aiConfidence pgtype.Numeric
		var selectorJSON []byte
		if err := rows.Scan(
			&effect.ID,
			&effect.SourceCardID,
			&effect.EffectCode,
			&effect.ExtractionVersion,
			&effect.SourceText,
			&conditionText,
			&costText,
			&effect.ActionText,
			&restrictionText,
			&effect.SelectorStatus,
			&selectorJSON,
			&aiConfidence,
			&effect.IsActive,
			&effect.SourceCard.ID,
			&effect.SourceCard.Name,
			&imageURL,
		); err != nil {
			return nil, fmt.Errorf("scan active resolved effect: %w", err)
		}
		effect.ConditionText = textPtr(conditionText)
		effect.CostText = textPtr(costText)
		effect.RestrictionText = textPtr(restrictionText)
		effect.SourceCardID = effect.SourceCard.ID
		effect.SourceCard.ImageURL = textPtr(imageURL)
		effect.SelectorJSON = json.RawMessage(selectorJSON)
		if aiConfidence.Valid {
			value, err := aiConfidence.Float64Value()
			if err == nil && value.Valid {
				effect.AIConfidence = &value.Float64
			}
		}
		effects = append(effects, effect)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate active resolved effects: %w", err)
	}
	return effects, nil
}

func scanCard(row pgx.Row) (domain.Card, error) {
	var card domain.Card
	var passcode, konamiID, cardType, frameType, race, attribute, spellTrapType, archetype, imageURL pgtype.Text
	var atk, def, level, rank, linkRating pgtype.Int4
	var aliases, aiProcessing, rawPayload []byte

	err := row.Scan(
		&card.ID,
		&card.UpstreamSource,
		&card.UpstreamID,
		&passcode,
		&konamiID,
		&card.Name,
		&card.NormalizedName,
		&aliases,
		&card.NormalizedAliases,
		&card.Description,
		&cardType,
		&frameType,
		&race,
		&attribute,
		&card.MonsterCategories,
		&spellTrapType,
		&atk,
		&def,
		&level,
		&rank,
		&linkRating,
		&archetype,
		&card.Mentions,
		&imageURL,
		&aiProcessing,
		&rawPayload,
		&card.ImportedAt,
		&card.UpdatedAt,
	)
	if err != nil {
		return domain.Card{}, err
	}

	if err := json.Unmarshal(aliases, &card.Aliases); err != nil {
		return domain.Card{}, fmt.Errorf("decode aliases: %w", err)
	}
	card.Passcode = textPtr(passcode)
	card.KonamiID = textPtr(konamiID)
	card.CardType = textPtr(cardType)
	card.FrameType = textPtr(frameType)
	card.Race = textPtr(race)
	card.Attribute = textPtr(attribute)
	card.SpellTrapType = textPtr(spellTrapType)
	card.ATK = intPtr(atk)
	card.DEF = intPtr(def)
	card.Level = intPtr(level)
	card.Rank = intPtr(rank)
	card.LinkRating = intPtr(linkRating)
	card.Archetype = textPtr(archetype)
	card.ImageURL = textPtr(imageURL)
	card.AIProcessing = json.RawMessage(aiProcessing)
	card.RawPayload = json.RawMessage(rawPayload)
	return card, nil
}

func textPtr(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func intPtr(value pgtype.Int4) *int {
	if !value.Valid {
		return nil
	}
	converted := int(value.Int32)
	return &converted
}

func normalizeName(value string) string {
	return strings.Join(strings.Fields(strings.ToLower(value)), " ")
}
