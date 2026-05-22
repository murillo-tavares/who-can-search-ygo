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

func TestMatchSelectorForSupportedCardFields(t *testing.T) {
	atk := 1200
	def := 800
	level := 4
	cardType := "Monster"
	target := domain.Card{
		Name:              "Feature Target",
		CardType:          &cardType,
		ATK:               &atk,
		DEF:               &def,
		Level:             &level,
		MonsterCategories: []string{"Ritual", "Effect"},
		TextFeatures:      []string{"coin_toss_effect"},
	}

	tests := []struct {
		name     string
		selector string
	}{
		{
			name:     "monster category contains",
			selector: `{"type":"comparison","field":"monster_categories","op":"contains","value":"Ritual"}`,
		},
		{
			name:     "text feature contains",
			selector: `{"type":"comparison","field":"text_features","op":"contains","value":"coin_toss_effect"}`,
		},
		{
			name:     "combined atk and def",
			selector: `{"type":"comparison","field":"combined_atk_def","op":"=","value":2000}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !MatchSelector([]byte(test.selector), target, ZoneDeck) {
				t.Fatalf("selector did not match target")
			}
		})
	}
}

func matchingSourceNames(effects []domain.CardEffect, target domain.Card) []string {
	var matches []string
	for _, effect := range effects {
		if effect.EffectCode != domain.EffectCodeAddDeckToHand || effect.SelectorStatus != domain.SelectorStatusResolved || !effect.IsActive {
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
