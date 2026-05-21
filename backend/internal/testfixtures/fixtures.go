package testfixtures

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"who-can-search-ygo/backend/internal/domain"
)

type CardsFixture struct {
	Cards []domain.Card `json:"cards"`
}

type EffectsFixture struct {
	Effects []domain.CardEffect `json:"card_effects"`
}

func LoadCards(t *testing.T) []domain.Card {
	t.Helper()

	var fixture CardsFixture
	readFixture(t, "cards.json", &fixture)
	return fixture.Cards
}

func LoadEffects(t *testing.T, cards []domain.Card) []domain.CardEffect {
	t.Helper()

	var fixture EffectsFixture
	readFixture(t, "card_effects.json", &fixture)

	byID := make(map[string]domain.Card, len(cards))
	for _, card := range cards {
		byID[card.ID] = card
	}

	for i := range fixture.Effects {
		source := byID[fixture.Effects[i].SourceCardID]
		fixture.Effects[i].SourceCard = domain.CardSummary{
			ID:       source.ID,
			Name:     source.Name,
			ImageURL: source.ImageURL,
		}
	}
	return fixture.Effects
}

func CardByName(t *testing.T, cards []domain.Card, name string) domain.Card {
	t.Helper()

	for _, card := range cards {
		if card.Name == name {
			return card
		}
	}
	t.Fatalf("card %q not found", name)
	return domain.Card{}
}

func readFixture(t *testing.T, name string, value any) {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate fixture helper")
	}
	path := filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata", "fixtures", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	if err := json.Unmarshal(data, value); err != nil {
		t.Fatalf("decode fixture %s: %v", name, err)
	}
}
