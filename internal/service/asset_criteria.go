package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

func parseAssetCriteria(raw string) (domain.AssetCriteria, error) {
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.UseNumber()
	return decodeAssetCriteria(decoder)
}

func parseOptionalAssetCriteria(criteriaJSON *string) (*domain.AssetCriteria, error) {
	if criteriaJSON == nil {
		return nil, nil
	}
	criteria, err := parseAssetCriteria(*criteriaJSON)
	if err != nil {
		return nil, err
	}
	return &criteria, nil
}

func decodeAssetCriteria(decoder *json.Decoder) (domain.AssetCriteria, error) {
	object, err := decodeCriteriaRoot(decoder)
	if err != nil {
		return nil, invalidAssetCriteria(err)
	}
	if err := validateCriteriaNumbers(object); err != nil {
		return nil, invalidAssetCriteria(err)
	}
	return domain.AssetCriteria(object), nil
}

func decodeCriteriaRoot(decoder *json.Decoder) (map[string]any, error) {
	value, err := decodeCriteriaValue(decoder)
	if err != nil {
		return nil, err
	}
	return checkCriteriaRoot(decoder, value)
}

func checkCriteriaRoot(decoder *json.Decoder, value any) (map[string]any, error) {
	object, ok := value.(map[string]any)
	if !ok {
		return nil, errors.New("criteria root must be an object")
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return nil, errors.New("criteria contains trailing data")
	}
	return object, nil
}

func validateCriteriaNumbers(value any) error {
	if number, ok := value.(json.Number); ok {
		return validateCriteriaNumber(number)
	}
	return validateCriteriaContainer(value)
}

func validateCriteriaNumber(number json.Number) error {
	if err := domain.ValidateAssetCriteriaNumber(number.String()); err != nil {
		return fmt.Errorf("number %q cannot be represented faithfully as BSON Decimal128: %w", number, err)
	}
	return nil
}

func validateCriteriaContainer(value any) error {
	switch typed := value.(type) {
	case map[string]any:
		return validateCriteriaObjectNumbers(typed)
	case []any:
		return validateCriteriaArrayNumbers(typed)
	}
	return nil
}

func validateCriteriaObjectNumbers(object map[string]any) error {
	for _, child := range object {
		if err := validateCriteriaNumbers(child); err != nil {
			return err
		}
	}
	return nil
}

func validateCriteriaArrayNumbers(array []any) error {
	for _, child := range array {
		if err := validateCriteriaNumbers(child); err != nil {
			return err
		}
	}
	return nil
}

func invalidAssetCriteria(err error) error {
	if err == nil {
		return domain.ErrInvalidAssetCriteria
	}
	return fmt.Errorf("%w: %v", domain.ErrInvalidAssetCriteria, err)
}
