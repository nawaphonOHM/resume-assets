package mongodb

import (
	"fmt"
	"slices"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
)

const (
	mongoExpressionEqual       = "$eq"
	mongoExpressionConditional = "$cond"
)

// BuildFilter builds an exact filename match and optional exact tags expression.
func BuildFilter(filename string, criteria *domain.AssetCriteria) (mongodb.D, error) {
	filter := mongodb.D{{Key: "filename", Value: filename}}
	if criteria == nil {
		return filter, nil
	}
	expression, err := criteriaObjectExpression("$tags", map[string]any(*criteria))
	if err != nil {
		return nil, err
	}
	return append(filter, mongodb.E{Key: "$expr", Value: expression}), nil
}

func criteriaObjectExpression(actual any, expected map[string]any) (any, error) {
	comparisons, err := objectFieldComparisons("$$assetPairs", expected)
	if err != nil {
		return nil, err
	}
	objectMatches := mongodb.D{{Key: "$let", Value: mongodb.D{
		{Key: "vars", Value: mongodb.D{{Key: "assetPairs", Value: mongodb.D{{Key: "$objectToArray", Value: actual}}}}},
		{Key: "in", Value: mongodb.D{{Key: "$and", Value: comparisons}}},
	}}}
	return mongodb.D{{Key: mongoExpressionConditional, Value: mongodb.A{
		objectTypeCheck(actual),
		objectMatches,
		false,
	}}}, nil
}

func objectTypeCheck(actual any) mongodb.D {
	return mongodb.D{{Key: mongoExpressionEqual, Value: mongodb.A{
		mongodb.D{{Key: "$type", Value: actual}},
		"object",
	}}}
}

func sortedObjectKeys(object map[string]any) []string {
	keys := make([]string, 0, len(object))
	for key := range object {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

func literalExpression(value any) mongodb.D {
	return mongodb.D{{Key: "$literal", Value: value}}
}

func invalidObjectValue(key string, err error) error {
	return fmt.Errorf("criteria field %q: %w", key, err)
}
