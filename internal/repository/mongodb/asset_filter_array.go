package mongodb

import (
	"fmt"

	"github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
)

func criteriaArrayExpression(actual any, expected []any) (any, error) {
	itemExpressions, err := criteriaArrayItemExpressions(actual, expected)
	if err != nil {
		return nil, err
	}
	elementsMatch := any(true)
	if len(itemExpressions) > 0 {
		elementsMatch = mongodb.D{{Key: "$and", Value: itemExpressions}}
	}
	return mongodb.D{{Key: mongoExpressionConditional, Value: mongodb.A{
		arrayTypeCheck(actual),
		mongodb.D{{Key: mongoExpressionConditional, Value: mongodb.A{
			objectFieldCountCheck(actual, len(expected)),
			elementsMatch,
			false,
		}}},
		false,
	}}}, nil
}

func criteriaArrayItemExpressions(actual any, expected []any) (mongodb.A, error) {
	items := make(mongodb.A, 0, len(expected))
	for index, expectedItem := range expected {
		item, err := criteriaValueExpression(arrayItemAt(actual, index), expectedItem)
		if err != nil {
			return nil, fmt.Errorf("array index %d: %w", index, err)
		}
		items = append(items, item)
	}
	return items, nil
}

func arrayItemAt(actual any, index int) mongodb.D {
	return mongodb.D{{Key: "$arrayElemAt", Value: mongodb.A{actual, index}}}
}

func arrayTypeCheck(actual any) mongodb.D {
	return mongodb.D{{Key: "$eq", Value: mongodb.A{
		mongodb.D{{Key: "$type", Value: actual}},
		"array",
	}}}
}
