package mongodb

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
)

const (
	integrationSameFilename    = "same.js"
	integrationVariantKey      = "variant"
	integrationNestedKey       = "nested"
	integrationStatusKey       = "status"
	integrationStatusAvailable = "available"
)

func insertMongoIntegrationFixtures(ctx context.Context, t *testing.T, collection *mongodb.Collection) {
	t.Helper()
	if _, err := collection.InsertMany(ctx, mongoIntegrationDocuments()); err != nil {
		t.Fatalf("insert integration fixtures: %v", err)
	}
}

func mongoIntegrationDocuments() []any {
	storedNested := mongodb.D{
		{Key: "sequence", Value: mongodb.A{int32(1), mongodb.D{{Key: "y", Value: true}, {Key: "x", Value: false}}}},
		{Key: "payload", Value: mongodb.D{{Key: "second", Value: "two"}, {Key: "first", Value: "one"}}},
	}
	arrayOrderNested := mongodb.D{
		{Key: "sequence", Value: mongodb.A{mongodb.D{{Key: "x", Value: false}, {Key: "y", Value: true}}, int32(1)}},
		{Key: "payload", Value: mongodb.D{{Key: "first", Value: "one"}, {Key: "second", Value: "two"}}},
	}
	return []any{
		integrationAsset(integrationSameFilename, integrationVariantTags("single", storedNested)),
		integrationAsset(integrationSameFilename, integrationVariantTags("array-order", arrayOrderNested)),
		integrationAsset(integrationSameFilename, integrationVariantTags("duplicate", storedNested)),
		integrationAsset(integrationSameFilename, integrationVariantTags("duplicate", storedNested)),
		integrationAsset("empty-tags.js", mongodb.D{}),
		integrationAsset("missing-tags.js", nil),
		integrationAsset("null-vs-missing.js", mongodb.D{{Key: "nullable", Value: nil}}),
		integrationAsset("null-vs-missing.js", mongodb.D{{Key: "other", Value: "present"}}),
	}
}

func integrationVariantTags(variant string, nested mongodb.D) mongodb.D {
	return mongodb.D{
		{Key: integrationVariantKey, Value: variant},
		{Key: integrationNestedKey, Value: nested},
		{Key: integrationStatusKey, Value: integrationStatusAvailable},
	}
}

func integrationNestedCriteria() map[string]any {
	return map[string]any{
		"payload":  map[string]any{"first": "one", "second": "two"},
		"sequence": []any{json.Number("1"), map[string]any{"x": false, "y": true}},
	}
}
