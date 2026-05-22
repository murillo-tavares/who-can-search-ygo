package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"who-can-search-ygo/backend/internal/domain"
	"who-can-search-ygo/backend/internal/service"
	"who-can-search-ygo/backend/internal/testfixtures"
)

func TestGetCardSearchersEndpoint(t *testing.T) {
	cards := testfixtures.LoadCards(t)
	effects := testfixtures.LoadEffects(t, cards)
	target := testfixtures.CardByName(t, cards, "Dark Magician")

	handler := NewHandler(service.New(&fixtureRepo{cards: cards, effects: effects}))
	request := httptest.NewRequest(http.MethodGet, "/cards/"+target.ID+"/searchers", nil)
	response := httptest.NewRecorder()

	handler.Routes().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusOK, response.Body.String())
	}

	var body service.SearchersResult
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.TargetCard.Name != "Dark Magician" {
		t.Fatalf("target card = %q, want Dark Magician", body.TargetCard.Name)
	}

	names := make([]string, 0, len(body.Results))
	for _, result := range body.Results {
		names = append(names, result.SourceCard.Name)
	}
	assertContains(t, names, "Dark Magical Circle")
	assertNotContains(t, names, "Illusion of Chaos")
}

func TestDocumentationEndpoints(t *testing.T) {
	handler := NewHandler(service.New(&fixtureRepo{}))

	openAPIResponse := httptest.NewRecorder()
	openAPIRequest := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	handler.Routes().ServeHTTP(openAPIResponse, openAPIRequest)

	if openAPIResponse.Code != http.StatusOK {
		t.Fatalf("openapi status = %d, want %d", openAPIResponse.Code, http.StatusOK)
	}
	if contentType := openAPIResponse.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("openapi content type = %q, want application/json", contentType)
	}

	var spec map[string]any
	if err := json.Unmarshal(openAPIResponse.Body.Bytes(), &spec); err != nil {
		t.Fatalf("decode openapi spec: %v", err)
	}
	if spec["openapi"] != "3.1.0" {
		t.Fatalf("openapi version = %v, want 3.1.0", spec["openapi"])
	}

	docsResponse := httptest.NewRecorder()
	docsRequest := httptest.NewRequest(http.MethodGet, "/docs", nil)
	handler.Routes().ServeHTTP(docsResponse, docsRequest)

	if docsResponse.Code != http.StatusOK {
		t.Fatalf("docs status = %d, want %d", docsResponse.Code, http.StatusOK)
	}
	if !strings.Contains(docsResponse.Body.String(), "Scalar.createApiReference") {
		t.Fatal("docs response does not render Scalar")
	}
}

func TestGetCardDoesNotExposeInternalPayloads(t *testing.T) {
	cards := testfixtures.LoadCards(t)
	target := testfixtures.CardByName(t, cards, "Dark Magician")

	handler := NewHandler(service.New(&fixtureRepo{cards: cards}))
	request := httptest.NewRequest(http.MethodGet, "/cards/"+target.ID, nil)
	response := httptest.NewRecorder()

	handler.Routes().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusOK, response.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, ok := body["raw_payload"]; ok {
		t.Fatal("response exposes raw_payload")
	}
	if _, ok := body["ai_processing"]; ok {
		t.Fatal("response exposes ai_processing")
	}
}

func TestSearchCardsRejectsTooLargeLimit(t *testing.T) {
	handler := NewHandler(service.New(&fixtureRepo{}))
	request := httptest.NewRequest(http.MethodGet, "/cards?query=dark&limit=51", nil)
	response := httptest.NewRecorder()

	handler.Routes().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusBadRequest, response.Body.String())
	}
}

type fixtureRepo struct {
	cards   []domain.Card
	effects []domain.CardEffect
}

func (r *fixtureRepo) GetCard(_ context.Context, id string) (domain.Card, error) {
	for _, card := range r.cards {
		if card.ID == id {
			return card, nil
		}
	}
	return domain.Card{}, service.ErrNotFound
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
