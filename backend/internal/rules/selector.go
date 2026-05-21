package rules

import (
	"encoding/json"
	"fmt"
	"strings"

	"who-can-search-ygo/backend/internal/domain"
)

const ZoneDeck = "deck"

type Expression interface {
	Match(card domain.Card, zone string) bool
}

type rawExpression struct {
	Type  string            `json:"type"`
	Op    string            `json:"op"`
	Field string            `json:"field"`
	Value json.RawMessage   `json:"value"`
	Args  []json.RawMessage `json:"args"`
}

type Comparison struct {
	Field string
	Op    string
	Value any
}

type Logical struct {
	Op   string
	Args []Expression
}

func ParseSelector(data json.RawMessage) (Expression, error) {
	if len(data) == 0 || string(data) == "null" {
		return nil, fmt.Errorf("selector is empty")
	}

	var raw rawExpression
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("decode selector: %w", err)
	}

	switch raw.Type {
	case "comparison":
		value, err := decodeValue(raw.Value)
		if err != nil {
			return nil, err
		}
		if !supportedField(raw.Field) || !supportedComparisonOp(raw.Op) {
			return nil, fmt.Errorf("unsupported comparison")
		}
		return Comparison{Field: raw.Field, Op: raw.Op, Value: value}, nil
	case "logical":
		if raw.Op != "and" && raw.Op != "or" {
			return nil, fmt.Errorf("unsupported logical operator")
		}
		if len(raw.Args) == 0 {
			return nil, fmt.Errorf("logical selector has no args")
		}
		args := make([]Expression, 0, len(raw.Args))
		for _, arg := range raw.Args {
			expr, err := ParseSelector(arg)
			if err != nil {
				return nil, err
			}
			args = append(args, expr)
		}
		return Logical{Op: raw.Op, Args: args}, nil
	default:
		return nil, fmt.Errorf("unsupported selector type")
	}
}

func MatchSelector(data json.RawMessage, card domain.Card, zone string) bool {
	expr, err := ParseSelector(data)
	if err != nil {
		return false
	}
	return expr.Match(card, zone)
}

func (l Logical) Match(card domain.Card, zone string) bool {
	switch l.Op {
	case "and":
		for _, arg := range l.Args {
			if !arg.Match(card, zone) {
				return false
			}
		}
		return true
	case "or":
		for _, arg := range l.Args {
			if arg.Match(card, zone) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func (c Comparison) Match(card domain.Card, zone string) bool {
	fieldValue, ok := cardFieldValue(card, c.Field, zone)
	if !ok {
		return false
	}

	switch c.Op {
	case "=":
		return equal(fieldValue, c.Value)
	case "!=":
		return !equal(fieldValue, c.Value)
	case "in":
		return in(fieldValue, c.Value)
	case "not_in":
		return !in(fieldValue, c.Value)
	case "contains":
		return contains(fieldValue, c.Value)
	case "not_contains":
		return !contains(fieldValue, c.Value)
	case "<", "<=", ">", ">=":
		return compareNumber(fieldValue, c.Value, c.Op)
	default:
		return false
	}
}

func decodeValue(data json.RawMessage) (any, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("comparison value is empty")
	}

	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("decode comparison value: %w", err)
	}
	return v, nil
}

func supportedField(field string) bool {
	switch field {
	case "name", "card_type", "race", "attribute", "level", "rank", "link_rating", "atk", "def", "archetype", "mentions", "spell_trap_type":
		return true
	default:
		return false
	}
}

func supportedComparisonOp(op string) bool {
	switch op {
	case "=", "!=", "in", "not_in", "contains", "not_contains", "<", "<=", ">", ">=":
		return true
	default:
		return false
	}
}

func cardFieldValue(card domain.Card, field string, zone string) (any, bool) {
	switch field {
	case "name":
		names := []string{card.Name}
		for _, alias := range card.Aliases {
			if alias.AliasKind == "exact_name" && aliasApplies(alias, zone) {
				names = append(names, alias.Alias)
			}
		}
		return names, true
	case "card_type":
		return pointerString(card.CardType)
	case "race":
		return pointerString(card.Race)
	case "attribute":
		return pointerString(card.Attribute)
	case "level":
		return pointerInt(card.Level)
	case "rank":
		return pointerInt(card.Rank)
	case "link_rating":
		return pointerInt(card.LinkRating)
	case "atk":
		return pointerInt(card.ATK)
	case "def":
		return pointerInt(card.DEF)
	case "archetype":
		values := make([]string, 0, 1+len(card.Aliases))
		if card.Archetype != nil {
			values = append(values, *card.Archetype)
		}
		for _, alias := range card.Aliases {
			if alias.AliasKind == "archetype_membership" && aliasApplies(alias, zone) {
				values = append(values, alias.Alias)
			}
		}
		return values, true
	case "mentions":
		return card.Mentions, true
	case "spell_trap_type":
		return pointerString(card.SpellTrapType)
	default:
		return nil, false
	}
}

func pointerString(value *string) (any, bool) {
	if value == nil {
		return nil, false
	}
	return *value, true
}

func pointerInt(value *int) (any, bool) {
	if value == nil {
		return nil, false
	}
	return *value, true
}

func aliasApplies(alias domain.Alias, zone string) bool {
	if len(alias.AppliesInZoneCodes) > 0 && !stringSliceContains(alias.AppliesInZoneCodes, zone) {
		return false
	}

	if len(alias.ConditionJSON) == 0 {
		return false
	}
	var condition struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(alias.ConditionJSON, &condition); err != nil {
		return false
	}
	return condition.Type == "always"
}

func equal(left any, right any) bool {
	switch l := left.(type) {
	case string:
		r, ok := right.(string)
		return ok && strings.EqualFold(l, r)
	case []string:
		r, ok := right.(string)
		if !ok {
			return false
		}
		for _, item := range l {
			if strings.EqualFold(item, r) {
				return true
			}
		}
		return false
	case int:
		r, ok := jsonNumberToInt(right)
		return ok && l == r
	default:
		return false
	}
}

func in(left any, right any) bool {
	values, ok := right.([]any)
	if !ok {
		return false
	}
	for _, value := range values {
		if equal(left, value) {
			return true
		}
	}
	return false
}

func contains(left any, right any) bool {
	needle, ok := right.(string)
	if !ok {
		return false
	}

	values, ok := left.([]string)
	if !ok {
		return false
	}
	return stringSliceContainsFold(values, needle)
}

func compareNumber(left any, right any, op string) bool {
	l, ok := left.(int)
	if !ok {
		return false
	}
	r, ok := jsonNumberToInt(right)
	if !ok {
		return false
	}

	switch op {
	case "<":
		return l < r
	case "<=":
		return l <= r
	case ">":
		return l > r
	case ">=":
		return l >= r
	default:
		return false
	}
}

func jsonNumberToInt(value any) (int, bool) {
	switch v := value.(type) {
	case float64:
		if v != float64(int(v)) {
			return 0, false
		}
		return int(v), true
	case int:
		return v, true
	default:
		return 0, false
	}
}

func stringSliceContains(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}

func stringSliceContainsFold(values []string, needle string) bool {
	for _, value := range values {
		if strings.EqualFold(value, needle) {
			return true
		}
	}
	return false
}
