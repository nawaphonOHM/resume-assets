package mongodb

import (
	"context"
	"errors"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	libmongodb "github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

const testFilename = "app.js"

func TestFindAssetNilCollection(t *testing.T) {
	repository := NewMongoAssetRepository(nil)
	asset, err := repository.FindAsset(context.Background(), testFilename, nil)
	if asset != nil {
		t.Fatalf("expected nil asset, got %#v", asset)
	}
	if !errors.Is(err, ErrNilCollection) {
		t.Fatalf("expected ErrNilCollection, got %v", err)
	}
}

func TestMapFindErrorMissingDocument(t *testing.T) {
	if err := mapFindError(libmongodb.ErrNoDocuments, testFilename); !errors.Is(err, domain.ErrAssetNotFound) {
		t.Fatalf("expected ErrAssetNotFound, got %v", err)
	}
}

func TestMapFindErrorBackend(t *testing.T) {
	backendErr := errors.New("database unavailable")
	mapped := mapFindError(backendErr, testFilename)
	if !errors.Is(mapped, backendErr) {
		t.Fatalf("expected backend error to be wrapped, got %v", mapped)
	}
	if errors.Is(mapped, domain.ErrAssetNotFound) {
		t.Fatalf("backend error must not be mapped to missing metadata: %v", mapped)
	}
}

func TestAssetDocumentBSON(t *testing.T) {
	asset := domain.AssetDocument{
		Filename: testFilename,
		Location: "s3://assets-bucket/app.js",
		Access:   "ALLOW",
	}
	data, err := libmongodb.Marshal(asset)
	if err != nil {
		t.Fatalf("marshal asset document: %v", err)
	}
	decoded := decodeAssetDocument(t, data)
	if decoded != asset {
		t.Fatalf("round-trip mismatch: got %#v, want %#v", decoded, asset)
	}
}

func decodeAssetDocument(t *testing.T, data []byte) domain.AssetDocument {
	t.Helper()
	var decoded domain.AssetDocument
	if err := libmongodb.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal asset document: %v", err)
	}
	return decoded
}
