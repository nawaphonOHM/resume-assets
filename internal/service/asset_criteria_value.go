package service

import (
	"encoding/json"
	"fmt"
)

func decodeCriteriaValue(decoder *json.Decoder) (any, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	return decodeCriteriaToken(token, decoder)
}

func decodeCriteriaToken(token any, decoder *json.Decoder) (any, error) {
	delimiter, ok := token.(json.Delim)
	if !ok {
		return token, nil
	}
	return decodeCriteriaCompound(decoder, delimiter)
}

func decodeCriteriaCompound(decoder *json.Decoder, delimiter json.Delim) (any, error) {
	if delimiter == '{' {
		return decodeCriteriaObject(decoder)
	}
	if delimiter == '[' {
		return decodeCriteriaArray(decoder)
	}
	return nil, fmt.Errorf("unexpected delimiter %q", delimiter)
}
