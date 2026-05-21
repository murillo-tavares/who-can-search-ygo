package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"who-can-search-ygo/backend/internal/domain"
)

type cardsFixture struct {
	Cards []domain.Card `json:"cards"`
}

type effectsFixture struct {
	Effects []domain.CardEffect `json:"card_effects"`
}

// SeedFixtures upserts local fixture data into PostgreSQL for development.
func SeedFixtures(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	cards, err := readCardsFixture(dir)
	if err != nil {
		return err
	}
	effects, err := readEffectsFixture(dir)
	if err != nil {
		return err
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin fixture seed: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, card := range cards {
		if err := seedCard(ctx, tx, card); err != nil {
			return err
		}
	}
	for _, effect := range effects {
		if err := seedEffect(ctx, tx, effect); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit fixture seed: %w", err)
	}
	return nil
}

func seedCard(ctx context.Context, tx pgx.Tx, card domain.Card) error {
	aliases, err := json.Marshal(card.Aliases)
	if err != nil {
		return fmt.Errorf("encode card aliases %s: %w", card.ID, err)
	}

	const query = `
INSERT INTO cards (
	id,
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
	imported_at,
	updated_at
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8::jsonb, $9, $10, $11, $12, $13, $14,
	$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25::jsonb, $26::jsonb,
	COALESCE(NULLIF($27, '')::timestamptz, now()),
	COALESCE(NULLIF($28, '')::timestamptz, now())
)
ON CONFLICT (id) DO UPDATE SET
	upstream_source = EXCLUDED.upstream_source,
	upstream_id = EXCLUDED.upstream_id,
	passcode = EXCLUDED.passcode,
	konami_id = EXCLUDED.konami_id,
	name = EXCLUDED.name,
	normalized_name = EXCLUDED.normalized_name,
	aliases = EXCLUDED.aliases,
	normalized_aliases = EXCLUDED.normalized_aliases,
	description = EXCLUDED.description,
	card_type = EXCLUDED.card_type,
	frame_type = EXCLUDED.frame_type,
	race = EXCLUDED.race,
	attribute = EXCLUDED.attribute,
	monster_categories = EXCLUDED.monster_categories,
	spell_trap_type = EXCLUDED.spell_trap_type,
	atk = EXCLUDED.atk,
	def = EXCLUDED.def,
	level = EXCLUDED.level,
	rank = EXCLUDED.rank,
	link_rating = EXCLUDED.link_rating,
	archetype = EXCLUDED.archetype,
	mentions = EXCLUDED.mentions,
	image_url = EXCLUDED.image_url,
	ai_processing = EXCLUDED.ai_processing,
	raw_payload = EXCLUDED.raw_payload,
	updated_at = EXCLUDED.updated_at
`

	_, err = tx.Exec(
		ctx,
		query,
		card.ID,
		card.UpstreamSource,
		card.UpstreamID,
		card.Passcode,
		card.KonamiID,
		card.Name,
		card.NormalizedName,
		string(aliases),
		card.NormalizedAliases,
		card.Description,
		card.CardType,
		card.FrameType,
		card.Race,
		card.Attribute,
		card.MonsterCategories,
		card.SpellTrapType,
		card.ATK,
		card.DEF,
		card.Level,
		card.Rank,
		card.LinkRating,
		card.Archetype,
		card.Mentions,
		card.ImageURL,
		jsonOrDefault(card.AIProcessing, "{}"),
		jsonOrDefault(card.RawPayload, "{}"),
		card.ImportedAt,
		card.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("seed card %s: %w", card.ID, err)
	}
	return nil
}

func seedEffect(ctx context.Context, tx pgx.Tx, effect domain.CardEffect) error {
	selectorHash, err := selectorHash(effect)
	if err != nil {
		return err
	}

	const query = `
INSERT INTO card_effects (
	id,
	source_card_id,
	effect_code,
	extraction_version,
	source_text,
	condition_text,
	cost_text,
	action_text,
	restriction_text,
	selector_status,
	selector_json,
	selector_hash,
	ai_confidence,
	is_active
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NULLIF($11, '')::jsonb, $12, $13, $14
)
ON CONFLICT (id) DO UPDATE SET
	source_card_id = EXCLUDED.source_card_id,
	effect_code = EXCLUDED.effect_code,
	extraction_version = EXCLUDED.extraction_version,
	source_text = EXCLUDED.source_text,
	condition_text = EXCLUDED.condition_text,
	cost_text = EXCLUDED.cost_text,
	action_text = EXCLUDED.action_text,
	restriction_text = EXCLUDED.restriction_text,
	selector_status = EXCLUDED.selector_status,
	selector_json = EXCLUDED.selector_json,
	selector_hash = EXCLUDED.selector_hash,
	ai_confidence = EXCLUDED.ai_confidence,
	is_active = EXCLUDED.is_active,
	updated_at = now()
`

	_, err = tx.Exec(
		ctx,
		query,
		effect.ID,
		effect.SourceCardID,
		effect.EffectCode,
		effect.ExtractionVersion,
		effect.SourceText,
		effect.ConditionText,
		effect.CostText,
		effect.ActionText,
		effect.RestrictionText,
		effect.SelectorStatus,
		nullableJSON(effect.SelectorJSON),
		selectorHash,
		effect.AIConfidence,
		effect.IsActive,
	)
	if err != nil {
		return fmt.Errorf("seed effect %s: %w", effect.ID, err)
	}
	return nil
}

func readCardsFixture(dir string) ([]domain.Card, error) {
	var fixture cardsFixture
	if err := readFixture(filepath.Join(dir, "cards.json"), &fixture); err != nil {
		return nil, err
	}
	return fixture.Cards, nil
}

func readEffectsFixture(dir string) ([]domain.CardEffect, error) {
	var fixture effectsFixture
	if err := readFixture(filepath.Join(dir, "card_effects.json"), &fixture); err != nil {
		return nil, err
	}
	return fixture.Effects, nil
}

func readFixture(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read fixture %s: %w", path, err)
	}
	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("decode fixture %s: %w", path, err)
	}
	return nil
}

func selectorHash(effect domain.CardEffect) (string, error) {
	if effect.SelectorStatus == "unresolved" {
		return "unresolved", nil
	}

	var selector any
	if err := json.Unmarshal(effect.SelectorJSON, &selector); err != nil {
		return "", fmt.Errorf("decode selector for effect %s: %w", effect.ID, err)
	}
	canonical, err := json.Marshal(selector)
	if err != nil {
		return "", fmt.Errorf("canonicalize selector for effect %s: %w", effect.ID, err)
	}

	sum := sha256.Sum256(canonical)
	return hex.EncodeToString(sum[:]), nil
}

func jsonOrDefault(value json.RawMessage, fallback string) string {
	if len(value) == 0 || string(value) == "null" {
		return fallback
	}
	return string(value)
}

func nullableJSON(value json.RawMessage) string {
	if len(value) == 0 || string(value) == "null" {
		return ""
	}
	return string(value)
}
