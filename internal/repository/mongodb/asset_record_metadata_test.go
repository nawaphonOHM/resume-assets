package mongodb

import (
	"reflect"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	libmongodb "github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

const (
	assetMetadataLength      = 1296
	assetMetadataContentType = "image/svg+xml"
)

func TestAssetRecordStringFieldsRoundTrip(t *testing.T) {
	original := assetRecord{
		Filename: testFilename,
		Location: "s3://bucket/key",
		Access:   "ALLOW",
		Metadata: domain.AssetMetadata{
			ForCaching: map[string]any{
				"ETag":          "CJzyvtqFq5cDEAE=",
				"Last-Modified": "2026-10-08T18:25:30+0000",
			},
			ForDownloading: map[string]any{
				"Content-Length": int32(assetMetadataLength),
				"Content-Type":   assetMetadataContentType,
			},
		},
	}
	encoded, err := libmongodb.Marshal(original)
	if err != nil {
		t.Fatalf("marshal asset record: %v", err)
	}
	document := decodeAssetRecord(t, encoded).toAssetDocument()
	if !reflect.DeepEqual(*document, expectedAssetMetadataDocument()) {
		t.Fatalf("unexpected metadata: got %#v, want %#v", document, expectedAssetMetadataDocument())
	}
}

func expectedAssetMetadataDocument() domain.AssetDocument {
	return domain.AssetDocument{
		Filename: testFilename,
		Location: "s3://bucket/key",
		Access:   "ALLOW",
		Metadata: domain.AssetMetadata{
			ForCaching: map[string]any{
				"ETag":          "CJzyvtqFq5cDEAE=",
				"Last-Modified": "2026-10-08T18:25:30+0000",
			},
			ForDownloading: map[string]any{
				"Content-Length": int32(assetMetadataLength),
				"Content-Type":   assetMetadataContentType,
			},
		},
	}
}
