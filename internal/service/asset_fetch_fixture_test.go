package service

import (
	"context"
	"testing"
	"time"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

func repositoryForLocation(location string) *repositoryFake {
	return &repositoryFake{document: &domain.AssetDocument{
		Filename: testAssetFilename,
		Location: location,
		Access:   testAccessAllowed,
	}}
}

func newStorageFixture(body *closeTracker, fixedTime time.Time) *storageFake {
	return &storageFake{output: &storage.GetObjectOutput{
		Body:          body,
		ContentType:   "application/javascript",
		ContentLength: testAssetContentLength,
		ETag:          "asset-etag",
		LastModified:  fixedTime,
	}}
}

func getAssetForTest(ctx context.Context, t *testing.T, service *AssetService) *domain.StaticAsset {
	t.Helper()
	asset, err := service.GetAsset(ctx, testAssetFilename, nil)
	if err != nil {
		t.Fatalf("GetAsset: %v", err)
	}
	return asset
}

func assertStreamOutput(t *testing.T, asset *domain.StaticAsset, body *closeTracker, fixedTime time.Time) {
	t.Helper()
	expected := domain.StaticAsset{
		LastModified:  fixedTime,
		Body:          body,
		ContentType:   "application/javascript",
		ETag:          "asset-etag",
		ContentLength: testAssetContentLength,
	}
	assertStaticAsset(t, asset, expected)
	assertBodyOpen(t, body)
	closeAssetBody(t, asset.Body)
	assertBodyClosed(t, body)
}
