package mongodb

import (
	"fmt"

	"github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
)

func objectFieldComparisons(actualPairs any, expected map[string]any) (mongodb.A, error) {
	comparisons := mongodb.A{objectFieldCountCheck(actualPairs, len(expected))}
	for _, key := range sortedObjectKeys(expected) {
		comparison, err := exactObjectFieldComparison(actualPairs, key, expected[key])
		if err != nil {
			return nil, invalidObjectValue(key, err)
		}
		comparisons = append(comparisons, comparison)
	}
	return comparisons, nil
}

func objectFieldCountCheck(actualPairs any, expectedCount int) mongodb.D {
	return mongodb.D{{Key: "$eq", Value: mongodb.A{
		mongodb.D{{Key: "$size", Value: actualPairs}},
		expectedCount,
	}}}
}

func exactObjectFieldComparison(actualPairs any, key string, expected any) (any, error) {
	valueExpression, err := criteriaValueExpression("$$assetValue", expected)
	if err != nil {
		return nil, fmt.Errorf("field value: %w", err)
	}
	return mongodb.D{{Key: "$let", Value: mongodb.D{
		{Key: "vars", Value: mongodb.D{{Key: "assetMatches", Value: matchingObjectFields(actualPairs, key)}}},
		{Key: "in", Value: mongodb.D{{Key: mongoExpressionConditional, Value: mongodb.A{
			objectFieldCountCheck("$$assetMatches", 1),
			matchingFieldValue(valueExpression),
			false,
		}}}},
	}}}, nil
}

func matchingObjectFields(actualPairs any, key string) mongodb.D {
	return mongodb.D{{Key: "$filter", Value: mongodb.D{
		{Key: "input", Value: actualPairs},
		{Key: "as", Value: "assetPair"},
		{Key: "cond", Value: mongodb.D{{Key: "$eq", Value: mongodb.A{
			"$$assetPair.k",
			literalExpression(key),
		}}}},
	}}}
}

func matchingFieldValue(valueExpression any) mongodb.D {
	return mongodb.D{{Key: "$let", Value: mongodb.D{
		{Key: "vars", Value: mongodb.D{{Key: "assetValue", Value: arrayItemAt("$$assetMatches.v", 0)}}},
		{Key: "in", Value: valueExpression},
	}}}
}
