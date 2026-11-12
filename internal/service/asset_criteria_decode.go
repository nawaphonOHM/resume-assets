package service

import (
	"encoding/json"
	"errors"
	"fmt"
)

func decodeCriteriaObject(decoder *json.Decoder) (map[string]any, error) {
	object, err := decodeCriteriaObjectFields(decoder)
	if err != nil {
		return nil, err
	}
	if err := consumeCriteriaDelimiter(decoder, '}'); err != nil {
		return nil, err
	}
	return object, nil
}

func decodeCriteriaObjectFields(decoder *json.Decoder) (map[string]any, error) {
	object := make(map[string]any)
	for decoder.More() {
		key, value, err := decodeCriteriaObjectEntry(decoder, object)
		if err != nil {
			return nil, err
		}
		object[key] = value
	}
	return object, nil
}

func decodeCriteriaObjectEntry(decoder *json.Decoder, object map[string]any) (string, any, error) {
	key, err := decodeCriteriaKey(decoder, object)
	if err != nil {
		return "", nil, err
	}
	value, err := decodeCriteriaValue(decoder)
	return key, value, err
}

func decodeCriteriaKey(decoder *json.Decoder, object map[string]any) (string, error) {
	token, err := decoder.Token()
	if err != nil {
		return "", err
	}
	return validateCriteriaKey(token, object)
}

func validateCriteriaKey(token any, object map[string]any) (string, error) {
	key, ok := token.(string)
	if !ok {
		return "", errors.New("object key is not a string")
	}
	if err := ensureUniqueCriteriaKey(key, object); err != nil {
		return "", err
	}
	return key, nil
}

func ensureUniqueCriteriaKey(key string, object map[string]any) error {
	if _, exists := object[key]; exists {
		return fmt.Errorf("duplicate object key %q", key)
	}
	return nil
}

func decodeCriteriaArray(decoder *json.Decoder) ([]any, error) {
	array, err := decodeCriteriaArrayItems(decoder)
	if err != nil {
		return nil, err
	}
	if err := consumeCriteriaDelimiter(decoder, ']'); err != nil {
		return nil, err
	}
	return array, nil
}

func decodeCriteriaArrayItems(decoder *json.Decoder) ([]any, error) {
	array := make([]any, 0)
	for decoder.More() {
		value, err := decodeCriteriaValue(decoder)
		if err != nil {
			return nil, err
		}
		array = append(array, value)
	}
	return array, nil
}

func consumeCriteriaDelimiter(decoder *json.Decoder, want rune) error {
	closing, err := decoder.Token()
	if err != nil {
		return err
	}
	if closing != json.Delim(want) {
		return fmt.Errorf("expected closing delimiter %q", want)
	}
	return nil
}
