package domain

import "encoding/json"

const EffectCodeAddDeckToHand = "add_deck_to_hand"

const (
	SelectorStatusResolved   = "resolved"
	SelectorStatusIgnored    = "ignored"
	SelectorStatusUnresolved = "unresolved"
)

type Card struct {
	ID                string          `json:"id"`
	UpstreamSource    string          `json:"upstream_source,omitempty"`
	UpstreamID        string          `json:"upstream_id,omitempty"`
	Passcode          *string         `json:"passcode,omitempty"`
	KonamiID          *string         `json:"konami_id,omitempty"`
	Name              string          `json:"name"`
	NormalizedName    string          `json:"normalized_name"`
	Aliases           []Alias         `json:"aliases"`
	NormalizedAliases []string        `json:"normalized_aliases"`
	Description       string          `json:"description,omitempty"`
	CardType          *string         `json:"card_type"`
	FrameType         *string         `json:"frame_type,omitempty"`
	Race              *string         `json:"race"`
	Attribute         *string         `json:"attribute"`
	MonsterCategories []string        `json:"monster_categories"`
	SpellTrapType     *string         `json:"spell_trap_type"`
	ATK               *int            `json:"atk"`
	DEF               *int            `json:"def"`
	Level             *int            `json:"level"`
	Rank              *int            `json:"rank"`
	LinkRating        *int            `json:"link_rating"`
	Archetype         *string         `json:"archetype"`
	Mentions          []string        `json:"mentions"`
	TextFeatures      []string        `json:"text_features"`
	ImageURL          *string         `json:"image_url,omitempty"`
	AIProcessing      json.RawMessage `json:"ai_processing,omitempty"`
	RawPayload        json.RawMessage `json:"raw_payload,omitempty"`
	ImportedAt        string          `json:"imported_at,omitempty"`
	UpdatedAt         string          `json:"updated_at,omitempty"`
}

type Alias struct {
	Alias              string          `json:"alias"`
	NormalizedAlias    string          `json:"normalized_alias"`
	AliasKind          string          `json:"alias_kind"`
	AppliesInZoneCodes []string        `json:"applies_in_zone_codes"`
	ConditionJSON      json.RawMessage `json:"condition_json"`
	Source             string          `json:"source"`
}

type CardSummary struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	ImageURL *string `json:"image_url"`
}

type CardEffect struct {
	ID                string          `json:"id"`
	SourceCardID      string          `json:"source_card_id"`
	SourceCardName    string          `json:"source_card_name,omitempty"`
	EffectCode        string          `json:"effect_code"`
	ExtractionVersion string          `json:"extraction_version,omitempty"`
	SourceText        string          `json:"source_text"`
	ConditionText     *string         `json:"condition_text"`
	CostText          *string         `json:"cost_text"`
	ActionText        string          `json:"action_text"`
	RestrictionText   *string         `json:"restriction_text"`
	SelectorStatus    string          `json:"selector_status"`
	SelectorJSON      json.RawMessage `json:"selector_json"`
	AIConfidence      *float64        `json:"ai_confidence"`
	IsActive          bool            `json:"is_active"`
	SourceCard        CardSummary     `json:"source_card"`
}
