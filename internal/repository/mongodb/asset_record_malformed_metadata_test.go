package mongodb

import (
	"reflect"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	libmongodb "github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

type malformedMetadataCase struct {
	metadata        any
	want            domain.AssetMetadata
	includeMetadata bool
}

func TestAssetRecordMalformedMetadataIsIgnored(t *testing.T) {
	assertMalformedMetadata(t, "missing", malformedMetadataCase{})
	assertMalformedMetadata(t, "null", malformedMetadataCase{includeMetadata: true})
	assertMalformedMetadata(t, "malformed root", malformedMetadataCase{
		metadata: "not a document", includeMetadata: true,
	})
	assertMalformedMetadata(t, "malformed header map", malformedMetadataCase{
		metadata: map[string]any{
			"forCaching":     "not a document",
			"forDownloading": map[string]any{"Content-Type": "application/octet-stream"},
		},
		includeMetadata: true,
		want: domain.AssetMetadata{
			ForDownloading: map[string]any{"Content-Type": "application/octet-stream"},
		},
	})
}

func assertMalformedMetadata(t *testing.T, name string, test malformedMetadataCase) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		checkMalformedMetadata(t, test)
	})
}

func checkMalformedMetadata(t *testing.T, test malformedMetadataCase) {
	t.Helper()
	document := map[string]any{"filename": testFilename}
	addMalformedMetadata(document, test)
	asset := decodeAssetRecord(t, marshalMalformedMetadata(t, document)).toAssetDocument()
	assertAssetMetadataEqual(t, asset.Metadata, test.want)
}

func addMalformedMetadata(document map[string]any, test malformedMetadataCase) {
	if test.includeMetadata {
		document["metaData"] = test.metadata
	}
}

func marshalMalformedMetadata(t *testing.T, document map[string]any) []byte {
	t.Helper()
	encoded, err := libmongodb.Marshal(document)
	if err != nil {
		t.Fatalf("marshal asset record: %v", err)
	}
	return encoded
}

func assertAssetMetadataEqual(t *testing.T, got, want domain.AssetMetadata) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected metadata: got %#v, want %#v", got, want)
	}
}
