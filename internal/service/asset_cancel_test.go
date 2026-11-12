package service

import (
	"context"
	"errors"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

func TestGetAssetPropagatesCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	storageClient := &storageFake{err: ctx.Err()}
	repository := &repositoryFake{document: &domain.AssetDocument{
		Filename: testAssetFilename,
		Location: "key",
		Access:   testAccessAllowed,
	}}
	assetService := newServiceForTest(t, repository, storageClient)

	if _, err := assetService.GetAsset(ctx, testAssetFilename, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation to propagate, got %v", err)
	}
	assertStorageContext(ctx, t, storageClient.call)
}
