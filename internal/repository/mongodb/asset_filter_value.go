package mongodb

import (
	"encoding/json"
	"fmt"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
)

func criteriaValueExpression(actual, expected any) (any, error) {
	if object, ok := criteriaObjectValue(expected); ok {
		return criteriaObjectExpression(actual, object)
	}
	if array, ok := expected.([]any); ok {
		return criteriaArrayExpression(actual, array)
	}
	return criteriaScalarExpression(actual, expected)
}

func criteriaObjectValue(value any) (map[string]any, bool) {
	if object, ok := value.(map[string]any); ok {
		return object, true
	}
	criteria, ok := value.(domain.AssetCriteria)
	return map[string]any(criteria), ok
}

func criteriaScalarExpression(actual, expected any) (any, error) {
	operand, err := criteriaScalarOperand(expected)
	if err != nil {
		return nil, err
	}
	return mongodb.D{{Key: mongoExpressionEqual, Value: mongodb.A{actual, operand}}}, nil
}

func criteriaScalarOperand(value any) (any, error) {
	if number, ok := value.(json.Number); ok {
		return decimalCriteriaExpression(number)
	}
	if isSupportedCriteriaScalar(value) {
		return literalExpression(value), nil
	}
	return nil, fmt.Errorf("unsupported criteria value type %T", value)
}

func isSupportedCriteriaScalar(value any) bool {
	switch value.(type) {
	case nil, string, bool:
		return true
	default:
		return false
	}
}

func decimalCriteriaExpression(number json.Number) (any, error) {
	if err := domain.ValidateAssetCriteriaNumber(number.String()); err != nil {
		return nil, fmt.Errorf("invalid BSON numeric criteria: %w", err)
	}
	return mongodb.D{{Key: "$toDecimal", Value: literalExpression(number.String())}}, nil
}
