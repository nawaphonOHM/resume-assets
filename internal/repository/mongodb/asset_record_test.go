package mongodb

import (
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	libmongodb "github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

func TestAssetRecordStringFieldsRoundTrip(t *testing.T) {
	original := assetRecord{Filename: testFilename, Location: "s3://bucket/key", Access: "ALLOW"}
	encoded, err := libmongodb.Marshal(original)
	if err != nil {
		t.Fatalf("marshal asset record: %v", err)
	}
	document := decodeAssetRecord(t, encoded).toAssetDocument()
	expected := domain.AssetDocument{Filename: testFilename, Location: "s3://bucket/key", Access: "ALLOW"}
	if *document != expected {
		t.Fatalf("unexpected metadata: got %#v, want %#v", document, expected)
	}
}

func TestAssetRecordUntrustedFieldsFailClosed(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]any
	}{
		{name: "non-string access", fields: map[string]any{"access": false}},
		{name: "non-string location", fields: map[string]any{"location": 42}},
		{name: "null fields", fields: map[string]any{"access": nil, "location": nil}},
		{name: "missing fields", fields: map[string]any{}},
		{
			name: "non-string access and location",
			fields: map[string]any{"access": false, "location": 42},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := map[string]any{"filename": testFilename}
			for key, value := range test.fields {
				original[key] = value
			}
			encoded, err := libmongodb.Marshal(original)
			if err != nil {
				t.Fatalf("marshal asset record: %v", err)
			}
			document := decodeAssetRecord(t, encoded).toAssetDocument()
			assertUntrustedFields(t, document)
		})
	}
}

func assertUntrustedFields(t *testing.T, document *domain.AssetDocument) {
	t.Helper()
	if document.Access != "" {
		t.Fatalf("non-string access must not be trusted: %#v", document)
	}
	if document.Location != "" {
		t.Fatalf("non-string location must not be trusted: %#v", document)
	}
}

func decodeAssetRecord(t *testing.T, data []byte) assetRecord {
	t.Helper()
	var record assetRecord
	if err := libmongodb.Unmarshal(data, &record); err != nil {
		t.Fatalf("unmarshal asset record: %v", err)
	}
	return record
}
