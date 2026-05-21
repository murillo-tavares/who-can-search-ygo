package service

import (
	"context"
	"errors"

	"who-can-search-ygo/backend/internal/domain"
	"who-can-search-ygo/backend/internal/rules"
)

var ErrNotFound = errors.New("not found")

type Repository interface {
	GetCard(ctx context.Context, id string) (domain.Card, error)
	SearchCards(ctx context.Context, query string, limit int) ([]domain.CardSummary, error)
	ListActiveResolvedEffects(ctx context.Context, effectCode string) ([]domain.CardEffect, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

type SearchersResult struct {
	TargetCard TargetCardSummary `json:"target_card"`
	EffectCode string            `json:"effect_code"`
	Results    []SearcherResult  `json:"results"`
}

type TargetCardSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SearcherResult struct {
	EffectID        string             `json:"effect_id"`
	SourceCard      domain.CardSummary `json:"source_card"`
	SourceText      string             `json:"source_text"`
	ConditionText   *string            `json:"condition_text"`
	CostText        *string            `json:"cost_text"`
	ActionText      string             `json:"action_text"`
	RestrictionText *string            `json:"restriction_text"`
}

func (s *Service) GetCard(ctx context.Context, id string) (domain.Card, error) {
	return s.repo.GetCard(ctx, id)
}

func (s *Service) SearchCards(ctx context.Context, query string, limit int) ([]domain.CardSummary, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	return s.repo.SearchCards(ctx, query, limit)
}

func (s *Service) GetAddDeckToHandSearchers(ctx context.Context, targetID string) (SearchersResult, error) {
	target, err := s.repo.GetCard(ctx, targetID)
	if err != nil {
		return SearchersResult{}, err
	}

	effects, err := s.repo.ListActiveResolvedEffects(ctx, domain.EffectCodeAddDeckToHand)
	if err != nil {
		return SearchersResult{}, err
	}

	result := SearchersResult{
		TargetCard: TargetCardSummary{ID: target.ID, Name: target.Name},
		EffectCode: domain.EffectCodeAddDeckToHand,
		Results:    []SearcherResult{},
	}

	for _, effect := range effects {
		if !rules.MatchSelector(effect.SelectorJSON, target, rules.ZoneDeck) {
			continue
		}
		result.Results = append(result.Results, SearcherResult{
			EffectID:        effect.ID,
			SourceCard:      effect.SourceCard,
			SourceText:      effect.SourceText,
			ConditionText:   effect.ConditionText,
			CostText:        effect.CostText,
			ActionText:      effect.ActionText,
			RestrictionText: effect.RestrictionText,
		})
	}

	return result, nil
}
