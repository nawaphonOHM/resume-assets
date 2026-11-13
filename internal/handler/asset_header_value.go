package handler

import (
	"reflect"
	"strconv"
)

func assetHeaderValue(candidate any) (string, bool) {
	if text, ok := candidate.(string); ok {
		return text, true
	}
	return integerHeaderValue(candidate)
}

func integerHeaderValue(candidate any) (string, bool) {
	number := reflect.ValueOf(candidate)
	if !number.IsValid() {
		return "", false
	}
	return formatIntegerHeaderValue(number)
}

func formatIntegerHeaderValue(number reflect.Value) (string, bool) {
	if isSignedIntegerKind(number.Kind()) {
		return strconv.FormatInt(number.Int(), decimalBase), true
	}
	if isUnsignedIntegerKind(number.Kind()) {
		return strconv.FormatUint(number.Uint(), decimalBase), true
	}
	return "", false
}

func isSignedIntegerKind(kind reflect.Kind) bool {
	return kind >= reflect.Int && kind <= reflect.Int64
}

func isUnsignedIntegerKind(kind reflect.Kind) bool {
	return kind >= reflect.Uint && kind <= reflect.Uint64 || kind == reflect.Uintptr
}
