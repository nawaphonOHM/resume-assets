package domain

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

const (
	decimal128ExponentBase = 10
	decimal128ExponentBits = 32
	decimal128MaxDigits    = 34
	decimal128MinScale     = -6176
	decimal128MaxScale     = 6111
)

var errExpectedJSONNumber = errors.New("expected one JSON number")

// ValidateAssetCriteriaNumber reports whether a JSON number fits Decimal128 without rounding.
func ValidateAssetCriteriaNumber(raw string) error {
	number, err := parseJSONNumber(raw)
	if err != nil {
		return err
	}
	return validateDecimal128Number(number.String())
}

func parseJSONNumber(raw string) (json.Number, error) {
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.UseNumber()
	value, err := decoder.Token()
	if err != nil {
		return "", err
	}
	return exactJSONNumber(value, raw)
}

func exactJSONNumber(value any, raw string) (json.Number, error) {
	number, ok := value.(json.Number)
	if !ok {
		return "", errExpectedJSONNumber
	}
	if number.String() != raw {
		return "", errExpectedJSONNumber
	}
	return number, nil
}

func validateDecimal128Number(raw string) error {
	significand, scale, err := decimal128Parts(raw)
	if err != nil || significand == "" {
		return err
	}
	significand, scale = normalizeDecimal128Parts(significand, scale)
	return validateDecimal128Parts(significand, scale)
}

func decimal128Parts(raw string) (string, int64, error) {
	unsigned := strings.TrimPrefix(raw, "-")
	exponent, unsigned, err := decimal128Exponent(unsigned)
	if err != nil {
		return "", 0, err
	}
	significand, fractionalDigits := decimal128Significand(unsigned)
	return strings.TrimLeft(significand, "0"), exponent - fractionalDigits, nil
}

func decimal128Exponent(unsigned string) (int64, string, error) {
	index := strings.IndexAny(unsigned, "eE")
	if index < 0 {
		return 0, unsigned, nil
	}
	exponent, err := strconv.ParseInt(unsigned[index+1:], decimal128ExponentBase, decimal128ExponentBits)
	if err != nil {
		return 0, "", errors.New("number exponent is outside the supported range")
	}
	return exponent, unsigned[:index], nil
}

func decimal128Significand(unsigned string) (string, int64) {
	index := strings.IndexByte(unsigned, '.')
	if index < 0 {
		return unsigned, 0
	}
	fractionalDigits := int64(len(unsigned) - index - 1)
	return unsigned[:index] + unsigned[index+1:], fractionalDigits
}

func normalizeDecimal128Parts(significand string, scale int64) (string, int64) {
	trimmed := strings.TrimRight(significand, "0")
	return trimmed, scale + int64(len(significand)-len(trimmed))
}

func validateDecimal128Parts(significand string, scale int64) error {
	if !fitsDecimal128Precision(significand, scale) {
		return errors.New("number cannot be represented faithfully as Decimal128")
	}
	if !fitsDecimal128UpperScale(significand, scale) {
		return errors.New("number cannot be represented faithfully as Decimal128")
	}
	return nil
}

func fitsDecimal128Precision(significand string, scale int64) bool {
	return len(significand) <= decimal128MaxDigits && scale >= decimal128MinScale
}

func fitsDecimal128UpperScale(significand string, scale int64) bool {
	return scale <= decimal128MaxScale || int64(len(significand))+scale-decimal128MaxScale <= decimal128MaxDigits
}
