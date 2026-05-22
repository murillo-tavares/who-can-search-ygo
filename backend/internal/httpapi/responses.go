package httpapi

import "who-can-search-ygo/backend/internal/domain"

type cardResponse struct {
	ID                string         `json:"id"`
	UpstreamSource    string         `json:"upstream_source,omitempty"`
	UpstreamID        string         `json:"upstream_id,omitempty"`
	Passcode          *string        `json:"passcode,omitempty"`
	KonamiID          *string        `json:"konami_id,omitempty"`
	Name              string         `json:"name"`
	NormalizedName    string         `json:"normalized_name"`
	Aliases           []domain.Alias `json:"aliases"`
	NormalizedAliases []string       `json:"normalized_aliases"`
	Description       string         `json:"description,omitempty"`
	CardType          *string        `json:"card_type"`
	FrameType         *string        `json:"frame_type,omitempty"`
	Race              *string        `json:"race"`
	Attribute         *string        `json:"attribute"`
	MonsterCategories []string       `json:"monster_categories"`
	SpellTrapType     *string        `json:"spell_trap_type"`
	ATK               *int           `json:"atk"`
	DEF               *int           `json:"def"`
	Level             *int           `json:"level"`
	Rank              *int           `json:"rank"`
	LinkRating        *int           `json:"link_rating"`
	Archetype         *string        `json:"archetype"`
	Mentions          []string       `json:"mentions"`
	TextFeatures      []string       `json:"text_features"`
	ImageURL          *string        `json:"image_url,omitempty"`
}

func newCardResponse(card domain.Card) cardResponse {
	return cardResponse{
		ID:                card.ID,
		UpstreamSource:    card.UpstreamSource,
		UpstreamID:        card.UpstreamID,
		Passcode:          card.Passcode,
		KonamiID:          card.KonamiID,
		Name:              card.Name,
		NormalizedName:    card.NormalizedName,
		Aliases:           card.Aliases,
		NormalizedAliases: card.NormalizedAliases,
		Description:       card.Description,
		CardType:          card.CardType,
		FrameType:         card.FrameType,
		Race:              card.Race,
		Attribute:         card.Attribute,
		MonsterCategories: card.MonsterCategories,
		SpellTrapType:     card.SpellTrapType,
		ATK:               card.ATK,
		DEF:               card.DEF,
		Level:             card.Level,
		Rank:              card.Rank,
		LinkRating:        card.LinkRating,
		Archetype:         card.Archetype,
		Mentions:          card.Mentions,
		TextFeatures:      card.TextFeatures,
		ImageURL:          card.ImageURL,
	}
}
