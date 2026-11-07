package service

import (
	"errors"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

type accessGateCase struct {
	repoErr  error
	wantErr  error
	document *domain.AssetDocument
	name     string
}

func TestGetAssetDoesNotFetchWithoutPermission(t *testing.T) {
	backendErr := errors.New("database unavailable")
	tests := []accessGateCase{
		{
			name: "denied despite invalid location", wantErr: domain.ErrAssetAccessDenied,
			document: &domain.AssetDocument{Access: "DENY", Location: "invalid://location"},
		},
		{name: "missing access", document: &domain.AssetDocument{}, wantErr: domain.ErrAssetAccessDenied},
		{
			name: "lowercase access", document: &domain.AssetDocument{Access: "allow"},
			wantErr: domain.ErrAssetAccessDenied,
		},
		{
			name: "mismatched filename", wantErr: domain.ErrInvalidAssetMetadata,
			document: &domain.AssetDocument{Filename: "other.js", Access: testAccessAllowed},
		},
		{name: "missing metadata", repoErr: domain.ErrAssetNotFound, wantErr: domain.ErrAssetNotFound},
		{name: "backend failure", repoErr: backendErr, wantErr: backendErr},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runAccessGateCase(t, test)
		})
	}
}

func runAccessGateCase(t *testing.T, test accessGateCase) {
	t.Helper()
	repository := &repositoryFake{document: test.document, err: test.repoErr}
	storageClient := &storageFake{}
	assetService := newServiceForTest(t, repository, storageClient)
	asset, err := assetService.GetAsset(t.Context(), testAssetFilename)
	assertNoStaticAsset(t, asset)
	assertPolicyError(t, err, test.wantErr)
	assertStorageNotCalled(t, storageClient)
}
