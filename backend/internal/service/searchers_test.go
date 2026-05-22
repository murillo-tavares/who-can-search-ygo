package service

import (
	"context"
	"strings"
	"testing"

	"who-can-search-ygo/backend/internal/domain"
	"who-can-search-ygo/backend/internal/testfixtures"
)

func TestGetAddDeckToHandSearchers(t *testing.T) {
	cards := testfixtures.LoadCards(t)
	effects := testfixtures.LoadEffects(t, cards)
	target := testfixtures.CardByName(t, cards, "Dark Magician")

	svc := New(newFixtureRepo(cards, effects))
	result, err := svc.GetAddDeckToHandSearchers(context.Background(), target.ID)
	if err != nil {
		t.Fatalf("GetAddDeckToHandSearchers returned error: %v", err)
	}

	if result.TargetCard.Name != "Dark Magician" {
		t.Fatalf("target card = %q, want Dark Magician", result.TargetCard.Name)
	}
	if result.EffectCode != domain.EffectCodeAddDeckToHand {
		t.Fatalf("effect code = %q, want %q", result.EffectCode, domain.EffectCodeAddDeckToHand)
	}

	names := make([]string, 0, len(result.Results))
	for _, item := range result.Results {
		names = append(names, item.SourceCard.Name)
	}
	assertContains(t, names, "Dark Magical Circle")
	assertNotContains(t, names, "Illusion of Chaos")
}

type fixtureRepo struct {
	cards   []domain.Card
	effects []domain.CardEffect
}

func newFixtureRepo(cards []domain.Card, effects []domain.CardEffect) *fixtureRepo {
	return &fixtureRepo{cards: cards, effects: effects}
}

func (r *fixtureRepo) GetCard(_ context.Context, id string) (domain.Card, error) {
	for _, card := range r.cards {
		if card.ID == id {
			return card, nil
		}
	}
	return domain.Card{}, ErrNotFound
}

func (r *fixtureRepo) SearchCards(_ context.Context, query string, limit int) ([]domain.CardSummary, error) {
	normalized := strings.ToLower(strings.TrimSpace(query))
	var results []domain.CardSummary
	for _, card := range r.cards {
		if strings.Contains(strings.ToLower(card.Name), normalized) {
			results = append(results, domain.CardSummary{ID: card.ID, Name: card.Name, ImageURL: card.ImageURL})
		}
		if len(results) == limit {
			break
		}
	}
	return results, nil
}

func (r *fixtureRepo) ListActiveResolvedEffects(_ context.Context, effectCode string) ([]domain.CardEffect, error) {
	var results []domain.CardEffect
	for _, effect := range r.effects {
		if effect.EffectCode == effectCode && effect.SelectorStatus == domain.SelectorStatusResolved && effect.IsActive {
			results = append(results, effect)
		}
	}
	return results, nil
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
