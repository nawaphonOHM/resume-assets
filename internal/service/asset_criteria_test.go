package service

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

func TestParseAssetCriteriaSupportsJSONValuesAndPreservesNumbers(t *testing.T) {
	raw := `{"text":"available","number":9007199254740993,"decimal":1.25e2,` +
		`"enabled":true,"nothing":null,"array":["x",false,{"nested":1}],"object":{"b":2,"a":1}}`
	criteria, err := parseAssetCriteria(raw)
	if err != nil {
		t.Fatalf("parseAssetCriteria: %v", err)
	}
	want := domain.AssetCriteria{
		"text":    "available",
		"number":  json.Number("9007199254740993"),
		"decimal": json.Number("1.25e2"),
		"enabled": true,
		"nothing": nil,
		"array":   []any{"x", false, map[string]any{"nested": json.Number("1")}},
		"object":  map[string]any{"b": json.Number("2"), "a": json.Number("1")},
	}
	if !reflect.DeepEqual(criteria, want) {
		t.Fatalf("unexpected criteria: got %#v, want %#v", criteria, want)
	}
}

func TestParseAssetCriteriaAcceptsEmptyObject(t *testing.T) {
	criteria, err := parseAssetCriteria(`{}`)
	if err != nil {
		t.Fatalf("parse empty criteria: %v", err)
	}
	assertEmptyCriteria(t, criteria)
}

func assertEmptyCriteria(t *testing.T, criteria domain.AssetCriteria) {
	t.Helper()
	if criteria == nil || len(criteria) != 0 {
		t.Fatalf("expected a non-nil empty object, got %#v", criteria)
	}
}

func TestParseAssetCriteriaRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{name: "malformed JSON", text: `{"status":`},
		{name: "array root", text: `[]`},
		{name: "null root", text: `null`},
		{name: "string root", text: `"available"`},
		{name: "duplicate root key", text: `{"status":"available","status":"hidden"}`},
		{name: "duplicate nested key", text: `{"nested":{"value":1,"value":2}}`},
		{name: "multiple JSON values", text: `{} {}`},
		{name: "unrepresentable number", text: `{"number":1e7000}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := parseAssetCriteria(test.text); !errors.Is(err, domain.ErrInvalidAssetCriteria) {
				t.Fatalf("expected ErrInvalidAssetCriteria, got %v", err)
			}
		})
	}
}

func stringPointer(value string) *string { return &value }
