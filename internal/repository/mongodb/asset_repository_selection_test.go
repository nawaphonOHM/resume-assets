package mongodb

import (
	"errors"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

func TestUniqueAssetRecordRejectsNoMatches(t *testing.T) {
	asset, err := uniqueAssetRecord(nil, testFilename)
	if asset != nil {
		t.Fatalf("expected no asset, got %#v", asset)
	}
	if !errors.Is(err, domain.ErrAssetNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestUniqueAssetRecordReturnsSingleMatch(t *testing.T) {
	record := assetRecord{Filename: testFilename, Location: "s3://assets/app.js", Access: "ALLOW"}
	asset, err := uniqueAssetRecord([]assetRecord{record}, testFilename)
	if err != nil {
		t.Fatalf("select single match: %v", err)
	}
	assertSelectedAssetRecord(t, asset)
}

func TestUniqueAssetRecordRejectsAmbiguousMatches(t *testing.T) {
	records := []assetRecord{{Filename: testFilename}, {Filename: testFilename}}
	asset, err := uniqueAssetRecord(records, testFilename)
	if asset != nil {
		t.Fatalf("expected no asset, got %#v", asset)
	}
	if !errors.Is(err, domain.ErrAmbiguousAssetMatch) {
		t.Fatalf("expected ambiguous match, got %v", err)
	}
}

func assertSelectedAssetRecord(t *testing.T, asset *domain.AssetDocument) {
	t.Helper()
	if asset == nil {
		t.Fatal("expected selected asset")
	}
	want := domain.AssetDocument{Filename: testFilename, Location: "s3://assets/app.js", Access: "ALLOW"}
	if *asset != want {
		t.Fatalf("selected asset = %#v, want %#v", asset, want)
	}
}
