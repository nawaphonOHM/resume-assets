package service

import (
	"errors"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

type assetErrorCase struct {
	storageErr error
	wantErr    error
	extraErr   error
	output     *storage.GetObjectOutput
	name       string
	location   string
	wantCalls  int
}

func TestGetAssetErrorCases(t *testing.T) {
	const objectLocation = "s3://bucket/key"
	backendErr := errors.New("storage unavailable")
	tests := []assetErrorCase{
		{
			name: "invalid location", location: "https://bucket/key",
			wantErr: domain.ErrInvalidAssetMetadata, extraErr: ErrInvalidLocation,
		},
		{
			name: "empty location", wantErr: domain.ErrInvalidAssetMetadata,
			extraErr: ErrInvalidLocation,
		},
		{
			name: "missing object", location: objectLocation, storageErr: storage.ErrObjectNotFound,
			wantErr: domain.ErrAssetNotFound, wantCalls: 1,
		},
		{
			name: "storage backend error", location: objectLocation,
			storageErr: backendErr, wantErr: backendErr, wantCalls: 1,
		},
		{name: "nil storage output", location: objectLocation, wantErr: ErrInvalidStorageResponse, wantCalls: 1},
		{
			name: "nil object body", location: objectLocation, output: &storage.GetObjectOutput{},
			wantErr: ErrInvalidStorageResponse, wantCalls: 1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runAssetErrorCase(t, test)
		})
	}
}

func runAssetErrorCase(t *testing.T, test assetErrorCase) {
	t.Helper()
	repository := &repositoryFake{document: &domain.AssetDocument{
		Filename: testAssetFilename,
		Location: test.location,
		Access:   testAccessAllowed,
	}}
	storageClient := &storageFake{output: test.output, err: test.storageErr}
	assetService := newServiceForTest(t, repository, storageClient)

	asset, err := assetService.GetAsset(t.Context(), testAssetFilename, nil)
	assertNoStaticAsset(t, asset)
	assertAssetError(t, err, test.wantErr)
	assertStorageCallCount(t, storageClient, test.wantCalls)
	assertAdditionalError(t, err, test.extraErr)
}

func assertNoStaticAsset(t *testing.T, asset *domain.StaticAsset) {
	t.Helper()
	if asset != nil {
		t.Fatalf("expected nil asset, got %#v", asset)
	}
}

func assertAssetError(t *testing.T, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Fatalf("expected error %v, got %v", want, got)
	}
}

func assertStorageCallCount(t *testing.T, client *storageFake, want int) {
	t.Helper()
	if client.calls != want {
		t.Fatalf("expected %d storage calls, got %d", want, client.calls)
	}
}

func assertAdditionalError(t *testing.T, err, extra error) {
	t.Helper()
	if extra == nil {
		return
	}
	if !errors.Is(err, extra) {
		t.Fatalf("expected additional error %v, got %v", extra, err)
	}
}
