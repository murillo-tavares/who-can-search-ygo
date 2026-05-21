package rules

import (
	"testing"

	"who-can-search-ygo/backend/internal/domain"
	"who-can-search-ygo/backend/internal/testfixtures"
)

func TestMatchSelectorForDarkMagicianFixture(t *testing.T) {
	cards := testfixtures.LoadCards(t)
	effects := testfixtures.LoadEffects(t, cards)
	target := testfixtures.CardByName(t, cards, "Dark Magician")

	matches := matchingSourceNames(effects, target)

	assertContains(t, matches, "Dark Magical Circle")
	assertNotContains(t, matches, "Illusion of Chaos")
}

func TestMatchSelectorForWarriorMonsterFixture(t *testing.T) {
	cards := testfixtures.LoadCards(t)
	effects := testfixtures.LoadEffects(t, cards)
	target := testfixtures.CardByName(t, cards, "Elemental HERO Stratos")

	matches := matchingSourceNames(effects, target)

	assertContains(t, matches, "Reinforcement of the Army")
	assertContains(t, matches, "E - Emergency Call")
}

func matchingSourceNames(effects []domain.CardEffect, target domain.Card) []string {
	var matches []string
	for _, effect := range effects {
		if effect.EffectCode != domain.EffectCodeAddDeckToHand || effect.SelectorStatus != "resolved" || !effect.IsActive {
			continue
		}
		if MatchSelector(effect.SelectorJSON, target, ZoneDeck) {
			matches = append(matches, effect.SourceCard.Name)
		}
	}
	return matches
}

func assertContains(t *testing.T, values []string, expected string) {
	t.Helper()

	for _, value := range values {
		if value == expected {
			return
		}
	}
	t.Fatalf("expected %q in %v", expected, values)
}

func assertNotContains(t *testing.T, values []string, unexpected string) {
	t.Helper()

	for _, value := range values {
		if value == unexpected {
			t.Fatalf("did not expect %q in %v", unexpected, values)
		}
	}
}
